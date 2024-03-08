package state

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/beacon"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/minipool"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"

	"golang.org/x/sync/errgroup"
)

const (
	minipoolBatchSize              int = 100
	minipoolCompleteShareBatchSize int = 500
	minipoolAddressBatchSize       int = 1000
	minipoolVersionBatchSize       int = 500
)

// Complete details for a minipool
type NativeMinipoolDetails struct {
	// Redstone
	Exists                            bool
	MinipoolAddress                   common.Address
	Pubkey                            beacon.ValidatorPubkey
	StatusRaw                         uint8
	StatusBlock                       *big.Int
	StatusTime                        *big.Int
	Finalised                         bool
	DepositTypeRaw                    uint8
	NodeFee                           *big.Int
	NodeDepositBalance                *big.Int
	NodeDepositAssigned               bool
	UserDepositBalance                *big.Int
	UserDepositAssigned               bool
	UserDepositAssignedTime           *big.Int
	UseLatestDelegate                 bool
	Delegate                          common.Address
	PreviousDelegate                  common.Address
	EffectiveDelegate                 common.Address
	PenaltyCount                      *big.Int
	PenaltyRate                       *big.Int
	NodeAddress                       common.Address
	Version                           uint8
	Balance                           *big.Int // Contract balance
	DistributableBalance              *big.Int // Contract balance minus node op refund
	NodeShareOfBalance                *big.Int // Result of calculateNodeShare(contract balance)
	UserShareOfBalance                *big.Int // Result of calculateUserShare(contract balance)
	NodeRefundBalance                 *big.Int
	WithdrawalCredentials             common.Hash
	Status                            types.MinipoolStatus
	DepositType                       types.MinipoolDeposit
	NodeShareOfBalanceIncludingBeacon *big.Int // Must call CalculateCompleteMinipoolShares to get this
	UserShareOfBalanceIncludingBeacon *big.Int // Must call CalculateCompleteMinipoolShares to get this
	NodeShareOfBeaconBalance          *big.Int // Must call CalculateCompleteMinipoolShares to get this
	UserShareOfBeaconBalance          *big.Int // Must call CalculateCompleteMinipoolShares to get this

	// Atlas
	UserDistributed              bool
	Slashed                      bool
	IsVacant                     bool
	LastBondReductionTime        *big.Int
	LastBondReductionPrevValue   *big.Int
	LastBondReductionPrevNodeFee *big.Int
	ReduceBondTime               *big.Int
	ReduceBondCancelled          bool
	ReduceBondValue              *big.Int
	PreMigrationBalance          *big.Int
}

// Gets the details for a minipool using the efficient multicall contract
func GetNativeMinipoolDetails(rp *rocketpool.RocketPool, contracts *NetworkContracts, minipoolAddress common.Address) (NativeMinipoolDetails, error) {
	opts := &bind.CallOpts{
		BlockNumber: contracts.ElBlockNumber,
	}

	mpMgr, err := minipool.NewMinipoolManager(rp)
	if err != nil {
		return NativeMinipoolDetails{}, fmt.Errorf("error creating minipool manager: %w", err)
	}

	details := NativeMinipoolDetails{}
	details.MinipoolAddress = minipoolAddress

	var version uint8
	err = rp.Query(func(mc *batch.MultiCaller) error {
		rocketpool.GetContractVersion(mc, &version, minipoolAddress)
		return nil
	}, opts)
	if err != nil {
		return NativeMinipoolDetails{}, fmt.Errorf("error getting minipool version: %w", err)
	}
	details.Version = version

	mc, err := batch.NewMultiCaller(rp.Client, contracts.MulticallerAddress)
	if err != nil {
		return NativeMinipoolDetails{}, fmt.Errorf("error creating multicaller: %w", err)
	}
	addMinipoolDetailsCalls(rp, mpMgr, contracts, mc, &details, opts)

	_, err = mc.FlexibleCall(true, opts)
	if err != nil {
		return NativeMinipoolDetails{}, fmt.Errorf("error executing multicall: %w", err)
	}

	fixupMinipoolDetails(rp, &details, opts)

	return details, nil
}

// Gets the minpool details for a node using the efficient multicall contract
func GetNodeNativeMinipoolDetails(rp *rocketpool.RocketPool, contracts *NetworkContracts, nodeAddress common.Address) ([]NativeMinipoolDetails, error) {
	opts := &bind.CallOpts{
		BlockNumber: contracts.ElBlockNumber,
	}

	// Get the list of minipool addresses for this node
	addresses, err := getNodeMinipoolAddressesFast(rp, contracts, nodeAddress, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool addresses: %w", err)
	}

	// Get the list of minipool versions
	versions, err := getMinipoolVersionsFast(rp, contracts, addresses, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool versions: %w", err)
	}

	// Get the minipool details
	return getBulkMinipoolDetails(rp, contracts, addresses, versions, opts)
}

// Gets all minpool details using the efficient multicall contract
func GetAllNativeMinipoolDetails(rp *rocketpool.RocketPool, contracts *NetworkContracts) ([]NativeMinipoolDetails, error) {
	opts := &bind.CallOpts{
		BlockNumber: contracts.ElBlockNumber,
	}

	// Get the list of all minipool addresses
	addresses, err := getAllMinipoolAddressesFast(rp, contracts, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool addresses: %w", err)
	}

	// Get the list of minipool versions
	versions, err := getMinipoolVersionsFast(rp, contracts, addresses, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool versions: %w", err)
	}

	// Get the minipool details
	return getBulkMinipoolDetails(rp, contracts, addresses, versions, opts)
}

// Calculate the node and user shares of the total minipool balance, including the portion on the Beacon chain
func CalculateCompleteMinipoolShares(rp *rocketpool.RocketPool, contracts *NetworkContracts, minipoolDetails []*NativeMinipoolDetails, beaconBalances []*big.Int) error {
	opts := &bind.CallOpts{
		BlockNumber: contracts.ElBlockNumber,
	}

	mpMgr, err := minipool.NewMinipoolManager(rp)
	if err != nil {
		return fmt.Errorf("error creating minipool manager: %w", err)
	}

	var wg errgroup.Group
	wg.SetLimit(threadLimit)
	count := len(minipoolDetails)
	for i := 0; i < count; i += minipoolCompleteShareBatchSize {
		i := i
		max := i + minipoolCompleteShareBatchSize
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

				// Make the minipool contract
				details := minipoolDetails[j]
				mp, err := mpMgr.NewMinipoolFromVersion(details.MinipoolAddress, details.Version)
				if err != nil {
					return err
				}
				mpCommon := mp.Common()

				// Calculate the Beacon shares
				beaconBalance := big.NewInt(0).Set(beaconBalances[j])
				if beaconBalance.Cmp(zero) > 0 {
					mpCommon.CalculateNodeShare(mc, &details.NodeShareOfBeaconBalance, beaconBalance)
					mpCommon.CalculateUserShare(mc, &details.UserShareOfBeaconBalance, beaconBalance)
				} else {
					details.NodeShareOfBeaconBalance = big.NewInt(0)
					details.UserShareOfBeaconBalance = big.NewInt(0)
				}

				// Calculate the total balance
				totalBalance := big.NewInt(0).Set(beaconBalances[j])      // Total balance = beacon balance
				totalBalance.Add(totalBalance, details.Balance)           // Add contract balance
				totalBalance.Sub(totalBalance, details.NodeRefundBalance) // Remove node refund

				// Calculate the node and user shares
				if totalBalance.Cmp(zero) > 0 {
					mpCommon.CalculateNodeShare(mc, &details.NodeShareOfBalanceIncludingBeacon, totalBalance)
					mpCommon.CalculateUserShare(mc, &details.UserShareOfBalanceIncludingBeacon, totalBalance)
				} else {
					details.NodeShareOfBalanceIncludingBeacon = big.NewInt(0)
					details.UserShareOfBalanceIncludingBeacon = big.NewInt(0)
				}
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return fmt.Errorf("error calculating minipool shares: %w", err)
	}

	return nil
}

// Get all minipool addresses using the multicaller
func getNodeMinipoolAddressesFast(rp *rocketpool.RocketPool, contracts *NetworkContracts, nodeAddress common.Address, opts *bind.CallOpts) ([]common.Address, error) {
	// Get node
	node, err := node.NewNode(rp, nodeAddress)
	if err != nil {
		return nil, fmt.Errorf("error getting node: %w", err)
	}

	// Get minipool count
	err = rp.Query(nil, opts, node.MinipoolCount)
	if err != nil {
		return []common.Address{}, err
	}
	minipoolCount := node.MinipoolCount.Formatted()

	// Sync
	var wg errgroup.Group
	wg.SetLimit(threadLimit)
	addresses := make([]common.Address, minipoolCount)

	// Run the getters in batches
	count := int(minipoolCount)
	for i := 0; i < count; i += minipoolAddressBatchSize {
		i := i
		max := i + minipoolAddressBatchSize
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
				core.AddCall(mc, contracts.RocketMinipoolManager, &addresses[j], "getNodeMinipoolAt", nodeAddress, big.NewInt(int64(j)))
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error getting minipool addresses for node %s: %w", nodeAddress.Hex(), err)
	}

	return addresses, nil
}

// Get all minipool addresses using the multicaller
func getAllMinipoolAddressesFast(rp *rocketpool.RocketPool, contracts *NetworkContracts, opts *bind.CallOpts) ([]common.Address, error) {
	// Get minipool manager
	mgr, err := minipool.NewMinipoolManager(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting node: %w", err)
	}

	// Get minipool count
	var minipoolCount uint64
	err = rp.Query(nil, opts, mgr.MinipoolCount)
	if err != nil {
		return []common.Address{}, err
	}

	// Sync
	var wg errgroup.Group
	wg.SetLimit(threadLimit)
	addresses := make([]common.Address, minipoolCount)

	// Run the getters in batches
	count := int(minipoolCount)
	for i := 0; i < count; i += minipoolAddressBatchSize {
		i := i
		max := i + minipoolAddressBatchSize
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
				core.AddCall(mc, contracts.RocketMinipoolManager, &addresses[j], "getMinipoolAt", big.NewInt(int64(j)))
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error getting all minipool addresses: %w", err)
	}

	return addresses, nil
}

// Get minipool versions using the multicaller
func getMinipoolVersionsFast(rp *rocketpool.RocketPool, contracts *NetworkContracts, addresses []common.Address, opts *bind.CallOpts) ([]uint8, error) {
	// Sync
	var wg errgroup.Group
	wg.SetLimit(threadLimit)

	// Run the getters in batches
	count := len(addresses)
	versions := make([]uint8, count)
	for i := 0; i < count; i += minipoolVersionBatchSize {
		i := i
		max := i + minipoolVersionBatchSize
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
				contract, err := rocketpool.GetRocketVersionContractForAddress(rp, addresses[j])
				if err != nil {
					return fmt.Errorf("error creating version contract for minipool %s: %w", addresses[j].Hex(), err)
				}
				core.AddCall(mc, contract, &versions[j], "version")
			}
			results, err := mc.FlexibleCall(false, opts) // Allow calls to fail - necessary for Prater
			for j, result := range results {
				if !result {
					versions[j+i] = 1 // Anything that failed the version check didn't have the method yet so it must be v1
				}
			}
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error getting minipool versions: %w", err)
	}

	return versions, nil
}

// Get multiple minipool details at once
func getBulkMinipoolDetails(rp *rocketpool.RocketPool, contracts *NetworkContracts, addresses []common.Address, versions []uint8, opts *bind.CallOpts) ([]NativeMinipoolDetails, error) {
	minipoolDetails := make([]NativeMinipoolDetails, len(addresses))

	mpMgr, err := minipool.NewMinipoolManager(rp)
	if err != nil {
		return nil, fmt.Errorf("error creating minipool manager: %w", err)
	}

	// Get the balances of the minipools
	balances, err := contracts.BalanceBatcher.GetEthBalances(addresses, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool balances: %w", err)
	}
	for i := range minipoolDetails {
		minipoolDetails[i].Balance = balances[i]
	}

	// Round 1: most of the details
	var wg errgroup.Group
	wg.SetLimit(threadLimit)
	count := len(addresses)
	for i := 0; i < count; i += minipoolBatchSize {
		i := i
		max := i + minipoolBatchSize
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
				details := &minipoolDetails[j]
				details.MinipoolAddress = address
				details.Version = versions[j]

				addMinipoolDetailsCalls(rp, mpMgr, contracts, mc, details, opts)
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error getting minipool details r1: %w", err)
	}

	// Round 2: NodeShare and UserShare once the refund amount has been populated
	var wg2 errgroup.Group
	wg2.SetLimit(threadLimit)
	for i := 0; i < count; i += minipoolBatchSize {
		i := i
		max := i + minipoolBatchSize
		if max > count {
			max = count
		}

		wg2.Go(func() error {
			var err error
			mc, err := batch.NewMultiCaller(rp.Client, contracts.MulticallerAddress)
			if err != nil {
				return err
			}
			for j := i; j < max; j++ {
				details := &minipoolDetails[j]
				details.Version = versions[j]
				addMinipoolShareCalls(rp, mpMgr, contracts, mc, details, opts)
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}

			return nil
		})
	}

	if err := wg2.Wait(); err != nil {
		return nil, fmt.Errorf("error getting minipool details r2: %w", err)
	}

	// Postprocess the minipools
	for i := range minipoolDetails {
		fixupMinipoolDetails(rp, &minipoolDetails[i], opts)
	}

	return minipoolDetails, nil
}

// Add all of the calls for the minipool details to the multicaller
func addMinipoolDetailsCalls(rp *rocketpool.RocketPool, mpMgr *minipool.MinipoolManager, contracts *NetworkContracts, mc *batch.MultiCaller, details *NativeMinipoolDetails, opts *bind.CallOpts) error {
	// Create the minipool contract binding
	address := details.MinipoolAddress
	mp, err := mpMgr.NewMinipoolFromVersion(address, details.Version)
	if err != nil {
		return err
	}
	mpContract := mp.GetContract()

	core.AddCall(mc, contracts.RocketMinipoolManager, &details.Exists, "getMinipoolExists", address)
	core.AddCall(mc, contracts.RocketMinipoolManager, &details.Pubkey, "getMinipoolPubkey", address)
	core.AddCall(mc, contracts.RocketMinipoolManager, &details.WithdrawalCredentials, "getMinipoolWithdrawalCredentials", address)
	core.AddCall(mc, contracts.RocketMinipoolManager, &details.Slashed, "getMinipoolRPLSlashed", address)
	core.AddCall(mc, mpContract, &details.StatusRaw, "getStatus")
	core.AddCall(mc, mpContract, &details.StatusBlock, "getStatusBlock")
	core.AddCall(mc, mpContract, &details.StatusTime, "getStatusTime")
	core.AddCall(mc, mpContract, &details.Finalised, "getFinalised")
	core.AddCall(mc, mpContract, &details.NodeFee, "getNodeFee")
	core.AddCall(mc, mpContract, &details.NodeDepositBalance, "getNodeDepositBalance")
	core.AddCall(mc, mpContract, &details.NodeDepositAssigned, "getNodeDepositAssigned")
	core.AddCall(mc, mpContract, &details.UserDepositBalance, "getUserDepositBalance")
	core.AddCall(mc, mpContract, &details.UserDepositAssigned, "getUserDepositAssigned")
	core.AddCall(mc, mpContract, &details.UserDepositAssignedTime, "getUserDepositAssignedTime")
	core.AddCall(mc, mpContract, &details.UseLatestDelegate, "getUseLatestDelegate")
	core.AddCall(mc, mpContract, &details.Delegate, "getDelegate")
	core.AddCall(mc, mpContract, &details.PreviousDelegate, "getPreviousDelegate")
	core.AddCall(mc, mpContract, &details.EffectiveDelegate, "getEffectiveDelegate")
	core.AddCall(mc, mpContract, &details.NodeAddress, "getNodeAddress")
	core.AddCall(mc, mpContract, &details.NodeRefundBalance, "getNodeRefundBalance")

	if details.Version < 3 {
		// These fields are all v3+ only
		details.UserDistributed = false
		details.LastBondReductionTime = big.NewInt(0)
		details.LastBondReductionPrevValue = big.NewInt(0)
		details.LastBondReductionPrevNodeFee = big.NewInt(0)
		details.IsVacant = false
		details.ReduceBondTime = big.NewInt(0)
		details.ReduceBondCancelled = false
		details.ReduceBondValue = big.NewInt(0)
		details.PreMigrationBalance = big.NewInt(0)
	} else {
		core.AddCall(mc, mpContract, &details.UserDistributed, "getUserDistributed")
		core.AddCall(mc, mpContract, &details.IsVacant, "getVacant")
		core.AddCall(mc, mpContract, &details.PreMigrationBalance, "getPreMigrationBalance")

		// If minipool v3 exists, RocketMinipoolBondReducer exists so this is safe
		core.AddCall(mc, contracts.RocketMinipoolBondReducer, &details.ReduceBondTime, "getReduceBondTime", address)
		core.AddCall(mc, contracts.RocketMinipoolBondReducer, &details.ReduceBondCancelled, "getReduceBondCancelled", address)
		core.AddCall(mc, contracts.RocketMinipoolBondReducer, &details.LastBondReductionTime, "getLastBondReductionTime", address)
		core.AddCall(mc, contracts.RocketMinipoolBondReducer, &details.LastBondReductionPrevValue, "getLastBondReductionPrevValue", address)
		core.AddCall(mc, contracts.RocketMinipoolBondReducer, &details.LastBondReductionPrevNodeFee, "getLastBondReductionPrevNodeFee", address)
		core.AddCall(mc, contracts.RocketMinipoolBondReducer, &details.ReduceBondValue, "getReduceBondValue", address)
	}

	penaltyCountKey := crypto.Keccak256Hash([]byte("network.penalties.penalty"), address.Bytes())
	core.AddCall(mc, contracts.RocketStorage, &details.PenaltyCount, "getUint", penaltyCountKey)

	penaltyRatekey := crypto.Keccak256Hash([]byte("minipool.penalty.rate"), address.Bytes())
	core.AddCall(mc, contracts.RocketStorage, &details.PenaltyRate, "getUint", penaltyRatekey)

	// Query the minipool manager using the delegate-invariant function
	core.AddCall(mc, contracts.RocketMinipoolManager, &details.DepositTypeRaw, "getMinipoolDepositType", address)

	return nil
}

// Add the calls for the minipool node and user share to the multicaller
func addMinipoolShareCalls(rp *rocketpool.RocketPool, mpMgr *minipool.MinipoolManager, contracts *NetworkContracts, mc *batch.MultiCaller, details *NativeMinipoolDetails, opts *bind.CallOpts) error {
	// Create the minipool contract binding
	address := details.MinipoolAddress
	mp, err := mpMgr.NewMinipoolFromVersion(address, details.Version)
	if err != nil {
		return err
	}
	mpContract := mp.GetContract()

	details.DistributableBalance = big.NewInt(0).Sub(details.Balance, details.NodeRefundBalance)
	if details.DistributableBalance.Cmp(zero) >= 0 {
		core.AddCall(mc, mpContract, &details.NodeShareOfBalance, "calculateNodeShare", details.DistributableBalance)
		core.AddCall(mc, mpContract, &details.UserShareOfBalance, "calculateUserShare", details.DistributableBalance)
	} else {
		details.NodeShareOfBalance = big.NewInt(0)
		details.UserShareOfBalance = big.NewInt(0)
	}

	return nil
}

// Fixes a minipool details struct with supplemental logic
func fixupMinipoolDetails(rp *rocketpool.RocketPool, details *NativeMinipoolDetails, opts *bind.CallOpts) error {

	details.Status = types.MinipoolStatus(details.StatusRaw)
	details.DepositType = types.MinipoolDeposit(details.DepositTypeRaw)

	return nil
}
