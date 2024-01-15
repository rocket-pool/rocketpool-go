package protocol

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// Settings
const (
	proposalBatchSize int = 100
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProtocol
type ProtocolDaoManager struct {
	// Settings for the Protocol DAO
	Settings *ProtocolDaoSettings

	// The time that the RPL rewards percentages were last updated
	LastRewardsPercentagesUpdate *core.FormattedUint256Field[time.Time]

	// Get the total number of Protocol DAO proposals
	ProposalCount *core.FormattedUint256Field[uint64]

	// The depth of a network or node voting tree pollard for each round of challenge / response
	DepthPerRound *core.FormattedUint256Field[uint64]

	// === Internal fields ===
	rp   *rocketpool.RocketPool
	cd   *core.Contract
	dp   *core.Contract
	dpp  *core.Contract
	dpps *core.Contract
	dpsr *core.Contract
	dpv  *core.Contract
}

// Rewards claimer percents
type RplRewardsPercentages struct {
	OdaoPercentage *big.Int `abi:"_trustedNodePercent"`
	PdaoPercentage *big.Int `abi:"_protocolPercent"`
	NodePercentage *big.Int `abi:"_nodePercent"`
}

// Structure of the RootSubmitted event
type RootSubmitted struct {
	ProposalID  *big.Int               `json:"proposalId"`
	Proposer    common.Address         `json:"proposer"`
	BlockNumber uint32                 `json:"blockNumber"`
	Index       *big.Int               `json:"index"`
	Root        types.VotingTreeNode   `json:"root"`
	TreeNodes   []types.VotingTreeNode `json:"treeNodes"`
	Timestamp   time.Time              `json:"timestamp"`
}

// Internal struct - returned by the RootSubmitted event
type rootSubmittedRaw struct {
	ProposalID  *big.Int               `json:"proposalId"`
	Proposer    common.Address         `json:"proposer"`
	BlockNumber uint32                 `json:"blockNumber"`
	Index       *big.Int               `json:"index"`
	Root        types.VotingTreeNode   `json:"root"`
	TreeNodes   []types.VotingTreeNode `json:"treeNodes"`
	Timestamp   *big.Int               `json:"timestamp"`
}

// Structure of the ChallengeSubmitted event
type ChallengeSubmitted struct {
	ProposalID *big.Int       `json:"proposalId"`
	Challenger common.Address `json:"challenger"`
	Index      *big.Int       `json:"index"`
	Timestamp  time.Time      `json:"timestamp"`
}

// Internal struct - returned by the ChallengeSubmitted event
type challengeSubmittedRaw struct {
	ProposalID *big.Int       `json:"proposalId"`
	Challenger common.Address `json:"challenger"`
	Index      *big.Int       `json:"index"`
	Timestamp  *big.Int       `json:"timestamp"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProtocolDaoManager contract binding
func NewProtocolDaoManager(rp *rocketpool.RocketPool) (*ProtocolDaoManager, error) {
	// Create the contracts
	cd, err := rp.GetContract(rocketpool.ContractName_RocketClaimDAO)
	if err != nil {
		return nil, fmt.Errorf("error getting claim DAO contract: %w", err)
	}
	dp, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocol)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO manager contract: %w", err)
	}
	dpp, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocolProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO protocol proposal contract: %w", err)
	}
	dpps, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocolProposals)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO protocol proposals contract: %w", err)
	}
	dpsr, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocolSettingsRewards)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO protocol settings rewards contract: %w", err)
	}
	dpv, err := rp.GetContract(rocketpool.ContractName_RocketDAOProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO protocol verifier contract: %w", err)
	}

	pdaoMgr := &ProtocolDaoManager{
		LastRewardsPercentagesUpdate: core.NewFormattedUint256Field[time.Time](dpsr, "getRewardsClaimersTimeUpdated"),
		ProposalCount:                core.NewFormattedUint256Field[uint64](dpp, "getTotal"),
		DepthPerRound:                core.NewFormattedUint256Field[uint64](dpv, "getDepthPerRound"),

		rp:   rp,
		cd:   cd,
		dp:   dp,
		dpp:  dpp,
		dpps: dpps,
		dpsr: dpsr,
		dpv:  dpv,
	}
	settings, err := newProtocolDaoSettings(pdaoMgr)
	if err != nil {
		return nil, fmt.Errorf("error creating Protocol DAO settings binding: %w", err)
	}
	pdaoMgr.Settings = settings
	return pdaoMgr, nil
}

// =============
// === Calls ===
// =============

// === ClaimDAO ===

// Check if a recurring spend exists with the given contract name
func (c *ProtocolDaoManager) GetContractExists(mc *batch.MultiCaller, out *bool, contractName string) {
	core.AddCall(mc, c.cd, out, "getContractExists", contractName)
}

// === DAOProtocolSettingsRewards ===

// Get the allocation of RPL rewards to the node operators, Oracle DAO, and the Protocol DAO
func (c *ProtocolDaoManager) GetRewardsPercentages(mc *batch.MultiCaller, out *RplRewardsPercentages) {
	core.AddCallRaw(mc, c.dpsr, out, "getRewardsClaimersPerc")
}

// Get the allocation of RPL rewards to the node operators
func (c *ProtocolDaoManager) GetNodeOperatorRewardsPercent(mc *batch.MultiCaller, out **big.Int) {
	core.AddCall(mc, c.dpsr, out, "getRewardsClaimersNodePerc")
}

// Get the allocation of RPL rewards to the Oracle DAO
func (c *ProtocolDaoManager) GetOracleDaoRewardsPercent(mc *batch.MultiCaller, out **big.Int) {
	core.AddCall(mc, c.dpsr, out, "getRewardsClaimersTrustedNodePerc")
}

// Get the allocation of RPL rewards to the Protocol DAO
func (c *ProtocolDaoManager) GetProtocolDaoRewardsPercent(mc *batch.MultiCaller, out **big.Int) {
	core.AddCall(mc, c.dpsr, out, "getRewardsClaimersProtocolPerc")
}

// ====================
// === Transactions ===
// ====================

// === DAOProtocol ===

// Get info for bootstrapping a bool setting
func (c *ProtocolDaoManager) BootstrapBool(contractName rocketpool.ContractName, setting SettingName, value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingBool", opts, contractName, string(setting), value)
}

// Get info for bootstrapping a uint256 setting
func (c *ProtocolDaoManager) BootstrapUint(contractName rocketpool.ContractName, setting SettingName, value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingUint", opts, contractName, string(setting), value)
}

// Get info for bootstrapping an address setting
func (c *ProtocolDaoManager) BootstrapAddress(contractName rocketpool.ContractName, setting SettingName, value common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingAddress", opts, contractName, string(setting), value)
}

// Get info for bootstrapping a rewards claimer
func (c *ProtocolDaoManager) BootstrapClaimer(contractName rocketpool.ContractName, amount float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingClaimer", opts, contractName, eth.EthToWei(amount))
}

// === DAOProtocolProposals ===

// Get info for submitting a proposal to update a bool Protocol DAO setting
func (c *ProtocolDaoManager) ProposeSetBool(message string, contractName rocketpool.ContractName, setting SettingName, value bool, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("set %s", setting)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSettingBool", contractName, string(setting), value)
}

// Get info for submitting a proposal to update a uint Protocol DAO setting
func (c *ProtocolDaoManager) ProposeSetUint(message string, contractName rocketpool.ContractName, setting SettingName, value *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("set %s", setting)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSettingUint", contractName, string(setting), value)
}

// Get info for submitting a proposal to update an address Protocol DAO setting
func (c *ProtocolDaoManager) ProposeSetAddress(message string, contractName rocketpool.ContractName, setting SettingName, value common.Address, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("set %s", setting)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSettingAddress", contractName, string(setting), value)
}

// Get info for submitting a proposal to update multiple Protocol DAO settings at once
func (c *ProtocolDaoManager) ProposeSetMulti(message string, contractNames []rocketpool.ContractName, settings []SettingName, settingTypes []types.ProposalSettingType, values []any, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	settingNameStrings := make([]string, len(settings))
	for i, setting := range settings {
		settingNameStrings[i] = string(setting)
	}
	if message == "" {
		message = fmt.Sprintf("set %s", strings.Join(settingNameStrings, ", "))
	}
	encodedValues, err := abiEncodeMultiValues(settingTypes, values)
	if err != nil {
		return nil, fmt.Errorf("error ABI encoding values: %w", err)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSettingMulti", contractNames, settingNameStrings, settingTypes, encodedValues)
}

// Get info for submitting a proposal to update the allocations of RPL rewards
func (c *ProtocolDaoManager) ProposeSetRewardsPercentages(message string, odaoPercentage *big.Int, pdaoPercentage *big.Int, nodePercentage *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = "set rewards percentages"
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSettingRewardsClaimers", odaoPercentage, pdaoPercentage, nodePercentage)
}

// Get info for submitting a proposal to spend a portion of the Rocket Pool treasury one time
func (c *ProtocolDaoManager) ProposeOneTimeTreasurySpend(message, invoiceID string, recipient common.Address, amount *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("propose one-time treasury spend - invoice %s", invoiceID)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalTreasuryOneTimeSpend", invoiceID, recipient, amount)
}

// Get info for submitting a proposal to spend a portion of the Rocket Pool treasury in a recurring manner
func (c *ProtocolDaoManager) ProposeRecurringTreasurySpend(message string, contractName string, recipient common.Address, amountPerPeriod *big.Int, periodLength time.Duration, startTime time.Time, numberOfPeriods uint64, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("propose recurring treasury spend - contract %s", contractName)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalTreasuryNewContract", contractName, recipient, amountPerPeriod, big.NewInt(int64(periodLength.Seconds())), big.NewInt(startTime.Unix()), numberOfPeriods)
}

// Get info for submitting a proposal to update a recurring Rocket Pool treasury spending plan
func (c *ProtocolDaoManager) ProposeRecurringTreasurySpendUpdate(message string, contractName string, recipient common.Address, amountPerPeriod *big.Int, periodLength time.Duration, numberOfPeriods uint64, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("propose recurring treasury spend update - contract %s", contractName)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalTreasuryUpdateContract", contractName, recipient, amountPerPeriod, big.NewInt(int64(periodLength.Seconds())), numberOfPeriods)
}

// Get info for submitting a proposal to invite a member to the security council
func (c *ProtocolDaoManager) ProposeInviteToSecurityCouncil(message string, id string, address common.Address, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("invite %s (%s) to the security council", id, address.Hex())
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSecurityInvite", id, address)
}

// Get info for submitting a proposal to kick a member from the security council
func (c *ProtocolDaoManager) ProposeKickFromSecurityCouncil(message string, address common.Address, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("kick %s from the security council", address.Hex())
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSecurityKick", address)
}

// Get info for submitting a proposal to kick multiple members from the security council
func (c *ProtocolDaoManager) ProposeKickMultiFromSecurityCouncil(message string, addresses []common.Address, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = "kick multiple members from the security council"
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSecurityKickMulti", addresses)
}

// Get info for submitting a proposal to replace a member of the security council with another one in a single TX
func (c *ProtocolDaoManager) ProposeReplaceSecurityCouncilMember(message string, existingMemberAddress common.Address, newMemberID string, newMemberAddress common.Address, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("replace %s on the security council with %s (%s)", existingMemberAddress.Hex(), newMemberID, newMemberAddress.Hex())
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSecurityReplace", existingMemberAddress, newMemberID, newMemberAddress)
}

// Submit a protocol DAO proposal
func (c *ProtocolDaoManager) submitProposal(opts *bind.TransactOpts, blockNumber uint32, treeNodes []types.VotingTreeNode, message string, method string, args ...interface{}) (*core.TransactionInfo, error) {
	payload, err := c.dpps.ABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("error encoding payload: %w", err)
	}
	err = c.simulateProposalExecution(payload)
	if err != nil {
		return nil, fmt.Errorf("error simulating proposal execution: %w", err)
	}
	return core.NewTransactionInfo(c.dpp, "propose", opts, message, payload, blockNumber, treeNodes)
}

/// =============
/// === Utils ===
/// =============

// Get RootSubmitted event info
func (c *ProtocolDaoManager) GetRootSubmittedEvents(proposalIDs []uint64, intervalSize *big.Int, startBlock *big.Int, endBlock *big.Int, opts *bind.CallOpts) ([]RootSubmitted, error) {
	// Construct a filter query for relevant logs
	idBuffers := make([]common.Hash, len(proposalIDs))
	for i, id := range proposalIDs {
		proposalIdBig := big.NewInt(0).SetUint64(id)
		proposalIdBig.FillBytes(idBuffers[i].Bytes())
	}
	rootSubmittedEvent := c.dpv.ABI.Events["RootSubmitted"]
	addressFilter := []common.Address{*c.dpv.Address}
	topicFilter := [][]common.Hash{{rootSubmittedEvent.ID}, idBuffers}

	// Get the event logs
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, startBlock, endBlock, nil)
	if err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return []RootSubmitted{}, nil
	}

	events := make([]RootSubmitted, 0, len(logs))
	for _, log := range logs {
		// Get the log info values
		values, err := rootSubmittedEvent.Inputs.Unpack(log.Data)
		if err != nil {
			return nil, fmt.Errorf("error unpacking RootSubmitted event data: %w", err)
		}

		// Convert to a native struct
		var raw rootSubmittedRaw
		err = rootSubmittedEvent.Inputs.Copy(&raw, values)
		if err != nil {
			return nil, fmt.Errorf("error converting RootSubmitted event data to struct: %w", err)
		}

		// Get the decoded data
		events = append(events, RootSubmitted{
			ProposalID:  raw.ProposalID,
			Proposer:    raw.Proposer,
			BlockNumber: raw.BlockNumber,
			Index:       raw.Index,
			Root:        raw.Root,
			TreeNodes:   raw.TreeNodes,
			Timestamp:   time.Unix(raw.Timestamp.Int64(), 0),
		})
	}

	return events, nil
}

// Get ChallengeSubmitted event info
func (c *ProtocolDaoManager) GetChallengeSubmittedEvents(proposalIDs []uint64, intervalSize *big.Int, startBlock *big.Int, endBlock *big.Int, opts *bind.CallOpts) ([]ChallengeSubmitted, error) {
	// Construct a filter query for relevant logs
	idBuffers := make([]common.Hash, len(proposalIDs))
	for i, id := range proposalIDs {
		proposalIdBig := big.NewInt(0).SetUint64(id)
		proposalIdBig.FillBytes(idBuffers[i].Bytes())
	}
	challengeSubmittedEvent := c.dpv.ABI.Events["ChallengeSubmitted"]
	addressFilter := []common.Address{*c.dpv.Address}
	topicFilter := [][]common.Hash{{challengeSubmittedEvent.ID}, idBuffers}

	// Get the event logs
	logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, startBlock, endBlock, nil)
	if err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return []ChallengeSubmitted{}, nil
	}

	events := make([]ChallengeSubmitted, 0, len(logs))
	for _, log := range logs {
		// Get the log info values
		values, err := challengeSubmittedEvent.Inputs.Unpack(log.Data)
		if err != nil {
			return nil, fmt.Errorf("error unpacking ChallengeSubmitted event data: %w", err)
		}

		// Convert to a native struct
		var raw challengeSubmittedRaw
		err = challengeSubmittedEvent.Inputs.Copy(&raw, values)
		if err != nil {
			return nil, fmt.Errorf("error converting ChallengeSubmitted event data to struct: %w", err)
		}

		// Get the decoded data
		events = append(events, ChallengeSubmitted{
			ProposalID: raw.ProposalID,
			Challenger: raw.Challenger,
			Index:      raw.Index,
			Timestamp:  time.Unix(raw.Timestamp.Int64(), 0),
		})
	}

	return events, nil
}

// Get all proposal details
func (c *ProtocolDaoManager) GetProposals(proposalCount uint64, includeDetails bool, opts *bind.CallOpts) ([]*ProtocolDaoProposal, error) {
	// Create prop commons for each one
	props := make([]*ProtocolDaoProposal, proposalCount)
	for i := uint64(1); i <= proposalCount; i++ { // Proposals are 1-indexed
		prop, err := NewProtocolDaoProposal(c.rp, i)
		if err != nil {
			return nil, fmt.Errorf("error creating Protocol DAO proposal %d: %w", i, err)
		}
		props[i-1] = prop
	}

	// Get all details if requested
	if includeDetails {
		err := c.rp.BatchQuery(int(proposalCount), proposalBatchSize, func(mc *batch.MultiCaller, index int) error {
			core.QueryAllFields(props[index], mc)
			return nil
		}, opts)
		if err != nil {
			return nil, fmt.Errorf("error getting proposal details: %w", err)
		}
	}

	// Return
	return props, nil
}

// Simulate a proposal's execution to verify it won't revert
func (c *ProtocolDaoManager) simulateProposalExecution(payload []byte) error {
	_, err := c.rp.Client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:     *c.dpp.Address,
		To:       c.dpps.Address,
		GasPrice: big.NewInt(0),
		Value:    nil,
		Data:     payload,
	})
	return err
}

// Get the ABI encoding of multiple values for a ProposeSettingMulti call
func abiEncodeMultiValues(settingTypes []types.ProposalSettingType, values []any) ([][]byte, error) {
	// Sanity check the lengths
	settingCount := len(settingTypes)
	if settingCount != len(values) {
		return nil, fmt.Errorf("settingTypes and values must be the same length")
	}
	if settingCount == 0 {
		return [][]byte{}, nil
	}

	// ABI encode each value
	results := make([][]byte, settingCount)
	for i, settingType := range settingTypes {
		var encodedArg []byte
		switch settingType {
		case types.ProposalSettingType_Uint256:
			arg, success := values[i].(*big.Int)
			if !success {
				return nil, fmt.Errorf("value %d is not a *big.Int, but the setting type is Uint256", i)
			}
			encodedArg = math.U256Bytes(big.NewInt(0).Set(arg))

		case types.ProposalSettingType_Bool:
			arg, success := values[i].(bool)
			if !success {
				return nil, fmt.Errorf("value %d is not a bool, but the setting type is Bool", i)
			}
			if arg {
				encodedArg = math.PaddedBigBytes(common.Big1, 32)
			} else {
				encodedArg = math.PaddedBigBytes(common.Big0, 32)
			}

		case types.ProposalSettingType_Address:
			arg, success := values[i].(common.Address)
			if !success {
				return nil, fmt.Errorf("value %d is not an address, but the setting type is Address", i)
			}
			encodedArg = common.LeftPadBytes(arg.Bytes(), 32)

		default:
			return nil, fmt.Errorf("unknown proposal setting type [%v]", settingType)
		}
		results[i] = encodedArg
	}

	return results, nil
}
