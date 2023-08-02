package node

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for a Rocket Pool Node
type Node struct {
	Details     NodeDetails
	rp          *rocketpool.RocketPool
	nodeMgr     *core.Contract
	nodeStaking *core.Contract
	mpFactory   *core.Contract
	mpMgr       *core.Contract
}

// Details for a Rocket Pool Node
type NodeDetails struct {
	// NodeManager
	Address                          common.Address            `json:"address"`
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

	// MinipoolManager
	MinipoolCount           core.Parameter[uint64] `json:"minipoolCount"`
	ActiveMinipoolCount     core.Parameter[uint64] `json:"activeMinipoolCount"`
	FinalisedMinipoolCount  core.Parameter[uint64] `json:"finalisedMinipoolCount"`
	ValidatingMinipoolCount core.Parameter[uint64] `json:"validatingMinipoolCount"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new Node instance
func NewNode(rp *rocketpool.RocketPool, address common.Address) (*Node, error) {
	nodeManager, err := rp.GetContract(rocketpool.ContractName_RocketNodeManager)
	if err != nil {
		return nil, fmt.Errorf("error getting node staking binding: %w", err)
	}
	nodeStaking, err := rp.GetContract(rocketpool.ContractName_RocketNodeStaking)
	if err != nil {
		return nil, fmt.Errorf("error getting node staking binding: %w", err)
	}
	minipoolFactory, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolFactory)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool factory binding: %w", err)
	}
	minipoolManager, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolManager)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool manager binding: %w", err)
	}

	return &Node{
		Details: NodeDetails{
			Address: address,
		},
		rp:          rp,
		nodeMgr:     nodeManager,
		nodeStaking: nodeStaking,
		mpFactory:   minipoolFactory,
		mpMgr:       minipoolManager,
	}, nil
}

// =============
// === Calls ===
// =============

// === NodeManager ===

// Check whether or not the node exists
func (c *Node) GetExists(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeMgr, &c.Details.Exists, "getNodeExists", c.Details.Address)
}

// Get the time that the user registered
func (c *Node) GetRegistrationTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeMgr, &c.Details.RegistrationTime.RawValue, "getNodeRegistrationTime", c.Details.Address)
}

// Get the node's timezone location
func (c *Node) GetTimezoneLocation(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeMgr, &c.Details.TimezoneLocation, "getNodeTimezoneLocation", c.Details.Address)
}

// Get the network ID for the node's rewards
func (c *Node) GetRewardNetwork(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeMgr, &c.Details.RewardNetwork.RawValue, "getRewardNetwork", c.Details.Address)
}

// Check if the node's fee distributor has been initialized yet
func (c *Node) GetFeeDistributorInitialized(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeMgr, &c.Details.IsFeeDistributorInitialized, "getFeeDistributorInitialised", c.Details.Address)
}

// Get a node's average minipool fee (commission)
func (c *Node) GetAverageFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeMgr, &c.Details.AverageFee.RawValue, "getAverageNodeFee", c.Details.Address)
}

// Get the node's smoothing pool opt-in status
func (c *Node) GetSmoothingPoolRegistrationState(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeMgr, &c.Details.SmoothingPoolRegistrationState, "getSmoothingPoolRegistrationState", c.Details.Address)
}

// Get the time of the node's last smoothing pool opt-in / opt-out
func (c *Node) GetSmoothingPoolRegistrationChanged(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeMgr, &c.Details.SmoothingPoolRegistrationChanged.RawValue, "getSmoothingPoolRegistrationChanged", c.Details.Address)
}

// === NodeStaking ===

// Get the node's RPL stake
func (c *Node) GetRplStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeStaking, &c.Details.RplStake, "getNodeRPLStake", c.Details.Address)
}

// Get the node's effective RPL stake
func (c *Node) GetEffectiveRplStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeStaking, &c.Details.EffectiveRplStake, "getNodeEffectiveRPLStake", c.Details.Address)
}

// Get the node's minimum RPL stake to collateralize its minipools
func (c *Node) GetMinimumRplStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeStaking, &c.Details.MinimumRplStake, "getNodeMinimumRPLStake", c.Details.Address)
}

// Get the node's maximum RPL stake to collateralize its minipools
func (c *Node) GetMaximumRplStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeStaking, &c.Details.MaximumRplStake, "getNodeMaximumRPLStake", c.Details.Address)
}

// Get the time the node last staked RPL
func (c *Node) GetRplStakedTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeStaking, &c.Details.RplStakedTime.RawValue, "getNodeRPLStakedTime", c.Details.Address)
}

// Get the amount of ETH the node has borrowed from the deposit pool to create its minipools
func (c *Node) GetEthMatched(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeStaking, &c.Details.EthMatched, "getNodeETHMatched", c.Details.Address)
}

// Get the amount of ETH the node can still borrow from the deposit pool to create any new minipools
func (c *Node) GetEthMatchedLimit(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.nodeStaking, &c.Details.EthMatchedLimit, "getNodeETHMatchedLimit", c.Details.Address)
}

// === MinipoolManager ===

// Get all basic details
func (c *Node) GetBasicDetails(mc *multicall.MultiCaller) {
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

// === MinipoolManager ===

// Get the node's minipool count
func (c *Node) GetMinipoolCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mpMgr, &c.Details.MinipoolCount.RawValue, "getNodeMinipoolCount", c.Details.Address)
}

// Get the number of minipools owned by a node that are not finalised
func (c *Node) GetActiveMinipoolCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mpMgr, &c.Details.ActiveMinipoolCount.RawValue, "getNodeActiveMinipoolCount", c.Details.Address)
}

// Get the number of minipools owned by a node that are finalised
func (c *Node) GetFinalisedMinipoolCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mpMgr, &c.Details.FinalisedMinipoolCount.RawValue, "getNodeFinalisedMinipoolCount", c.Details.Address)
}

// Get the number of minipools owned by a node that are validating
func (c *Node) GetValidatingMinipoolCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mpMgr, &c.Details.ValidatingMinipoolCount.RawValue, "getNodeValidatingMinipoolCount", c.Details.Address)
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
	return core.NewTransactionInfo(c.nodeMgr, "registerNode", opts, timezoneLocation)
}

// Get info for setting a node's timezone location
func (c *Node) SetTimezoneLocation(timezoneLocation string, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	_, err := time.LoadLocation(timezoneLocation)
	if err != nil {
		return nil, fmt.Errorf("error verifying timezone [%s]: %w", timezoneLocation, err)
	}
	return core.NewTransactionInfo(c.nodeMgr, "setTimezoneLocation", opts, timezoneLocation)
}

// Get info for initializing (creating) the node's fee distributor
func (c *Node) InitializeFeeDistributor(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nodeMgr, "initialiseFeeDistributor", opts)
}

// Get info for opting in or out of the smoothing pool
func (c *Node) SetSmoothingPoolRegistrationState(optIn bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nodeMgr, "setSmoothingPoolRegistrationState", opts, optIn)
}

// Get info for staking RPL
func (c *Node) StakeRpl(rplAmount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nodeMgr, "stakeRPL", opts, rplAmount)
}

// Get info for adding or removing an address from the stake-RPL-on-behalf allowlist
func (c *Node) SetStakeRplForAllowed(caller common.Address, allowed bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nodeMgr, "setStakeRPLForAllowed", opts, caller, allowed)
}

// Get info for withdrawing staked RPL
func (c *Node) WithdrawRpl(rplAmount *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nodeMgr, "withdrawRPL", opts, rplAmount)
}

// ===================
// === Sub-Getters ===
// ===================

// === MinipoolManager ===

// Get one of the node's minipool addresses by index
func (c *Node) GetMinipoolAddress(mc *multicall.MultiCaller, address_Out *common.Address, index uint64) {
	multicall.AddCall(mc, c.mpMgr, address_Out, "getNodeMinipoolAt", c.Details.Address, big.NewInt(int64(index)))
}

// Get one of the node's validating minipool addresses by index
func (c *Node) GetValidatingMinipoolAddress(mc *multicall.MultiCaller, address_Out *common.Address, index uint64) {
	multicall.AddCall(mc, c.mpMgr, address_Out, "getNodeValidatingMinipoolAt", c.Details.Address, big.NewInt(int64(index)))
}

// Get all of the node's minipool addresses in a standalone call.
// This will use an internal batched multicall invocation to retrieve all of them.
// Provide the value returned from GetMinipoolCount() in minipoolCount.
func (c *Node) GetMinipoolAddresses(mc *multicall.MultiCaller, minipoolCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, minipoolCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(minipoolCount), c.rp.AddressBatchSize,
		func(mc *multicall.MultiCaller, index int) error {
			c.GetMinipoolAddress(mc, &addresses[index], uint64(index))
			return nil
		}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool addresses: %w", err)
	}

	// Return
	return addresses, nil
}

// Get all of the node's validating minipool addresses in a standalone call.
// This will use an internal batched multicall invocation to retrieve all of them.
// Provide the value returned from GetValidatingMinipoolCount() in minipoolCount.
func (c *Node) GetValidatingMinipoolAddresses(mc *multicall.MultiCaller, minipoolCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, minipoolCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(minipoolCount), c.rp.AddressBatchSize,
		func(mc *multicall.MultiCaller, index int) error {
			c.GetValidatingMinipoolAddress(mc, &addresses[index], uint64(index))
			return nil
		}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool addresses: %w", err)
	}

	// Return
	return addresses, nil
}

// =============
// === Utils ===
// =============

// === MinipoolFactory ===

// Get the address of a minipool based on the node's address and a salt
func (c *Node) GetExpectedMinipoolAddress(mc *multicall.MultiCaller, address_Out *common.Address, salt *big.Int) {
	multicall.AddCall(mc, c.mpFactory, address_Out, "getExpectedAddress", c.Details.Address, salt)
}
