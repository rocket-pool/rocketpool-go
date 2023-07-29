package node

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for a Rocket Pool Node
type Node struct {
	Details NodeDetails
	mgr     *NodeManager
	staking *NodeStaking
}

// Details for a Rocket Pool Node
type NodeDetails struct {
	// Primitives
	Index   uint64         `json:"index"`
	Address common.Address `json:"address"`

	// NodeManager
	Exists                           bool                      `json:"exists"`
	RegistrationTime                 core.Parameter[time.Time] `json:"registrationTime"`
	TimezoneLocation                 string                    `json:"timezoneLocation"`
	RewardNetwork                    core.Parameter[uint64]    `json:"rewardNetwork"`
	IsFeeDistributorInitialized      bool                      `json:"isFeeDistributorInitialized"`
	AverageFee                       core.Parameter[float64]   `json:"averageFee"`
	SmoothingPoolRegistrationState   bool                      `json:"smoothingPoolRegistrationState"`
	SmoothingPoolRegistrationChanged core.Parameter[time.Time] `json:"smoothingPoolRegistrationChanged"`

	// NodeStaking
	RplStake          *big.Int                  `json:"rplStake"`
	EffectiveRplStake *big.Int                  `json:"effectiveRplStake"`
	MinimumRplStake   *big.Int                  `json:"minimumRplStake"`
	MaximumRplStake   *big.Int                  `json:"maximumRplStake"`
	RplStakedTime     core.Parameter[time.Time] `json:"rplStakedTime"`
	EthMatched        *big.Int                  `json:"ethMatched"`
	EthMatchedLimit   *big.Int                  `json:"ethMatchedLimit"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new Node instance
func NewNode(mgr *NodeManager, staking *NodeStaking, index uint64, address common.Address) *Node {
	return &Node{
		Details: NodeDetails{
			Index:   index,
			Address: address,
		},
		mgr:     mgr,
		staking: staking,
	}
}

// =============
// === Calls ===
// =============

// Check whether or not the node exists
func (c *Node) GetExists(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.Exists, "getNodeExists", c.Details.Address)
}

// Get the time that the user registered
func (c *Node) GetRegistrationTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.RegistrationTime.RawValue, "getNodeRegistrationTime", c.Details.Address)
}

// Get the node's timezone location
func (c *Node) GetTimezoneLocation(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.TimezoneLocation, "getNodeTimezoneLocation", c.Details.Address)
}

// Get the network ID for the node's rewards
func (c *Node) GetRewardNetwork(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.RewardNetwork.RawValue, "getRewardNetwork", c.Details.Address)
}

// Check if the node's fee distributor has been initialized yet
func (c *Node) GetFeeDistributorInitialized(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.IsFeeDistributorInitialized, "getFeeDistributorInitialised", c.Details.Address)
}

// Get a node's average minipool fee (commission)
func (c *Node) GetAverageFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.AverageFee.RawValue, "getAverageNodeFee", c.Details.Address)
}

// Get the node's smoothing pool opt-in status
func (c *Node) GetSmoothingPoolRegistrationState(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.SmoothingPoolRegistrationState, "getSmoothingPoolRegistrationState", c.Details.Address)
}

// Get the time of the node's last smoothing pool opt-in / opt-out
func (c *Node) GetSmoothingPoolRegistrationChanged(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.contract, &c.Details.SmoothingPoolRegistrationChanged.RawValue, "getSmoothingPoolRegistrationChanged", c.Details.Address)
}

// Get the node's RPL stake
func (c *Node) GetRplStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.staking.contract, &c.Details.RplStake, "getNodeRPLStake", c.Details.Address)
}

// Get the node's effective RPL stake
func (c *Node) GetEffectiveRplStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.staking.contract, &c.Details.EffectiveRplStake, "getNodeEffectiveRPLStake", c.Details.Address)
}

// Get the node's minimum RPL stake to collateralize its minipools
func (c *Node) GetMinimumRplStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.staking.contract, &c.Details.MinimumRplStake, "getNodeMinimumRPLStake", c.Details.Address)
}

// Get the node's maximum RPL stake to collateralize its minipools
func (c *Node) GetMaximumRplStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.staking.contract, &c.Details.MaximumRplStake, "getNodeMaximumRPLStake", c.Details.Address)
}

// Get the time the node last staked RPL
func (c *Node) GetRplStakedTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.staking.contract, &c.Details.RplStakedTime.RawValue, "getNodeRPLStakedTime", c.Details.Address)
}

// Get the amount of ETH the node has borrowed from the deposit pool to create its minipools
func (c *Node) GetEthMatched(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.staking.contract, &c.Details.EthMatched, "getNodeETHMatched", c.Details.Address)
}

// Get the amount of ETH the node can still borrow from the deposit pool to create any new minipools
func (c *Node) GetEthMatchedLimit(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.staking.contract, &c.Details.EthMatchedLimit, "getNodeETHMatchedLimit", c.Details.Address)
}

// Get all basic details
func (c *Node) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetExists(mc)
	c.GetRegistrationTime(mc)
	c.GetTimezoneLocation(mc)
	c.GetRewardNetwork(mc)
	c.GetFeeDistributorInitialized(mc)
	c.GetAverageFee(mc)
	c.GetSmoothingPoolRegistrationState(mc)
	c.GetSmoothingPoolRegistrationChanged(mc)
	c.GetRplStake(mc)
	c.GetEffectiveRplStake(mc)
	c.GetMinimumRplStake(mc)
	c.GetMaximumRplStake(mc)
	c.GetRplStakedTime(mc)
	c.GetEthMatched(mc)
	c.GetEthMatchedLimit(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for registering a node
func (c *Node) Register(timezoneLocation string, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	_, err := time.LoadLocation(timezoneLocation)
	if err != nil {
		return nil, fmt.Errorf("error verifying timezone [%s]: %w", timezoneLocation, err)
	}
	return core.NewTransactionInfo(c.mgr.contract, "registerNode", opts, timezoneLocation)
}

// Get info for setting a node's timezone location
func (c *Node) SetTimezoneLocation(timezoneLocation string, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	_, err := time.LoadLocation(timezoneLocation)
	if err != nil {
		return nil, fmt.Errorf("error verifying timezone [%s]: %w", timezoneLocation, err)
	}
	return core.NewTransactionInfo(c.mgr.contract, "setTimezoneLocation", opts, timezoneLocation)
}

// Get info for initializing (creating) the node's fee distributor
func (c *Node) InitializeFeeDistributor(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.mgr.contract, "initialiseFeeDistributor", opts)
}

// Get info for opting in or out of the smoothing pool
func (c *Node) SetSmoothingPoolRegistrationState(optIn bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.mgr.contract, "setSmoothingPoolRegistrationState", opts, optIn)
}

// Get info for staking RPL
func (c *Node) StakeRpl(rplAmount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.mgr.contract, "stakeRPL", opts, rplAmount)
}

// Get info for adding or removing an address from the stake-RPL-on-behalf allowlist
func (c *Node) SetStakeRplForAllowed(caller common.Address, allowed bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.mgr.contract, "setStakeRPLForAllowed", opts, caller, allowed)
}

// Get info for withdrawing staked RPL
func (c *Node) WithdrawRpl(rplAmount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.mgr.contract, "withdrawRPL", opts, rplAmount)
}
