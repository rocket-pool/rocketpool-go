package minipool

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	rptypes "github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ==================
// === Interfaces ===
// ==================

type Minipool interface {
	// Calls
	GetStatus(mc *multicall.MultiCaller)
	GetStatusBlock(mc *multicall.MultiCaller)
	GetStatusTime(mc *multicall.MultiCaller)
	GetFinalised(mc *multicall.MultiCaller)
	GetDepositType(mc *multicall.MultiCaller)
	GetNodeAddress(mc *multicall.MultiCaller)
	GetNodeFee(mc *multicall.MultiCaller)
	GetNodeDepositBalance(mc *multicall.MultiCaller)
	GetNodeRefundBalance(mc *multicall.MultiCaller)
	GetNodeDepositAssigned(mc *multicall.MultiCaller)
	GetUserDepositBalance(mc *multicall.MultiCaller)
	GetUserDepositAssigned(mc *multicall.MultiCaller)
	GetUserDepositAssignedTime(mc *multicall.MultiCaller)
	GetUseLatestDelegate(mc *multicall.MultiCaller)
	GetDelegate(mc *multicall.MultiCaller)
	GetPreviousDelegate(mc *multicall.MultiCaller)
	GetEffectiveDelegate(mc *multicall.MultiCaller)

	// Transactions
	Refund(opts *bind.TransactOpts) (*core.TransactionInfo, error)
	Stake(validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, opts *bind.TransactOpts) (*core.TransactionInfo, error)
	Dissolve(opts *bind.TransactOpts) (*core.TransactionInfo, error)
	Close(opts *bind.TransactOpts) (*core.TransactionInfo, error)
	Finalise(opts *bind.TransactOpts) (*core.TransactionInfo, error)
	DelegateUpgrade(opts *bind.TransactOpts) (*core.TransactionInfo, error)
	DelegateRollback(opts *bind.TransactOpts) (*core.TransactionInfo, error)
	SetUseLatestDelegate(setting bool, opts *bind.TransactOpts) (*core.TransactionInfo, error)
	VoteScrub(opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Utils
	CalculateNodeShare(mc *multicall.MultiCaller, share_Out **big.Int, balance *big.Int)
	CalculateUserShare(mc *multicall.MultiCaller, share_Out **big.Int, balance *big.Int)
	GetPrestakeEvent(intervalSize *big.Int, opts *bind.CallOpts) (PrestakeData, error)
}

// ===============
// === Structs ===
// ===============

// Creates a new Minipool instance
/*
func NewMinipool(mgr *MinipoolManager, address common.Address, version uint8) (Minipool, error) {
	// Get the contract version
	version, err := rocketpool.GetContractVersion(rp, address, opts)
	if err != nil {
		errMsg := err.Error()
		errMsg = strings.ToLower(errMsg)
		if strings.Contains(errMsg, "execution reverted") ||
			strings.Contains(errMsg, "vm execution error") {
			// Reversions happen for minipool v1 on Prater which didn't have version() yet
			version = 1
		} else {
			return nil, fmt.Errorf("error getting minipool contract version: %w", err)
		}
	}

	switch version {
	case 1, 2:
		return newMinipool_v2(rp, address)
	case 3:
		return newMinipool_v3(rp, address, opts)
	default:
		return nil, fmt.Errorf("unexpected minipool contract version [%d]", version)
	}
}
*/

// Create a minipool binding from an explicit version number
func NewMinipoolFromVersion(rp *rocketpool.RocketPool, address common.Address, version uint8) (Minipool, error) {
	switch version {
	case 1, 2:
		return newMinipool_v2(rp, address)
	case 3:
		return newMinipool_v3(rp, address)
	default:
		return nil, fmt.Errorf("unexpected minipool contract version [%d]", version)
	}
}

// Create a minipool contract directly from its ABI, encoded in string form
func createMinipoolContractFromEncodedAbi(rp *rocketpool.RocketPool, address common.Address, encodedAbi string) (*core.Contract, error) {
	// Decode ABI
	abi, err := core.DecodeAbi(encodedAbi)
	if err != nil {
		return nil, fmt.Errorf("Could not decode minipool %s ABI: %w", address, err)
	}

	// Create and return
	return &core.Contract{
		Contract: bind.NewBoundContract(address, *abi, rp.Client, rp.Client, rp.Client),
		Address:  &address,
		ABI:      abi,
		Client:   rp.Client,
	}, nil
}

// Create a minipool contract directly from its ABI
func createMinipoolContractFromAbi(rp *rocketpool.RocketPool, address common.Address, abi *abi.ABI) (*core.Contract, error) {
	// Create and return
	return &core.Contract{
		Contract: bind.NewBoundContract(address, *abi, rp.Client, rp.Client, rp.Client),
		Address:  &address,
		ABI:      abi,
		Client:   rp.Client,
	}, nil
}

// =============
// === Calls ===
// =============

// Get all minipool details
func GetMinipools(rp *rocketpool.RocketPool, opts *bind.CallOpts) ([]MinipoolDetails, error) {
	minipoolAddresses, err := GetMinipoolAddresses(rp, opts)
	if err != nil {
		return []MinipoolDetails{}, err
	}
	return loadMinipoolDetails(rp, minipoolAddresses, opts)
}

// Get a node's minipool details
func GetNodeMinipools(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) ([]MinipoolDetails, error) {
	minipoolAddresses, err := GetNodeMinipoolAddresses(rp, nodeAddress, opts)
	if err != nil {
		return []MinipoolDetails{}, err
	}
	return loadMinipoolDetails(rp, minipoolAddresses, opts)
}

// Load minipool details
func loadMinipoolDetails(rp *rocketpool.RocketPool, minipoolAddresses []common.Address, opts *bind.CallOpts) ([]MinipoolDetails, error) {

	// Load minipool details in batches
	details := make([]MinipoolDetails, len(minipoolAddresses))
	for bsi := 0; bsi < len(minipoolAddresses); bsi += MinipoolDetailsBatchSize {

		// Get batch start & end index
		msi := bsi
		mei := bsi + MinipoolDetailsBatchSize
		if mei > len(minipoolAddresses) {
			mei = len(minipoolAddresses)
		}

		// Load details
		var wg errgroup.Group
		for mi := msi; mi < mei; mi++ {
			mi := mi
			wg.Go(func() error {
				minipoolAddress := minipoolAddresses[mi]
				minipoolDetails, err := GetMinipoolDetails(rp, minipoolAddress, opts)
				if err == nil {
					details[mi] = minipoolDetails
				}
				return err
			})
		}
		if err := wg.Wait(); err != nil {
			return []MinipoolDetails{}, err
		}

	}

	// Return
	return details, nil

}

// Get all minipool addresses
func GetMinipoolAddresses(rp *rocketpool.RocketPool, opts *bind.CallOpts) ([]common.Address, error) {

	// Get minipool count
	minipoolCount, err := GetMinipoolCount(rp, opts)
	if err != nil {
		return []common.Address{}, err
	}

	// Load minipool addresses in batches
	addresses := make([]common.Address, minipoolCount)
	for bsi := uint64(0); bsi < minipoolCount; bsi += MinipoolAddressBatchSize {

		// Get batch start & end index
		msi := bsi
		mei := bsi + MinipoolAddressBatchSize
		if mei > minipoolCount {
			mei = minipoolCount
		}

		// Load addresses
		var wg errgroup.Group
		for mi := msi; mi < mei; mi++ {
			mi := mi
			wg.Go(func() error {
				address, err := GetMinipoolAt(rp, mi, opts)
				if err == nil {
					addresses[mi] = address
				}
				return err
			})
		}
		if err := wg.Wait(); err != nil {
			return []common.Address{}, err
		}

	}

	// Return
	return addresses, nil

}

// Get the addresses of all minipools in prelaunch status
func GetPrelaunchMinipoolAddresses(rp *rocketpool.RocketPool, opts *bind.CallOpts) ([]common.Address, error) {

	rocketMinipoolManager, err := getRocketMinipoolManager(rp, opts)
	if err != nil {
		return []common.Address{}, err
	}

	// Get the total number of minipools
	totalMinipoolsUint, err := GetMinipoolCount(rp, nil)
	if err != nil {
		return []common.Address{}, err
	}

	totalMinipools := int64(totalMinipoolsUint)
	addresses := []common.Address{}
	limit := big.NewInt(MinipoolPrelaunchBatchSize)
	for i := int64(0); i < totalMinipools; i += MinipoolPrelaunchBatchSize {
		// Get a batch of addresses
		offset := big.NewInt(i)
		newAddresses := new([]common.Address)
		if err := rocketMinipoolManager.Call(opts, newAddresses, "getPrelaunchMinipools", offset, limit); err != nil {
			return []common.Address{}, fmt.Errorf("Could not get prelaunch minipool addresses: %w", err)
		}
		addresses = append(addresses, *newAddresses...)
	}

	return addresses, nil
}

// Get a node's minipool addresses
func GetNodeMinipoolAddresses(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) ([]common.Address, error) {

	// Get minipool count
	minipoolCount, err := GetNodeMinipoolCount(rp, nodeAddress, opts)
	if err != nil {
		return []common.Address{}, err
	}

	// Load minipool addresses in batches
	addresses := make([]common.Address, minipoolCount)
	for bsi := uint64(0); bsi < minipoolCount; bsi += MinipoolAddressBatchSize {

		// Get batch start & end index
		msi := bsi
		mei := bsi + MinipoolAddressBatchSize
		if mei > minipoolCount {
			mei = minipoolCount
		}

		// Load addresses
		var wg errgroup.Group
		for mi := msi; mi < mei; mi++ {
			mi := mi
			wg.Go(func() error {
				address, err := GetNodeMinipoolAt(rp, nodeAddress, mi, opts)
				if err == nil {
					addresses[mi] = address
				}
				return err
			})
		}
		if err := wg.Wait(); err != nil {
			return []common.Address{}, err
		}

	}

	// Return
	return addresses, nil

}

// Get a node's validating minipool pubkeys
func GetNodeValidatingMinipoolPubkeys(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) ([]rptypes.ValidatorPubkey, error) {

	// Get minipool count
	minipoolCount, err := GetNodeValidatingMinipoolCount(rp, nodeAddress, opts)
	if err != nil {
		return []rptypes.ValidatorPubkey{}, err
	}

	// Load pubkeys in batches
	var lock = sync.RWMutex{}
	pubkeys := make([]rptypes.ValidatorPubkey, minipoolCount)
	for bsi := uint64(0); bsi < minipoolCount; bsi += MinipoolAddressBatchSize {

		// Get batch start & end index
		msi := bsi
		mei := bsi + MinipoolAddressBatchSize
		if mei > minipoolCount {
			mei = minipoolCount
		}

		// Load pubkeys
		var wg errgroup.Group
		for mi := msi; mi < mei; mi++ {
			mi := mi
			wg.Go(func() error {
				minipoolAddress, err := GetNodeValidatingMinipoolAt(rp, nodeAddress, mi, opts)
				if err != nil {
					return err
				}
				pubkey, err := GetMinipoolPubkey(rp, minipoolAddress, opts)
				if err != nil {
					return err
				}
				lock.Lock()
				pubkeys[mi] = pubkey
				lock.Unlock()
				return nil
			})
		}
		if err := wg.Wait(); err != nil {
			return []rptypes.ValidatorPubkey{}, err
		}

	}

	// Return
	return pubkeys, nil

}

// Get a minipool's details
func GetMinipoolDetails(rp *rocketpool.RocketPool, minipoolAddress common.Address, opts *bind.CallOpts) (MinipoolDetails, error) {

	// Data
	var wg errgroup.Group
	var exists bool
	var pubkey rptypes.ValidatorPubkey

	// Load data
	wg.Go(func() error {
		var err error
		exists, err = GetMinipoolExists(rp, minipoolAddress, opts)
		return err
	})
	wg.Go(func() error {
		var err error
		pubkey, err = GetMinipoolPubkey(rp, minipoolAddress, opts)
		return err
	})

	// Wait for data
	if err := wg.Wait(); err != nil {
		return MinipoolDetails{}, err
	}

	// Return
	return MinipoolDetails{
		Address: minipoolAddress,
		Exists:  exists,
		Pubkey:  pubkey,
	}, nil

}

// Get the minipool count by status
func GetMinipoolCountPerStatus(rp *rocketpool.RocketPool, opts *bind.CallOpts) (MinipoolCountsPerStatus, error) {
	rocketMinipoolManager, err := getRocketMinipoolManager(rp, opts)
	if err != nil {
		return MinipoolCountsPerStatus{}, err
	}

	// Get the total number of minipools
	totalMinipoolsUint, err := GetMinipoolCount(rp, nil)
	if err != nil {
		return MinipoolCountsPerStatus{}, err
	}

	totalMinipools := int64(totalMinipoolsUint)
	minipoolCounts := MinipoolCountsPerStatus{
		Initialized:  big.NewInt(0),
		Prelaunch:    big.NewInt(0),
		Staking:      big.NewInt(0),
		Dissolved:    big.NewInt(0),
		Withdrawable: big.NewInt(0),
	}
	limit := big.NewInt(MinipoolPrelaunchBatchSize)
	for i := int64(0); i < totalMinipools; i += MinipoolPrelaunchBatchSize {
		// Get a batch of counts
		offset := big.NewInt(i)
		newMinipoolCounts := new(MinipoolCountsPerStatus)
		if err := rocketMinipoolManager.Call(opts, newMinipoolCounts, "getMinipoolCountPerStatus", offset, limit); err != nil {
			return MinipoolCountsPerStatus{}, fmt.Errorf("Could not get minipool counts: %w", err)
		}
		if newMinipoolCounts != nil {
			if newMinipoolCounts.Initialized != nil {
				minipoolCounts.Initialized.Add(minipoolCounts.Initialized, newMinipoolCounts.Initialized)
			}
			if newMinipoolCounts.Prelaunch != nil {
				minipoolCounts.Prelaunch.Add(minipoolCounts.Prelaunch, newMinipoolCounts.Prelaunch)
			}
			if newMinipoolCounts.Staking != nil {
				minipoolCounts.Staking.Add(minipoolCounts.Staking, newMinipoolCounts.Staking)
			}
			if newMinipoolCounts.Dissolved != nil {
				minipoolCounts.Dissolved.Add(minipoolCounts.Dissolved, newMinipoolCounts.Dissolved)
			}
			if newMinipoolCounts.Withdrawable != nil {
				minipoolCounts.Withdrawable.Add(minipoolCounts.Withdrawable, newMinipoolCounts.Withdrawable)
			}
		}
	}
	return minipoolCounts, nil
}

// Get a minipool address by index
func GetMinipoolAt(rp *rocketpool.RocketPool, index uint64, opts *bind.CallOpts) (common.Address, error) {
	rocketMinipoolManager, err := getRocketMinipoolManager(rp, opts)
	if err != nil {
		return common.Address{}, err
	}
	minipoolAddress := new(common.Address)
	if err := rocketMinipoolManager.Call(opts, minipoolAddress, "getMinipoolAt", big.NewInt(int64(index))); err != nil {
		return common.Address{}, fmt.Errorf("Could not get minipool %d address: %w", index, err)
	}
	return *minipoolAddress, nil
}

func GetNodeMinipoolAt(rp *rocketpool.RocketPool, nodeAddress common.Address, index uint64, opts *bind.CallOpts) (common.Address, error) {
	rocketMinipoolManager, err := getRocketMinipoolManager(rp, opts)
	if err != nil {
		return common.Address{}, err
	}
	minipoolAddress := new(common.Address)
	if err := rocketMinipoolManager.Call(opts, minipoolAddress, "getNodeMinipoolAt", nodeAddress, big.NewInt(int64(index))); err != nil {
		return common.Address{}, fmt.Errorf("Could not get node %s minipool %d address: %w", nodeAddress.Hex(), index, err)
	}
	return *minipoolAddress, nil
}

// Get a node's validating minipool count
func GetNodeValidatingMinipoolCount(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) (uint64, error) {
	rocketMinipoolManager, err := getRocketMinipoolManager(rp, opts)
	if err != nil {
		return 0, err
	}
	minipoolCount := new(*big.Int)
	if err := rocketMinipoolManager.Call(opts, minipoolCount, "getNodeValidatingMinipoolCount", nodeAddress); err != nil {
		return 0, fmt.Errorf("Could not get node %s validating minipool count: %w", nodeAddress.Hex(), err)
	}
	return (*minipoolCount).Uint64(), nil
}

// Get a node's validating minipool address by index
func GetNodeValidatingMinipoolAt(rp *rocketpool.RocketPool, nodeAddress common.Address, index uint64, opts *bind.CallOpts) (common.Address, error) {
	rocketMinipoolManager, err := getRocketMinipoolManager(rp, opts)
	if err != nil {
		return common.Address{}, err
	}
	minipoolAddress := new(common.Address)
	if err := rocketMinipoolManager.Call(opts, minipoolAddress, "getNodeValidatingMinipoolAt", nodeAddress, big.NewInt(int64(index))); err != nil {
		return common.Address{}, fmt.Errorf("Could not get node %s validating minipool %d address: %w", nodeAddress.Hex(), index, err)
	}
	return *minipoolAddress, nil
}

// Get a minipool address by validator pubkey
func GetMinipoolByPubkey(rp *rocketpool.RocketPool, pubkey rptypes.ValidatorPubkey, opts *bind.CallOpts) (common.Address, error) {
	rocketMinipoolManager, err := getRocketMinipoolManager(rp, opts)
	if err != nil {
		return common.Address{}, err
	}
	minipoolAddress := new(common.Address)
	if err := rocketMinipoolManager.Call(opts, minipoolAddress, "getMinipoolByPubkey", pubkey[:]); err != nil {
		return common.Address{}, fmt.Errorf("Could not get validator %s minipool address: %w", pubkey.Hex(), err)
	}
	return *minipoolAddress, nil
}

// Get a vacant minipool address by index
func GetVacantMinipoolAt(rp *rocketpool.RocketPool, index uint64, opts *bind.CallOpts) (common.Address, error) {
	rocketMinipoolManager, err := getRocketMinipoolManager(rp, opts)
	if err != nil {
		return common.Address{}, err
	}
	vacantMinipoolAddress := new(common.Address)
	if err := rocketMinipoolManager.Call(opts, vacantMinipoolAddress, "getVacantMinipoolAt", big.NewInt(int64(index))); err != nil {
		return common.Address{}, fmt.Errorf("Could not get vacant minipool %d address: %w", index, err)
	}
	return *vacantMinipoolAddress, nil
}

// Get contracts
var rocketMinipoolManagerLock sync.Mutex

func getRocketMinipoolManager(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*core.Contract, error) {
	rocketMinipoolManagerLock.Lock()
	defer rocketMinipoolManagerLock.Unlock()
	return rp.GetContract("rocketMinipoolManager", opts)
}
