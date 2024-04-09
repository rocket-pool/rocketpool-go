package state

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/v2/core"
	"github.com/rocket-pool/rocketpool-go/v2/node"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
	"github.com/rocket-pool/rocketpool-go/v2/types"

	"golang.org/x/sync/errgroup"
)

const (
	legacyNodeBatchSize  int = 100
	nodeAddressBatchSize int = 1000
)

// Complete details for a node
type NativeNodeDetails struct {
	Exists                           bool
	RegistrationTime                 *big.Int
	TimezoneLocation                 string
	FeeDistributorInitialised        bool
	FeeDistributorAddress            common.Address
	RewardNetwork                    *big.Int
	RplStake                         *big.Int
	EffectiveRPLStake                *big.Int
	MinimumRPLStake                  *big.Int
	MaximumRPLStake                  *big.Int
	EthMatched                       *big.Int
	EthMatchedLimit                  *big.Int
	MinipoolCount                    *big.Int
	BalanceETH                       *big.Int
	BalanceRETH                      *big.Int
	BalanceRPL                       *big.Int
	BalanceOldRPL                    *big.Int
	DepositCreditBalance             *big.Int
	DistributorBalanceUserETH        *big.Int // Must call CalculateAverageFeeAndDistributorShares to get this
	DistributorBalanceNodeETH        *big.Int // Must call CalculateAverageFeeAndDistributorShares to get this
	WithdrawalAddress                common.Address
	PendingWithdrawalAddress         common.Address
	SmoothingPoolRegistrationState   bool
	SmoothingPoolRegistrationChanged *big.Int
	NodeAddress                      common.Address
	AverageNodeFee                   *big.Int // Must call CalculateAverageFeeAndDistributorShares to get this
	CollateralisationRatio           *big.Int
	DistributorBalance               *big.Int
}

// Gets the details for a node using the efficient multicall contract
func GetNativeNodeDetails(rp *rocketpool.RocketPool, contracts *NetworkContracts, nodeAddress common.Address) (NativeNodeDetails, error) {
	opts := &bind.CallOpts{
		BlockNumber: contracts.ElBlockNumber,
	}
	details := NativeNodeDetails{
		NodeAddress:               nodeAddress,
		AverageNodeFee:            big.NewInt(0),
		CollateralisationRatio:    big.NewInt(0),
		DistributorBalanceUserETH: big.NewInt(0),
		DistributorBalanceNodeETH: big.NewInt(0),
	}

	mc, err := batch.NewMultiCaller(rp.Client, contracts.MulticallerAddress)
	if err != nil {
		return NativeNodeDetails{}, fmt.Errorf("error creating multicaller: %w", err)
	}
	addNodeDetailsCalls(contracts, mc, &details, nodeAddress)

	_, err = mc.FlexibleCall(true, opts)
	if err != nil {
		return NativeNodeDetails{}, fmt.Errorf("error executing multicall: %w", err)
	}

	// Get the node's ETH balance
	details.BalanceETH, err = rp.Client.BalanceAt(context.Background(), nodeAddress, opts.BlockNumber)
	if err != nil {
		return NativeNodeDetails{}, err
	}

	// Get the distributor balance
	distributorBalance, err := rp.Client.BalanceAt(context.Background(), details.FeeDistributorAddress, opts.BlockNumber)
	if err != nil {
		return NativeNodeDetails{}, err
	}

	// Do some postprocessing on the node data
	details.DistributorBalance = distributorBalance

	// Fix the effective stake
	if details.EffectiveRPLStake.Cmp(details.MinimumRPLStake) == -1 {
		details.EffectiveRPLStake.SetUint64(0)
	}

	return details, nil
}

// Gets the details for all nodes using the efficient multicall contract
func GetAllNativeNodeDetails(rp *rocketpool.RocketPool, contracts *NetworkContracts) ([]NativeNodeDetails, error) {
	opts := &bind.CallOpts{
		BlockNumber: contracts.ElBlockNumber,
	}

	// Get the list of node addresses
	addresses, err := getNodeAddressesFast(rp, contracts, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node addresses: %w", err)
	}
	count := len(addresses)
	nodeDetails := make([]NativeNodeDetails, count)

	// Sync
	var wg errgroup.Group
	wg.SetLimit(threadLimit)

	// Run the getters in batches
	for i := 0; i < count; i += legacyNodeBatchSize {
		i := i
		max := i + legacyNodeBatchSize
		if max > count {
			max = count
		}

		wg.Go(func() error {
			var err error
			mc, err := batch.NewMultiCaller(rp.Client, contracts.MulticallerAddress)
			if err != nil {
				return err
			}
			for j := i; j < max; j++ {
				address := addresses[j]
				details := &nodeDetails[j]
				details.NodeAddress = address
				details.AverageNodeFee = big.NewInt(0)
				details.DistributorBalanceUserETH = big.NewInt(0)
				details.DistributorBalanceNodeETH = big.NewInt(0)
				details.CollateralisationRatio = big.NewInt(0)

				addNodeDetailsCalls(contracts, mc, details, address)
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error getting node details: %w", err)
	}

	// Get the balances of the nodes
	distributorAddresses := make([]common.Address, count)
	balances, err := contracts.BalanceBatcher.GetEthBalances(addresses, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node balances: %w", err)
	}
	for i, details := range nodeDetails {
		nodeDetails[i].BalanceETH = balances[i]
		distributorAddresses[i] = details.FeeDistributorAddress
	}

	// Get the balances of the distributors
	balances, err = contracts.BalanceBatcher.GetEthBalances(distributorAddresses, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting distributor balances: %w", err)
	}

	// Do some postprocessing on the node data
	for i := range nodeDetails {
		details := &nodeDetails[i]
		details.DistributorBalance = balances[i]

		// Fix the effective stake
		if details.EffectiveRPLStake.Cmp(details.MinimumRPLStake) == -1 {
			details.EffectiveRPLStake.SetUint64(0)
		}
	}

	return nodeDetails, nil
}

// Calculate the average node fee and user/node shares of the distributor's balance
func CalculateAverageFeeAndDistributorShares(rp *rocketpool.RocketPool, contracts *NetworkContracts, node NativeNodeDetails, minipoolDetails []*NativeMinipoolDetails) error {

	// Calculate the total of all fees for staking minipools that aren't finalized
	totalFee := big.NewInt(0)
	eligibleMinipools := int64(0)
	for _, mpd := range minipoolDetails {
		if mpd.Status == types.MinipoolStatus_Staking && !mpd.Finalised {
			totalFee.Add(totalFee, mpd.NodeFee)
			eligibleMinipools++
		}
	}

	// Get the average fee (0 if there aren't any minipools)
	if eligibleMinipools > 0 {
		node.AverageNodeFee.Div(totalFee, big.NewInt(eligibleMinipools))
	}

	// Get the user and node portions of the distributor balance
	distributorBalance := big.NewInt(0).Set(node.DistributorBalance)
	if distributorBalance.Cmp(big.NewInt(0)) > 0 {
		nodeBalance := big.NewInt(0)
		nodeBalance.Mul(distributorBalance, big.NewInt(1e18))
		nodeBalance.Div(nodeBalance, node.CollateralisationRatio)

		userBalance := big.NewInt(0)
		userBalance.Sub(distributorBalance, nodeBalance)

		if eligibleMinipools == 0 {
			// Split it based solely on the collateralisation ratio if there are no minipools (and hence no average fee)
			node.DistributorBalanceNodeETH = big.NewInt(0).Set(nodeBalance)
			node.DistributorBalanceUserETH = big.NewInt(0).Sub(distributorBalance, nodeBalance)
		} else {
			// Amount of ETH given to the NO as a commission
			commissionEth := big.NewInt(0)
			commissionEth.Mul(userBalance, node.AverageNodeFee)
			commissionEth.Div(commissionEth, big.NewInt(1e18))

			node.DistributorBalanceNodeETH.Add(nodeBalance, commissionEth)                         // Node gets their portion + commission on user portion
			node.DistributorBalanceUserETH.Sub(distributorBalance, node.DistributorBalanceNodeETH) // User gets balance - node share
		}

	} else {
		// No distributor balance
		node.DistributorBalanceNodeETH = big.NewInt(0)
		node.DistributorBalanceUserETH = big.NewInt(0)
	}

	return nil
}

// Get all node addresses using the multicaller
func getNodeAddressesFast(rp *rocketpool.RocketPool, contracts *NetworkContracts, opts *bind.CallOpts) ([]common.Address, error) {
	nodeMgr, err := node.NewNodeManager(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting node manager: %w", err)
	}

	// Get node count
	err = rp.Query(nil, opts, nodeMgr.NodeCount)
	if err != nil {
		return []common.Address{}, err
	}
	nodeCount := nodeMgr.NodeCount.Formatted()

	// Sync
	var wg errgroup.Group
	wg.SetLimit(threadLimit)
	addresses := make([]common.Address, nodeCount)

	// Run the getters in batches
	count := int(nodeCount)
	for i := 0; i < count; i += nodeAddressBatchSize {
		i := i
		max := i + nodeAddressBatchSize
		if max > count {
			max = count
		}

		wg.Go(func() error {
			var err error
			mc, err := batch.NewMultiCaller(rp.Client, contracts.MulticallerAddress)
			if err != nil {
				return err
			}
			for j := i; j < max; j++ {
				core.AddCall(mc, contracts.RocketNodeManager, &addresses[j], "getNodeAt", big.NewInt(int64(j)))
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error getting node addresses: %w", err)
	}

	return addresses, nil
}

// Add all of the calls for the node details to the multicaller
func addNodeDetailsCalls(contracts *NetworkContracts, mc *batch.MultiCaller, details *NativeNodeDetails, address common.Address) {
	core.AddCall(mc, contracts.RocketNodeManager, &details.Exists, "getNodeExists", address)
	core.AddCall(mc, contracts.RocketNodeManager, &details.RegistrationTime, "getNodeRegistrationTime", address)
	core.AddCall(mc, contracts.RocketNodeManager, &details.TimezoneLocation, "getNodeTimezoneLocation", address)
	core.AddCall(mc, contracts.RocketNodeManager, &details.FeeDistributorInitialised, "getFeeDistributorInitialised", address)
	core.AddCall(mc, contracts.RocketNodeDistributorFactory, &details.FeeDistributorAddress, "getProxyAddress", address)
	core.AddCall(mc, contracts.RocketNodeManager, &details.RewardNetwork, "getRewardNetwork", address)
	core.AddCall(mc, contracts.RocketNodeStaking, &details.RplStake, "getNodeRPLStake", address)
	core.AddCall(mc, contracts.RocketNodeStaking, &details.EffectiveRPLStake, "getNodeEffectiveRPLStake", address)
	core.AddCall(mc, contracts.RocketNodeStaking, &details.MinimumRPLStake, "getNodeMinimumRPLStake", address)
	core.AddCall(mc, contracts.RocketNodeStaking, &details.MaximumRPLStake, "getNodeMaximumRPLStake", address)
	core.AddCall(mc, contracts.RocketNodeStaking, &details.EthMatched, "getNodeETHMatched", address)
	core.AddCall(mc, contracts.RocketNodeStaking, &details.EthMatchedLimit, "getNodeETHMatchedLimit", address)
	core.AddCall(mc, contracts.RocketMinipoolManager, &details.MinipoolCount, "getNodeMinipoolCount", address)
	core.AddCall(mc, contracts.RocketTokenRETH, &details.BalanceRETH, "balanceOf", address)
	core.AddCall(mc, contracts.RocketTokenRPL, &details.BalanceRPL, "balanceOf", address)
	core.AddCall(mc, contracts.RocketTokenRPLFixedSupply, &details.BalanceOldRPL, "balanceOf", address)
	core.AddCall(mc, contracts.RocketStorage, &details.WithdrawalAddress, "getNodeWithdrawalAddress", address)
	core.AddCall(mc, contracts.RocketStorage, &details.PendingWithdrawalAddress, "getNodePendingWithdrawalAddress", address)
	core.AddCall(mc, contracts.RocketNodeManager, &details.SmoothingPoolRegistrationState, "getSmoothingPoolRegistrationState", address)
	core.AddCall(mc, contracts.RocketNodeManager, &details.SmoothingPoolRegistrationChanged, "getSmoothingPoolRegistrationChanged", address)

	// Atlas
	core.AddCall(mc, contracts.RocketNodeDeposit, &details.DepositCreditBalance, "getNodeDepositCredit", address)
	core.AddCall(mc, contracts.RocketNodeStaking, &details.CollateralisationRatio, "getNodeETHCollateralisationRatio", address)
}
