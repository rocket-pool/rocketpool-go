package trustednode

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	quorumSettingPath                 = "members.quorum"
	rplBondSettingPath                = "members.rplbond"
	minipoolUnbondedMaxSettingPath    = "members.minipool.unbonded.max"
	minipoolUnbondedMinFeeSettingPath = "members.minipool.unbonded.min.fee"
	challengeCooldownSettingPath      = "members.challenge.cooldown"
	challengeWindowSettingPath        = "members.challenge.window"
	challengeCostSettingPath          = "members.challenge.cost"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAONodeTrustedSettingsMembers
type DaoNodeTrustedSettingsMembers struct {
	Details                         DaoNodeTrustedSettingsMembersDetails
	rp                              *rocketpool.RocketPool
	contract                        *core.Contract
	daoNodeTrustedContract          *trustednode.DaoNodeTrusted
	daoNodeTrustedProposalsContract *trustednode.DaoNodeTrustedProposals
}

// Details for RocketDAONodeTrustedSettingsMembers
type DaoNodeTrustedSettingsMembersDetails struct {
	Quorum                 core.Parameter[float64] `json:"quorum"`
	RplBond                *big.Int                `json:"rplBond"`
	UnbondedMinipoolMax    core.Parameter[uint64]  `json:"unbondedMinipoolMax"`
	UnbondedMinipoolMinFee core.Parameter[float64] `json:"unbondedMinipoolMinFee"`
	ChallengeCooldown      core.Parameter[uint64]  `json:"challengeCooldown"`
	ChallengeWindow        core.Parameter[uint64]  `json:"challengeWindow"`
	ChallengeCost          *big.Int                `json:"challengeCost"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoNodeTrustedSettingsMembers contract binding
func NewDaoNodeTrustedSettingsMembers(rp *rocketpool.RocketPool, daoNodeTrustedContract *trustednode.DaoNodeTrusted, daoNodeTrustedProposalsContract *trustednode.DaoNodeTrustedProposals, opts *bind.CallOpts) (*DaoNodeTrustedSettingsMembers, error) {
	// Create the contract
	contract, err := rp.GetContract(membersSettingsContractName, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted settings members contract: %w", err)
	}

	return &DaoNodeTrustedSettingsMembers{
		Details:                         DaoNodeTrustedSettingsMembersDetails{},
		rp:                              rp,
		contract:                        contract,
		daoNodeTrustedContract:          daoNodeTrustedContract,
		daoNodeTrustedProposalsContract: daoNodeTrustedProposalsContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the member proposal quorum threshold
func (c *DaoNodeTrustedSettingsMembers) GetQuorum(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.Quorum.RawValue, "getQuorum")
}

// Get the RPL bond required for a member
func (c *DaoNodeTrustedSettingsMembers) GetRplBond(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.RplBond, "getRPLBond")
}

// Get the maximum number of unbonded minipools a member can run
func (c *DaoNodeTrustedSettingsMembers) GetUnbondedMinipoolMax(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.UnbondedMinipoolMax.RawValue, "getMinipoolUnbondedMax")
}

// Get the minimum commission rate before unbonded minipools are allowed
func (c *DaoNodeTrustedSettingsMembers) GetUnbondedMinipoolMinFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.UnbondedMinipoolMinFee.RawValue, "getMinipoolUnbondedMinFee")
}

// Get the period a member must wait for before submitting another challenge, in blocks
func (c *DaoNodeTrustedSettingsMembers) GetChallengeCooldown(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ChallengeCooldown.RawValue, "getChallengeCooldown")
}

// Get the period during which a member can respond to a challenge, in blocks
func (c *DaoNodeTrustedSettingsMembers) GetChallengeWindow(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ChallengeWindow.RawValue, "getChallengeWindow")
}

// Get the fee for a non-member to challenge a member, in wei
func (c *DaoNodeTrustedSettingsMembers) GetChallengeCost(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ChallengeCost, "getChallengeCost")
}

// Get all basic details
func (c *DaoNodeTrustedSettingsMembers) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetQuorum(mc)
	c.GetRplBond(mc)
	c.GetUnbondedMinipoolMax(mc)
	c.GetUnbondedMinipoolMinFee(mc)
	c.GetChallengeCooldown(mc)
	c.GetChallengeWindow(mc)
	c.GetChallengeCost(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for setting the member proposal quorum threshold
func (c *DaoNodeTrustedSettingsMembers) BootstrapQuorum(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(membersSettingsContractName, quorumSettingPath, eth.EthToWei(value), opts)
}

// Get info for setting the RPL bond required for a member
func (c *DaoNodeTrustedSettingsMembers) BootstrapRplBond(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(membersSettingsContractName, rplBondSettingPath, value, opts)
}

// Get info for setting the maximum number of unbonded minipools a member can run
func (c *DaoNodeTrustedSettingsMembers) BootstrapUnbondedMinipoolMax(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(membersSettingsContractName, minipoolUnbondedMaxSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the minimum commission rate before unbonded minipools are allowed
func (c *DaoNodeTrustedSettingsMembers) BootstrapUnbondedMinipoolMinFee(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(membersSettingsContractName, minipoolUnbondedMinFeeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the period a member must wait for before submitting another challenge, in blocks
func (c *DaoNodeTrustedSettingsMembers) BootstrapChallengeCooldown(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(membersSettingsContractName, challengeCooldownSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the period during which a member can respond to a challenge, in blocks
func (c *DaoNodeTrustedSettingsMembers) BootstrapChallengeWindow(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(membersSettingsContractName, challengeWindowSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the fee for a non-member to challenge a member, in wei
func (c *DaoNodeTrustedSettingsMembers) BootstrapChallengeCost(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedContract.BootstrapUint(membersSettingsContractName, challengeCostSettingPath, value, opts)
}

// Get info for setting the member proposal quorum threshold
func (c *DaoNodeTrustedSettingsMembers) ProposeQuorum(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", quorumSettingPath), membersSettingsContractName, quorumSettingPath, eth.EthToWei(value), opts)
}

// Get info for setting the RPL bond required for a member
func (c *DaoNodeTrustedSettingsMembers) ProposeRplBond(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", rplBondSettingPath), membersSettingsContractName, rplBondSettingPath, value, opts)
}

// Get info for setting the maximum number of unbonded minipools a member can run
func (c *DaoNodeTrustedSettingsMembers) ProposeUnbondedMinipoolMax(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", minipoolUnbondedMaxSettingPath), membersSettingsContractName, minipoolUnbondedMaxSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the minimum commission rate before unbonded minipools are allowed
func (c *DaoNodeTrustedSettingsMembers) ProposeUnbondedMinipoolMinFee(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", minipoolUnbondedMinFeeSettingPath), membersSettingsContractName, minipoolUnbondedMinFeeSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the period a member must wait for before submitting another challenge, in blocks
func (c *DaoNodeTrustedSettingsMembers) ProposeChallengeCooldown(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", challengeCooldownSettingPath), membersSettingsContractName, challengeCooldownSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the period during which a member can respond to a challenge, in blocks
func (c *DaoNodeTrustedSettingsMembers) ProposeChallengeWindow(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", challengeWindowSettingPath), membersSettingsContractName, challengeWindowSettingPath, big.NewInt(int64(value)), opts)
}

// Get info for setting the fee for a non-member to challenge a member, in wei
func (c *DaoNodeTrustedSettingsMembers) ProposeChallengeCost(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoNodeTrustedProposalsContract.ProposeSetUint(fmt.Sprintf("set %s", challengeCostSettingPath), membersSettingsContractName, challengeCostSettingPath, value, opts)
}
