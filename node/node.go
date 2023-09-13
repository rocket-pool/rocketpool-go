package node

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// ===============
// === Structs ===
// ===============

// Binding for a Rocket Pool Node
type Node struct {
	Details     NodeDetails
	rp          *rocketpool.RocketPool
	distFactory *core.Contract
	nodeDeposit *core.Contract
	nodeMgr     *core.Contract
	nodeStaking *core.Contract
	mpFactory   *core.Contract
	mpMgr       *core.Contract
	storage     *core.Contract
}

// Details for a Rocket Pool Node
type NodeDetails struct {
	// DistributorFactory
	DistributorAddress common.Address `json:"distributorAddress"`

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

	// NodeDeposit
	Credit *big.Int `json:"credit"`

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

	// Storage
	WithdrawalAddress        common.Address `json:"withdrawalAddress"`
	PendingWithdrawalAddress common.Address `json:"pendingWithdrawalAddress"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new Node instance
func NewNode(rp *rocketpool.RocketPool, address common.Address) (*Node, error) {
	distFactory, err := rp.GetContract(rocketpool.ContractName_RocketNodeDistributorFactory)
	if err != nil {
		return nil, fmt.Errorf("error getting distributor factory binding: %w", err)
	}
	nodeDeposit, err := rp.GetContract(rocketpool.ContractName_RocketNodeDeposit)
	if err != nil {
		return nil, fmt.Errorf("error getting node deposit binding: %w", err)
	}
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
		distFactory: distFactory,
		nodeDeposit: nodeDeposit,
		nodeMgr:     nodeManager,
		nodeStaking: nodeStaking,
		mpFactory:   minipoolFactory,
		mpMgr:       minipoolManager,
		storage:     rp.Storage.Contract,
	}, nil
}

// =============
// === Calls ===
// =============

// === DistributorFactory ===

// Get the node's fee distributor address
func (c *Node) GetDistributorAddress(mc *batch.MultiCaller) {
	core.AddCall(mc, c.distFactory, &c.Details.DistributorAddress, "getProxyAddress", c.Details.Address)
}

// === NodeDeposit ===

// Get the amount of ETH in the node's deposit credit bank
func (c *Node) GetDepositCredit(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeDeposit, &c.Details.Credit, "getNodeDepositCredit", c.Details.Address)
}

// === NodeManager ===

// Check whether or not the node exists
func (c *Node) GetExists(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeMgr, &c.Details.Exists, "getNodeExists", c.Details.Address)
}

// Get the time that the user registered
func (c *Node) GetRegistrationTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeMgr, &c.Details.RegistrationTime.RawValue, "getNodeRegistrationTime", c.Details.Address)
}

// Get the node's timezone location
func (c *Node) GetTimezoneLocation(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeMgr, &c.Details.TimezoneLocation, "getNodeTimezoneLocation", c.Details.Address)
}

// Get the network ID for the node's rewards
func (c *Node) GetRewardNetwork(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeMgr, &c.Details.RewardNetwork.RawValue, "getRewardNetwork", c.Details.Address)
}

// Check if the node's fee distributor has been initialized yet
func (c *Node) GetFeeDistributorInitialized(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeMgr, &c.Details.IsFeeDistributorInitialized, "getFeeDistributorInitialised", c.Details.Address)
}

// Get a node's average minipool fee (commission)
func (c *Node) GetAverageFee(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeMgr, &c.Details.AverageFee.RawValue, "getAverageNodeFee", c.Details.Address)
}

// Get the node's smoothing pool opt-in status
func (c *Node) GetSmoothingPoolRegistrationState(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeMgr, &c.Details.SmoothingPoolRegistrationState, "getSmoothingPoolRegistrationState", c.Details.Address)
}

// Get the time of the node's last smoothing pool opt-in / opt-out
func (c *Node) GetSmoothingPoolRegistrationChanged(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeMgr, &c.Details.SmoothingPoolRegistrationChanged.RawValue, "getSmoothingPoolRegistrationChanged", c.Details.Address)
}

// === NodeStaking ===

// Get the node's RPL stake
func (c *Node) GetRplStake(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeStaking, &c.Details.RplStake, "getNodeRPLStake", c.Details.Address)
}

// Get the node's effective RPL stake
func (c *Node) GetEffectiveRplStake(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeStaking, &c.Details.EffectiveRplStake, "getNodeEffectiveRPLStake", c.Details.Address)
}

// Get the node's minimum RPL stake to collateralize its minipools
func (c *Node) GetMinimumRplStake(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeStaking, &c.Details.MinimumRplStake, "getNodeMinimumRPLStake", c.Details.Address)
}

// Get the node's maximum RPL stake to collateralize its minipools
func (c *Node) GetMaximumRplStake(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeStaking, &c.Details.MaximumRplStake, "getNodeMaximumRPLStake", c.Details.Address)
}

// Get the time the node last staked RPL
func (c *Node) GetRplStakedTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeStaking, &c.Details.RplStakedTime.RawValue, "getNodeRPLStakedTime", c.Details.Address)
}

// Get the amount of ETH the node has borrowed from the deposit pool to create its minipools
func (c *Node) GetEthMatched(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeStaking, &c.Details.EthMatched, "getNodeETHMatched", c.Details.Address)
}

// Get the amount of ETH the node can still borrow from the deposit pool to create any new minipools
func (c *Node) GetEthMatchedLimit(mc *batch.MultiCaller) {
	core.AddCall(mc, c.nodeStaking, &c.Details.EthMatchedLimit, "getNodeETHMatchedLimit", c.Details.Address)
}

// === MinipoolManager ===

// Get the node's minipool count
func (c *Node) GetMinipoolCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.Details.MinipoolCount.RawValue, "getNodeMinipoolCount", c.Details.Address)
}

// Get the number of minipools owned by a node that are not finalised
func (c *Node) GetActiveMinipoolCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.Details.ActiveMinipoolCount.RawValue, "getNodeActiveMinipoolCount", c.Details.Address)
}

// Get the number of minipools owned by a node that are finalised
func (c *Node) GetFinalisedMinipoolCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.Details.FinalisedMinipoolCount.RawValue, "getNodeFinalisedMinipoolCount", c.Details.Address)
}

// Get the number of minipools owned by a node that are validating
func (c *Node) GetValidatingMinipoolCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mpMgr, &c.Details.ValidatingMinipoolCount.RawValue, "getNodeValidatingMinipoolCount", c.Details.Address)
}

// === Storage ===

// Get the node's withdrawal address
func (c *Node) GetWithdrawalAddress(mc *batch.MultiCaller) {
	core.AddCall(mc, c.storage, &c.Details.WithdrawalAddress, "getNodeWithdrawalAddress", c.Details.Address)
}

// Get the node's pending withdrawal address
func (c *Node) GetPendingWithdrawalAddress(mc *batch.MultiCaller) {
	core.AddCall(mc, c.storage, &c.Details.PendingWithdrawalAddress, "getNodePendingWithdrawalAddress", c.Details.Address)
}

// ====================
// === Transactions ===
// ====================

// === NodeDeposit ===

// Get info for making a node deposit and creating a new minipool
func (c *Node) Deposit(bondAmount *big.Int, minimumNodeFee float64, validatorPubkey types.ValidatorPubkey, validatorSignature types.ValidatorSignature, depositDataRoot common.Hash, salt *big.Int, expectedMinipoolAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nodeDeposit, "deposit", opts, bondAmount, eth.EthToWei(minimumNodeFee), validatorPubkey[:], validatorSignature[:], depositDataRoot, salt, expectedMinipoolAddress)
}

// Get info for making a node deposit and creating a new minipool by using the credit balance
func (c *Node) DepositWithCredit(bondAmount *big.Int, minimumNodeFee float64, validatorPubkey types.ValidatorPubkey, validatorSignature types.ValidatorSignature, depositDataRoot common.Hash, salt *big.Int, expectedMinipoolAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nodeDeposit, "depositWithCredit", opts, bondAmount, eth.EthToWei(minimumNodeFee), validatorPubkey[:], validatorSignature[:], depositDataRoot, salt, expectedMinipoolAddress)
}

// Get info for making a vacant minipool for solo staker migration
func (c *Node) CreateVacantMinipool(bondAmount *big.Int, minimumNodeFee float64, validatorPubkey types.ValidatorPubkey, salt *big.Int, expectedMinipoolAddress common.Address, currentBalance *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nodeDeposit, "createVacantMinipool", opts, bondAmount, eth.EthToWei(minimumNodeFee), validatorPubkey[:], salt, expectedMinipoolAddress, currentBalance)
}

// === NodeManager ===

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

// === Storage ===

// Get info for setting the node's withdrawal address
func (c *Node) SetWithdrawalAddress(withdrawalAddress common.Address, confirm bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.storage, "setWithdrawalAddress", opts, c.Details.Address, withdrawalAddress, confirm)
}

// Get info for confirming the node's withdrawal address
func (c *Node) ConfirmWithdrawalAddress(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.storage, "confirmWithdrawalAddress", opts, c.Details.Address)
}

// ===================
// === Sub-Getters ===
// ===================

// === DistributorFactory ===

// Get a node's distributor with details
func (c *Node) GetNodeDistributor(distributorAddress common.Address, includeDetails bool, opts *bind.CallOpts) (*NodeDistributor, error) {
	// Create the distributor
	distributor, err := NewNodeDistributor(c.rp, c.Details.Address, distributorAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating node distributor binding for node %s at %s: %w", c.Details.Address, distributorAddress.Hex(), err)
	}

	// Get details via a multicall query
	if includeDetails {
		err = c.rp.Query(func(mc *batch.MultiCaller) error {
			distributor.GetAllDetails(mc)
			return nil
		}, opts)
		if err != nil {
			return nil, fmt.Errorf("error getting node distributor for node %s at %s: %w", c.Details.Address, distributorAddress.Hex(), err)
		}
	}

	// Return
	return distributor, nil
}

// === MinipoolManager ===

// Get one of the node's minipool addresses by index
func (c *Node) GetMinipoolAddress(mc *batch.MultiCaller, address_Out *common.Address, index uint64) {
	core.AddCall(mc, c.mpMgr, address_Out, "getNodeMinipoolAt", c.Details.Address, big.NewInt(int64(index)))
}

// Get one of the node's validating minipool addresses by index
func (c *Node) GetValidatingMinipoolAddress(mc *batch.MultiCaller, address_Out *common.Address, index uint64) {
	core.AddCall(mc, c.mpMgr, address_Out, "getNodeValidatingMinipoolAt", c.Details.Address, big.NewInt(int64(index)))
}

// Get all of the node's minipool addresses in a standalone call.
// This will use an internal batched multicall invocation to retrieve all of them.
// Provide the value returned from GetMinipoolCount() in minipoolCount.
func (c *Node) GetMinipoolAddresses(minipoolCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, minipoolCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(minipoolCount), c.rp.AddressBatchSize,
		func(mc *batch.MultiCaller, index int) error {
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
func (c *Node) GetValidatingMinipoolAddresses(minipoolCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, minipoolCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(minipoolCount), c.rp.AddressBatchSize,
		func(mc *batch.MultiCaller, index int) error {
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
func (c *Node) GetExpectedMinipoolAddress(mc *batch.MultiCaller, address_Out *common.Address, salt *big.Int) {
	core.AddCall(mc, c.mpFactory, address_Out, "getExpectedAddress", c.Details.Address, salt)
}
