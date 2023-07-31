package trustednode

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/strings"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAONodeTrustedProposals
type DaoNodeTrustedProposals struct {
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoNodeTrustedProposals contract binding
func NewDaoNodeTrustedProposals(rp *rocketpool.RocketPool) (*DaoNodeTrustedProposals, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrustedProposals)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted proposals contract: %w", err)
	}

	return &DaoNodeTrustedProposals{
		rp:       rp,
		contract: contract,
	}, nil
}

// ====================
// === Transactions ===
// ====================

// Get info for proposing to invite a new member to the Oracle DAO
func (c *DaoNodeTrustedProposals) ProposeInviteMember(message string, newMemberAddress common.Address, newMemberId, string, newMemberUrl string, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	newMemberUrl = strings.Sanitize(newMemberUrl)
	return c.submitProposal(opts, message, "proposalInvite", newMemberId, newMemberUrl, newMemberAddress)
}

// Get info for proposing to leave the Oracle DAO
func (c *DaoNodeTrustedProposals) ProposeMemberLeave(message string, memberAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.submitProposal(opts, message, "proposalLeave", memberAddress)
}

// Get info for proposing to replace the address of an Oracle DAO member
func (c *DaoNodeTrustedProposals) ProposeReplaceMember(message string, memberAddress common.Address, newMemberAddress common.Address, newMemberId string, newMemberUrl string, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	newMemberUrl = strings.Sanitize(newMemberUrl)
	return c.submitProposal(opts, message, "proposalReplace", memberAddress, newMemberId, newMemberUrl, newMemberAddress)
}

// Get info for proposing to kick a member from the Oracle DAO
func (c *DaoNodeTrustedProposals) ProposeKickMember(message string, memberAddress common.Address, rplFineAmount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.submitProposal(opts, message, "proposalKick", memberAddress, rplFineAmount)
}

// Get info for proposing a bool setting
func (c *DaoNodeTrustedProposals) ProposeSetBool(message string, contractName string, settingPath string, value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.submitProposal(opts, message, "proposalSettingBool", contractName, settingPath, value)
}

// Get info for proposing a uint setting
func (c *DaoNodeTrustedProposals) ProposeSetUint(message string, contractName string, settingPath string, value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.submitProposal(opts, message, "proposalSettingUint", contractName, settingPath, value)
}

// Get info for proposing a contract upgrade
func (c *DaoNodeTrustedProposals) ProposeUpgradeContract(message string, upgradeType string, contractName string, contractAbi string, contractAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	compressedAbi, err := core.EncodeAbiStr(contractAbi)
	if err != nil {
		return nil, fmt.Errorf("error compressing ABI: %w", err)
	}
	return c.submitProposal(opts, message, "proposalUpgrade", upgradeType, contractName, compressedAbi, contractAddress)
}

// Get info for cancelling a proposal
func (c *DaoNodeTrustedProposals) CancelProposal(proposalId uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "cancel", opts, big.NewInt(int64(proposalId)))
}

// Get info for voting on a proposal
func (c *DaoNodeTrustedProposals) VoteOnProposal(proposalId uint64, support bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "vote", opts, big.NewInt(int64(proposalId)), support)
}

// Get info for executing a proposal
func (c *DaoNodeTrustedProposals) ExecuteProposal(proposalId uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "execute", opts, big.NewInt(int64(proposalId)))
}

// Internal method used for actually constructing and submitting a proposal
func (c *DaoNodeTrustedProposals) submitProposal(opts *bind.TransactOpts, message string, method string, args ...interface{}) (*core.TransactionInfo, error) {
	payload, err := c.contract.ABI.Pack(method, args...)
	if err != nil {
		return nil, fmt.Errorf("error encoding payload: %w", err)
	}
	return core.NewTransactionInfo(c.contract, "propose", opts, message, payload)
}
