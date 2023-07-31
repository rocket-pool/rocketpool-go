package trustednode

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	cooldownTimeSettingPath  = "proposal.cooldown.time"
	voteTimeSettingPath      = "proposal.vote.time"
	voteDelayTimeSettingPath = "proposal.vote.delay.time"
	executeTimeSettingPath   = "proposal.execute.time"
	actionTimeSettingPath    = "proposal.action.time"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAONodeTrustedSettingsProposals
type DaoNodeTrustedSettingsProposals struct {
	Details                         DaoNodeTrustedSettingsProposalsDetails
	rp                              *rocketpool.RocketPool
	contract                        *core.Contract
	daoNodeTrustedContract          *trustednode.DaoNodeTrusted
	daoNodeTrustedProposalsContract *trustednode.DaoNodeTrustedProposals
}

// Details for RocketDAONodeTrustedSettingsProposals
type DaoNodeTrustedSettingsProposalsDetails struct {
	CooldownTime  core.Parameter[time.Duration] `json:"cooldownTime"`
	VoteTime      core.Parameter[time.Duration] `json:"voteTime"`
	VoteDelayTime core.Parameter[time.Duration] `json:"voteDelayTime"`
	ExecuteTime   core.Parameter[time.Duration] `json:"executeTime"`
	ActionTime    core.Parameter[time.Duration] `json:"actionTime"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoNodeTrustedSettingsProposals contract binding
func NewDaoNodeTrustedSettingsProposals(rp *rocketpool.RocketPool) (*DaoNodeTrustedSettingsProposals, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted settings proposals contract: %w", err)
	}
	daoNodeTrustedContract, err := trustednode.NewDaoNodeTrusted(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted contract: %w", err)
	}
	daoNodeTrustedProposalsContract, err := trustednode.NewDaoNodeTrustedProposals(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted proposals contract: %w", err)
	}

	return &DaoNodeTrustedSettingsProposals{
		Details:                         DaoNodeTrustedSettingsProposalsDetails{},
		rp:                              rp,
		contract:                        contract,
		daoNodeTrustedContract:          daoNodeTrustedContract,
		daoNodeTrustedProposalsContract: daoNodeTrustedProposalsContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the cooldown period a member must wait, in seconds, after making a proposal before making another
func (c *DaoNodeTrustedSettingsProposals) GetCooldownTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.CooldownTime.RawValue, "getCooldownTime")
}

// Get the period, in seconds, a proposal can be voted on
func (c *DaoNodeTrustedSettingsProposals) GetVoteTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.VoteTime.RawValue, "getVoteTime")
}

// Get the delay, in seconds, after creation before a proposal can be voted on
func (c *DaoNodeTrustedSettingsProposals) GetVoteDelayTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.VoteDelayTime.RawValue, "getVoteDelayTime")
}

// Get the period, in seconds, during which a passed proposal can be executed
func (c *DaoNodeTrustedSettingsProposals) GetExecuteTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ExecuteTime.RawValue, "getExecuteTime")
}

// Get the period, in seconds, during which an action can be performed on an executed proposal
func (c *DaoNodeTrustedSettingsProposals) GetActionTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ActionTime.RawValue, "getActionTime")
}

// Get all basic details
func (c *DaoNodeTrustedSettingsProposals) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetCooldownTime(mc)
	c.GetVoteTime(mc)
	c.GetVoteDelayTime(mc)
	c.GetExecuteTime(mc)
	c.GetActionTime(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for setting the cooldown period a member must wait, in seconds, after making a proposal before making another
func (c *DaoNodeTrustedSettingsProposals) BootstrapProposalCooldownTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, cooldownTimeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the period, in seconds, a proposal can be voted on
func (c *DaoNodeTrustedSettingsProposals) BootstrapProposalVoteTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, voteTimeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the delay, in seconds, after creation before a proposal can be voted on
func (c *DaoNodeTrustedSettingsProposals) BootstrapProposalVoteDelayTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, voteDelayTimeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the period, in seconds, during which a passed proposal can be executed
func (c *DaoNodeTrustedSettingsProposals) BootstrapProposalExecuteTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, executeTimeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the period, in seconds, during which an action can be performed on an executed proposal
func (c *DaoNodeTrustedSettingsProposals) BootstrapProposalActionTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, actionTimeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the cooldown period a member must wait, in seconds, after making a proposal before making another
func (c *DaoNodeTrustedSettingsProposals) ProposeProposalCooldownTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", cooldownTimeSettingPath), rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, cooldownTimeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the period, in seconds, a proposal can be voted on
func (c *DaoNodeTrustedSettingsProposals) ProposeProposalVoteTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", voteTimeSettingPath), rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, voteTimeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the delay, in seconds, after creation before a proposal can be voted on
func (c *DaoNodeTrustedSettingsProposals) ProposeProposalVoteDelayTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", voteDelayTimeSettingPath), rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, voteDelayTimeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the period, in seconds, during which a passed proposal can be executed
func (c *DaoNodeTrustedSettingsProposals) ProposeProposalExecuteTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", executeTimeSettingPath), rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, executeTimeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the period, in seconds, during which an action can be performed on an executed proposal
func (c *DaoNodeTrustedSettingsProposals) ProposeProposalActionTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", actionTimeSettingPath), rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, actionTimeSettingPath, big.NewInt(int64(value)), opts)
}
