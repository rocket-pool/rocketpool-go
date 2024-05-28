package registry

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// Estimate the gas of SetSnapshotAddress
func EstimateSetSnapshotAddress(rp *rocketpool.RocketPool, snapshotAddress common.Address, v uint8, r [32]byte, s [32]byte, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	return rp.RocketStorageContract.GetTransactionGasInfo(opts, "setSigningDelegate", snapshotAddress, v, r, s)
}

// Set the snapshot address for the node
func SetSnapshotAddress(rp *rocketpool.RocketPool, snapshotAddress common.Address, v uint8, r [32]byte, s [32]byte, opts *bind.TransactOpts) (common.Hash, error) {
	tx, err := rp.RocketSignerRegistry.SetSigningDelegate(opts, snapshotAddress, v, r, s)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error setting snapshot address: %w", err)
	}
	return tx.Hash(), nil
}
