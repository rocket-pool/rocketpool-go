package minipool

import (
	"fmt"
	"math/big"
	"sync"

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
	QueryAllDetails(mc *multicall.MultiCaller)
	GetMinipoolCommon() *MinipoolCommon
}

// ====================
// === Constructors ===
// ====================

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

// =============
// === Calls ===
// =============

// Get a node's minipool details
func GetNodeMinipools(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) ([]MinipoolDetails, error) {
	minipoolAddresses, err := GetNodeMinipoolAddresses(rp, nodeAddress, opts)
	if err != nil {
		return []MinipoolDetails{}, err
	}
	return loadMinipoolDetails(rp, minipoolAddresses, opts)
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

// Get contracts
var rocketMinipoolManagerLock sync.Mutex

func getRocketMinipoolManager(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*core.Contract, error) {
	rocketMinipoolManagerLock.Lock()
	defer rocketMinipoolManagerLock.Unlock()
	return rp.GetContract("rocketMinipoolManager", opts)
}
