package protocol

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
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

	// === Internal fields ===
	rp   *rocketpool.RocketPool
	dp   *core.Contract
	dpp  *core.Contract
	dpsr *core.Contract
}

// Rewards claimer percents
type RplRewardsPercentages struct {
	OdaoPercentage *big.Int `abi:"_trustedNodePercent"`
	PdaoPercentage *big.Int `abi:"_protocolPercent"`
	NodePercentage *big.Int `abi:"_nodePercent"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProtocolDaoManager contract binding
func NewProtocolDaoManager(rp *rocketpool.RocketPool) (*ProtocolDaoManager, error) {
	// Create the contracts
	dp, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocol)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO manager contract: %w", err)
	}
	dpp, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocolProposals)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO protocol proposals contract: %w", err)
	}
	dpsr, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocolSettingsRewards)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO protocol settings rewards contract: %w", err)
	}

	pdaoMgr := &ProtocolDaoManager{
		LastRewardsPercentagesUpdate: core.NewFormattedUint256Field[time.Time](dpsr, "getRewardsClaimersTimeUpdated"),

		rp:   rp,
		dp:   dp,
		dpp:  dpp,
		dpsr: dpsr,
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

// Get the allocation of RPL rewards to the node operators, Oracle DAO, and the Protocol DAO
func (c *ProtocolDaoManager) GetRewardsPercentages(mc *batch.MultiCaller, out *RplRewardsPercentages) {
	core.AddCallRaw(mc, c.dpsr, out, "getRewardsClaimersPerc")
}

// ====================
// === Transactions ===
// ====================

// === DAOProtocol ===

// Get info for bootstrapping a bool setting
func (c *ProtocolDaoManager) BootstrapBool(contractName rocketpool.ContractName, settingPath string, value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingBool", opts, contractName, settingPath, value)
}

// Get info for bootstrapping a uint256 setting
func (c *ProtocolDaoManager) BootstrapUint(contractName rocketpool.ContractName, settingPath string, value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingUint", opts, contractName, settingPath, value)
}

// Get info for bootstrapping an address setting
func (c *ProtocolDaoManager) BootstrapAddress(contractName rocketpool.ContractName, settingPath string, value common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingAddress", opts, contractName, settingPath, value)
}

// Get info for bootstrapping a rewards claimer
func (c *ProtocolDaoManager) BootstrapClaimer(contractName rocketpool.ContractName, amount float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dp, "bootstrapSettingClaimer", opts, contractName, eth.EthToWei(amount))
}

// === DAOProtocolProposals ===

// Get info for submitting a proposal to update a bool Protocol DAO setting
func (c *ProtocolDaoManager) ProposeSetBool(message string, contractName rocketpool.ContractName, settingPath string, value bool, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("set %s", settingPath)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSettingBool", contractName, settingPath, value)
}

// Get info for submitting a proposal to update a uint Protocol DAO setting
func (c *ProtocolDaoManager) ProposeSetUint(message string, contractName rocketpool.ContractName, settingPath string, value *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("set %s", settingPath)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSettingUint", contractName, settingPath, value)
}

// Get info for submitting a proposal to update an address Protocol DAO setting
func (c *ProtocolDaoManager) ProposeSetAddress(message string, contractName rocketpool.ContractName, settingPath string, value common.Address, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("set %s", settingPath)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSettingAddress", contractName, settingPath, value)
}

// Get info for submitting a proposal to update multiple Protocol DAO settings at once
func (c *ProtocolDaoManager) ProposeSetMulti(message string, contractNames []rocketpool.ContractName, settingPaths []string, settingTypes []types.ProposalSettingType, values []any, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("set %s", strings.Join(settingPaths, ", "))
	}
	encodedValues, err := abiEncodeMultiValues(settingTypes, values)
	if err != nil {
		return nil, fmt.Errorf("error ABI encoding values: %w", err)
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSettingMulti", contractNames, settingPaths, settingTypes, encodedValues)
}

// Get info for submitting a proposal to update the allocations of RPL rewards
func (c *ProtocolDaoManager) ProposeSetRewardsPercentages(message string, odaoPercentage *big.Int, pdaoPercentage *big.Int, nodePercentage *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	if message == "" {
		message = "set rewards percentages"
	}
	return c.submitProposal(opts, blockNumber, treeNodes, message, "proposalSettingRewardsClaimers", odaoPercentage, pdaoPercentage, nodePercentage)
}

// Get info for submitting a proposal to spend a portion of the Rocket Pool treasury one time
func (c *ProtocolDaoManager) ProposeOneTimeTreasurySpend(rp *rocketpool.RocketPool, message, invoiceID string, recipient common.Address, amount *big.Int, blockNumber uint32, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
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

// Submit a protocol DAO proposal
func (c *ProtocolDaoManager) submitProposal(opts *bind.TransactOpts, blockNumber uint32, treeNodes []types.VotingTreeNode, message string, method string, args ...interface{}) (*core.TransactionInfo, error) {
	payload, err := c.dpp.ABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("error encoding payload: %w", err)
	}
	return core.NewTransactionInfo(c.dpp, "propose", opts, message, payload, blockNumber, treeNodes)
}

/// =============
/// === Utils ===
/// =============

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
