package network

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// Estimate the gas of SetSnapshotAddress
func EstimateSetSnapshotAddress(rp *rocketpool.RocketPool, snapshotAddress common.Address, v uint8, r *[32]byte, s *[32]byte, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	RocketSignerRegistry, err := getRocketSignerRegistry(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	return RocketSignerRegistry.GetTransactionGasInfo(opts, "setSnapshotAddress", snapshotAddress, v, r, s)
}

// Set the snapshot address for the node
func SetSnapshotAddress(rp *rocketpool.RocketPool, snapshotAddress common.Address, v uint8, r *[32]byte, s *[32]byte, opts *bind.TransactOpts) (common.Hash, error) {
	RocketSignerRegistry, err := getRocketSignerRegistry(rp, nil)
	if err != nil {
		return common.Hash{}, err
	}
	tx, err := RocketSignerRegistry.Transact(opts, "setSnapshotAddress", snapshotAddress, v, r, s)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error setting voting delegate: %w", err)
	}
	return tx.Hash(), nil
}

// Get contracts
var rocketSignerRegistryLock sync.Mutex

func getRocketSignerRegistry(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*rocketpool.Contract, error) {
	rocketSignerRegistryLock.Lock()
	defer rocketSignerRegistryLock.Unlock()
	return rp.GetContract("rocketNetworkVoting", opts)
}
