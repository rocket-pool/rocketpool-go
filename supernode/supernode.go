package supernode

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	rptypes "github.com/rocket-pool/rocketpool-go/types"
)

// Estimate the gas of Deposit
func EstimateDepositGas(rp *rocketpool.RocketPool, supernodeAddress common.Address, validatorPubkey rptypes.ValidatorPubkey, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, salt *big.Int, expectedMinipoolAddress common.Address, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketSupernodeManager, err := getRocketSupernodeManager(rp)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	return rocketSupernodeManager.GetTransactionGasInfo(opts, "deposit", supernodeAddress, validatorPubkey[:], validatorSignature[:], depositDataRoot, salt, expectedMinipoolAddress)
}

// Make a node deposit
func Deposit(rp *rocketpool.RocketPool, supernodeAddress common.Address, validatorPubkey rptypes.ValidatorPubkey, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, salt *big.Int, expectedMinipoolAddress common.Address, opts *bind.TransactOpts) (common.Hash, error) {
	rocketSupernodeManager, err := getRocketSupernodeManager(rp)
	if err != nil {
		return common.Hash{}, err
	}
	hash, err := rocketSupernodeManager.Transact(opts, "deposit", supernodeAddress, validatorPubkey[:], validatorSignature[:], depositDataRoot, salt, expectedMinipoolAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("Could not make node deposit: %w", err)
	}
	return hash, nil
}

// Get contracts
var rocketSupernodeManagerLock sync.Mutex

func getRocketSupernodeManager(rp *rocketpool.RocketPool) (*rocketpool.Contract, error) {
	rocketSupernodeManagerLock.Lock()
	defer rocketSupernodeManagerLock.Unlock()
	return rp.GetContract("rocketSupernodeManager")
}
