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
	// The address of this node
	Address common.Address

	// True if the node exists (i.e. there is a registered node at this address)
	Exists *core.SimpleField[bool]

	// The time that the node was registered with the network
	RegistrationTime *core.FormattedUint256Field[time.Time]

	// The node's timezone location
	TimezoneLocation *core.SimpleField[string]

	// The network ID for the node's rewards
	RewardNetwork *core.FormattedUint256Field[uint64]

	// The node's average minipool fee (commission)
	AverageFee *core.FormattedUint256Field[float64]

	// The node's smoothing pool opt-in status
	SmoothingPoolRegistrationState *core.SimpleField[bool]

	// The time of the node's last smoothing pool opt-in / opt-out
	SmoothingPoolRegistrationChanged *core.FormattedUint256Field[time.Time]

	// True if the node's fee distributor has been initialized yet
	IsFeeDistributorInitialized *core.SimpleField[bool]

	// The node's fee distributor address
	DistributorAddress *core.SimpleField[common.Address]

	// The amount of ETH in the node's deposit credit bank
	Credit *core.SimpleField[*big.Int]

	// The node's RPL stake
	RplStake *core.SimpleField[*big.Int]

	// The node's effective RPL stake
	EffectiveRplStake *core.SimpleField[*big.Int]

	// The node's minimum RPL stake to collateralize its minipools
	MinimumRplStake *core.SimpleField[*big.Int]

	// The node's maximum RPL stake to collateralize its minipools
	MaximumRplStake *core.SimpleField[*big.Int]

	// The time the node last staked RPL
	RplStakedTime *core.FormattedUint256Field[time.Time]

	// The amount of ETH the node has borrowed from the deposit pool to create its minipools
	EthMatched *core.SimpleField[*big.Int]

	// The amount of ETH the node can still borrow from the deposit pool to create any new minipools
	EthMatchedLimit *core.SimpleField[*big.Int]

	// The number of minipools owned by the node count
	MinipoolCount *core.FormattedUint256Field[uint64]

	// The number of minipools owned by the node that are not finalised
	ActiveMinipoolCount *core.FormattedUint256Field[uint64]

	// The number of minipools owned by a node that are finalised
	FinalisedMinipoolCount *core.FormattedUint256Field[uint64]

	// The number of minipools owned by a node that are validating
	ValidatingMinipoolCount *core.FormattedUint256Field[uint64]

	// The node's primary withdrawal address
	PrimaryWithdrawalAddress *core.SimpleField[common.Address]

	// The node's pending primary withdrawal address
	PendingPrimaryWithdrawalAddress *core.SimpleField[common.Address]

	// The node's RPL withdrawal address
	IsRplWithdrawalAddressSet *core.SimpleField[bool]

	// The node's RPL withdrawal address
	RplWithdrawalAddress *core.SimpleField[common.Address]

	// The node's pending RPL withdrawal address
	PendingRplWithdrawalAddress *core.SimpleField[common.Address]

	// The amount of RPL locked as part of active PDAO proposals or challenges
	RplLocked *core.SimpleField[*big.Int]

	// The address that the provided node has currently delegated voting power to
	CurrentVotingDelegate *core.SimpleField[common.Address]

	// Whether or not on-chain voting has been initialized for the given node
	IsVotingInitialized *core.SimpleField[bool]

	// === Internal fields ===
	rp            *rocketpool.RocketPool
	distFactory   *core.Contract
	networkVoting *core.Contract
	nodeDeposit   *core.Contract
	nodeMgr       *core.Contract
	nodeStaking   *core.Contract
	mpFactory     *core.Contract
	mpMgr         *core.Contract
	storage       *core.Contract
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
	networkVoting, err := rp.GetContract(rocketpool.ContractName_RocketNetworkVoting)
	if err != nil {
		return nil, fmt.Errorf("error getting network voting binding: %w", err)
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
		Address: address,

		// DistributorFactory
		DistributorAddress: core.NewSimpleField[common.Address](distFactory, "getProxyAddress", address),

		// NetworkVoting
		CurrentVotingDelegate: core.NewSimpleField[common.Address](networkVoting, "getCurrentDelegate", address),
		IsVotingInitialized:   core.NewSimpleField[bool](networkVoting, "getVotingInitialised", address),

		// NodeDeposit
		Credit: core.NewSimpleField[*big.Int](nodeDeposit, "getNodeDepositCredit", address),

		// NodeManager
		Exists:                           core.NewSimpleField[bool](nodeManager, "getNodeExists", address),
		RegistrationTime:                 core.NewFormattedUint256Field[time.Time](nodeManager, "getNodeRegistrationTime", address),
		TimezoneLocation:                 core.NewSimpleField[string](nodeManager, "getNodeTimezoneLocation", address),
		RewardNetwork:                    core.NewFormattedUint256Field[uint64](nodeManager, "getRewardNetwork", address),
		IsFeeDistributorInitialized:      core.NewSimpleField[bool](nodeManager, "getFeeDistributorInitialised", address),
		AverageFee:                       core.NewFormattedUint256Field[float64](nodeManager, "getAverageNodeFee", address),
		SmoothingPoolRegistrationState:   core.NewSimpleField[bool](nodeManager, "getSmoothingPoolRegistrationState", address),
		SmoothingPoolRegistrationChanged: core.NewFormattedUint256Field[time.Time](nodeManager, "getSmoothingPoolRegistrationChanged", address),
		IsRplWithdrawalAddressSet:        core.NewSimpleField[bool](nodeManager, "getNodeRPLWithdrawalAddressIsSet", address),
		RplWithdrawalAddress:             core.NewSimpleField[common.Address](nodeManager, "getNodeRPLWithdrawalAddress", address),
		PendingRplWithdrawalAddress:      core.NewSimpleField[common.Address](nodeManager, "getNodePendingRPLWithdrawalAddress", address),

		// NodeStaking
		RplStake:          core.NewSimpleField[*big.Int](nodeStaking, "getNodeRPLStake", address),
		EffectiveRplStake: core.NewSimpleField[*big.Int](nodeStaking, "getNodeEffectiveRPLStake", address),
		MinimumRplStake:   core.NewSimpleField[*big.Int](nodeStaking, "getNodeMinimumRPLStake", address),
		MaximumRplStake:   core.NewSimpleField[*big.Int](nodeStaking, "getNodeMaximumRPLStake", address),
		RplStakedTime:     core.NewFormattedUint256Field[time.Time](nodeStaking, "getNodeRPLStakedTime", address),
		EthMatched:        core.NewSimpleField[*big.Int](nodeStaking, "getNodeETHMatched", address),
		EthMatchedLimit:   core.NewSimpleField[*big.Int](nodeStaking, "getNodeETHMatchedLimit", address),
		RplLocked:         core.NewSimpleField[*big.Int](nodeStaking, "getNodeRPLLocked", address),

		// MinipoolManager
		MinipoolCount:           core.NewFormattedUint256Field[uint64](minipoolManager, "getNodeMinipoolCount", address),
		ActiveMinipoolCount:     core.NewFormattedUint256Field[uint64](minipoolManager, "getNodeActiveMinipoolCount", address),
		FinalisedMinipoolCount:  core.NewFormattedUint256Field[uint64](minipoolManager, "getNodeFinalisedMinipoolCount", address),
		ValidatingMinipoolCount: core.NewFormattedUint256Field[uint64](minipoolManager, "getNodeValidatingMinipoolCount", address),

		// Storage
		PrimaryWithdrawalAddress:        core.NewSimpleField[common.Address](rp.Storage.Contract, "getNodeWithdrawalAddress", address),
		PendingPrimaryWithdrawalAddress: core.NewSimpleField[common.Address](rp.Storage.Contract, "getNodePendingWithdrawalAddress", address),

		rp:            rp,
		distFactory:   distFactory,
		networkVoting: networkVoting,
		nodeDeposit:   nodeDeposit,
		nodeMgr:       nodeManager,
		nodeStaking:   nodeStaking,
		mpFactory:     minipoolFactory,
		mpMgr:         minipoolManager,
		storage:       rp.Storage.Contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the address that the node has delegated voting power to at the given block
func (c *Node) GetVotingDelegateAtBlock(mc *batch.MultiCaller, delegate_Out *common.Address, blockNumber uint32) {
	core.AddCall(mc, c.networkVoting, delegate_Out, "getDelegate", c.Address, blockNumber)
}

// Get the voting power of the given node at the provided block
func (c *Node) GetVotingPowerAtBlock(mc *batch.MultiCaller, power_Out **big.Int, blockNumber uint32) {
	core.AddCall(mc, c.networkVoting, power_Out, "getVotingPower", c.Address, blockNumber)
}

// ====================
// === Transactions ===
// ====================

// === NetworkVoting ===

// Get info for initializing on-chain voting for the node
func (c *Node) InitializeVoting(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.networkVoting, "initialiseVoting", opts)
}

// Get info for setting the voting delegate for the node
func (c *Node) SetVotingDelegate(newDelegate common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.networkVoting, "setDelegate", opts, newDelegate)
}

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

// Get info for setting the RPL withdrawal address
func (c *Node) SetRplWithdrawalAddress(withdrawalAddress common.Address, confirm bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nodeMgr, "setRPLWithdrawalAddress", opts, c.Address, withdrawalAddress, confirm)
}

// Get info for confirming the RPL withdrawal address
func (c *Node) ConfirmRplWithdrawalAddress(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.nodeMgr, "confirmRPLWithdrawalAddress", opts, c.Address)
}

// === Storage ===

// Get info for setting the node's primary withdrawal address
func (c *Node) SetPrimaryWithdrawalAddress(withdrawalAddress common.Address, confirm bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.storage, "setWithdrawalAddress", opts, c.Address, withdrawalAddress, confirm)
}

// Get info for confirming the node's primary withdrawal address
func (c *Node) ConfirmPrimaryWithdrawalAddress(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.storage, "confirmWithdrawalAddress", opts, c.Address)
}

// ===================
// === Sub-Getters ===
// ===================

// === DistributorFactory ===

// Get a node's distributor with details
func (c *Node) GetNodeDistributor(distributorAddress common.Address, includeDetails bool, opts *bind.CallOpts) (*NodeDistributor, error) {
	// Create the distributor
	distributor, err := NewNodeDistributor(c.rp, c.Address, distributorAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating node distributor binding for node %s at %s: %w", c.Address, distributorAddress.Hex(), err)
	}

	// Get details via a multicall query
	if includeDetails {
		err = c.rp.Query(func(mc *batch.MultiCaller) error {
			core.QueryAllFields(distributor, mc)
			return nil
		}, opts)
		if err != nil {
			return nil, fmt.Errorf("error getting node distributor for node %s at %s: %w", c.Address, distributorAddress.Hex(), err)
		}
	}

	// Return
	return distributor, nil
}

// === MinipoolManager ===

// Get one of the node's minipool addresses by index
func (c *Node) GetMinipoolAddress(mc *batch.MultiCaller, address_Out *common.Address, index uint64) {
	core.AddCall(mc, c.mpMgr, address_Out, "getNodeMinipoolAt", c.Address, big.NewInt(int64(index)))
}

// Get one of the node's validating minipool addresses by index
func (c *Node) GetValidatingMinipoolAddress(mc *batch.MultiCaller, address_Out *common.Address, index uint64) {
	core.AddCall(mc, c.mpMgr, address_Out, "getNodeValidatingMinipoolAt", c.Address, big.NewInt(int64(index)))
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
	core.AddCall(mc, c.mpFactory, address_Out, "getExpectedAddress", c.Address, salt)
}
