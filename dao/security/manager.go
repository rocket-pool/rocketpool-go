package security

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

const (
	securityCouncilMemberBatchSize int = 200
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOSecurity
type SecurityCouncilManager struct {
	// The number of members in the security council
	MemberCount *core.FormattedUint256Field[uint64]

	// The total amount of votesneeded for a proposal to pass
	MemberQuorumVotesRequired *core.FormattedUint256Field[float64]

	// Settings for the Protocol DAO
	Settings *SecurityCouncilSettings

	// === Internal fields ===
	rp    *rocketpool.RocketPool
	ds    *core.Contract
	dsa   *core.Contract
	dsp   *core.Contract
	txMgr *eth.TransactionManager
}

// ====================
// === Constructors ===
// ====================

// Creates a new SecurityCouncilManager contract binding
func NewSecurityCouncilManager(rp *rocketpool.RocketPool, pSettings *protocol.ProtocolDaoSettings) (*SecurityCouncilManager, error) {
	// Create the contracts
	ds, err := rp.GetContract(rocketpool.ContractName_RocketDAOSecurity)
	if err != nil {
		return nil, fmt.Errorf("error getting security council manager contract: %w", err)
	}
	dsa, err := rp.GetContract(rocketpool.ContractName_RocketDAOSecurityActions)
	if err != nil {
		return nil, fmt.Errorf("error getting security council actions contract: %w", err)
	}
	dsp, err := rp.GetContract(rocketpool.ContractName_RocketDAOSecurityProposals)
	if err != nil {
		return nil, fmt.Errorf("error getting security council proposals contract: %w", err)
	}

	secMgr := &SecurityCouncilManager{
		MemberCount:               core.NewFormattedUint256Field[uint64](ds, "getMemberCount"),
		MemberQuorumVotesRequired: core.NewFormattedUint256Field[float64](ds, "getMemberQuorumVotesRequired"),

		rp:    rp,
		ds:    ds,
		dsa:   dsa,
		dsp:   dsp,
		txMgr: rp.GetTransactionManager(),
	}
	settings, err := newSecurityCouncilSettings(secMgr, pSettings)
	if err != nil {
		return nil, fmt.Errorf("error creating security council settings binding: %w", err)
	}
	secMgr.Settings = settings
	return secMgr, nil
}

// =============
// === Calls ===
// =============

// ====================
// === Transactions ===
// ====================

// === DAOSecurityActions ===

// Get info for joining the security council
func (c *SecurityCouncilManager) Join(opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.dsa.Contract, "actionJoin", opts)
}

// Get info for removing a member from the security council
func (c *SecurityCouncilManager) Kick(address common.Address, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.dsa.Contract, "actionKick", opts, address)
}

// Get info for removing multiple members from the security council
func (c *SecurityCouncilManager) KickMulti(addresses []common.Address, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.dsa.Contract, "actionKickMulti", opts, addresses)
}

// Get info for requesting to leave the security council
func (c *SecurityCouncilManager) RequestLeave(opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.dsa.Contract, "actionRequestLeave", opts)
}

// Get info for leaving the security council
func (c *SecurityCouncilManager) Leave(opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.dsa.Contract, "actionLeave", opts)
}

// === DAOSecurityProposals ===

// Get info for proposing a uint setting
func (c *SecurityCouncilManager) ProposeSetUint(message string, contractName rocketpool.ContractName, setting protocol.SettingName, value *big.Int, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("set %s", setting)
	}
	return c.submitProposal(opts, message, "proposalSettingUint", contractName, string(setting), value)
}

// Get info for proposing a bool setting
func (c *SecurityCouncilManager) ProposeSetBool(message string, contractName rocketpool.ContractName, setting protocol.SettingName, value bool, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("set %s", setting)
	}
	return c.submitProposal(opts, message, "proposalSettingBool", contractName, string(setting), value)
}

// Get info for proposing to invite a new member to the security council
func (c *SecurityCouncilManager) ProposeInvite(message string, newMemberID string, newMemberAddress common.Address, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("invite %s (%s)", newMemberID, newMemberAddress.Hex())
	}
	return c.submitProposal(opts, message, "proposalInvite", newMemberID, newMemberAddress)
}

// Get info for proposing to kick a member from the security council
func (c *SecurityCouncilManager) ProposeKick(message string, memberAddress common.Address, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("kick %s", memberAddress.Hex())
	}
	return c.submitProposal(opts, message, "proposalKick", memberAddress)
}

// Get info for proposing to kick multiple members from the security council
func (c *SecurityCouncilManager) ProposeKickMulti(message string, memberAddresses []common.Address, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	if message == "" {
		message = "kick multiple members"
	}
	return c.submitProposal(opts, message, "proposalKick", memberAddresses)
}

// Get info for proposing to kick a member from the security council and replace it with a new member
func (c *SecurityCouncilManager) ProposeReplace(message string, existingMemberAddress common.Address, newMemberID string, newMemberAddress common.Address, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	if message == "" {
		message = fmt.Sprintf("replace %s with %s (%s)", existingMemberAddress.Hex(), newMemberID, newMemberAddress.Hex())
	}
	return c.submitProposal(opts, message, "proposalReplace", existingMemberAddress, newMemberID, newMemberAddress)
}

// Internal method used for actually constructing and submitting a proposal
func (c *SecurityCouncilManager) submitProposal(opts *bind.TransactOpts, message string, method string, args ...interface{}) (*eth.TransactionInfo, error) {
	payload, err := c.dsp.ABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("error encoding payload: %w", err)
	}
	return c.txMgr.CreateTransactionInfo(c.dsp.Contract, "propose", opts, message, payload)
}

// =================
// === Addresses ===
// =================

// Get a security council member address by index
func (c *SecurityCouncilManager) GetMemberAddress(mc *batch.MultiCaller, address_Out *common.Address, index uint64) {
	core.AddCall(mc, c.ds, address_Out, "getMemberAt", big.NewInt(int64(index)))
}

// Get the list of security council member addresses.
func (c *SecurityCouncilManager) GetMemberAddresses(memberCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, memberCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(memberCount), c.rp.AddressBatchSize,
		func(mc *batch.MultiCaller, index int) error {
			c.GetMemberAddress(mc, &addresses[index], uint64(index))
			return nil
		},
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting security council member addresses: %w", err)
	}

	// Return
	return addresses, nil
}

// Get a security council member by address.
func (c *SecurityCouncilManager) CreateMemberFromAddress(address common.Address, includeDetails bool, opts *bind.CallOpts) (*SecurityCouncilMember, error) {
	// Create the member binding
	member, err := NewSecurityCouncilMember(c.rp, address)
	if err != nil {
		return nil, fmt.Errorf("error creating security council member binding for %s: %w", address.Hex(), err)
	}

	if includeDetails {
		err = c.rp.Query(func(mc *batch.MultiCaller) error {
			eth.QueryAllFields(member, mc)
			return nil
		}, opts)
		if err != nil {
			return nil, fmt.Errorf("error getting security council member details: %w", err)
		}
	}

	// Return
	return member, nil
}

// Get the list of all security council members.
func (c *SecurityCouncilManager) CreateMembersFromAddresses(addresses []common.Address, includeDetails bool, opts *bind.CallOpts) ([]*SecurityCouncilMember, error) {
	// Create the member bindings
	memberCount := len(addresses)
	members := make([]*SecurityCouncilMember, memberCount)
	for i, address := range addresses {
		member, err := NewSecurityCouncilMember(c.rp, address)
		if err != nil {
			return nil, fmt.Errorf("error creating security council member binding for %s: %w", address.Hex(), err)
		}
		members[i] = member
	}

	if includeDetails {
		err := c.rp.BatchQuery(int(memberCount), securityCouncilMemberBatchSize,
			func(mc *batch.MultiCaller, index int) error {
				member := members[index]
				eth.QueryAllFields(member, mc)
				return nil
			},
			opts,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting security council member details: %w", err)
		}
	}

	// Return
	return members, nil
}
