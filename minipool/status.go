package minipool

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// Get the node reward amount for a minipool by node fee, user deposit balance, and staking start & end balances
func GetMinipoolNodeRewardAmount(rp *rocketpool.RocketPool, nodeFee float64, userDepositBalance, startBalance, endBalance *big.Int, opts *bind.CallOpts) (*big.Int, error) {
    rocketMinipoolStatus, err := getRocketMinipoolStatus(rp)
    if err != nil {
        return nil, err
    }
    nodeAmount := new(*big.Int)
    if err := rocketMinipoolStatus.Call(opts, nodeAmount, "getMinipoolNodeRewardAmount", eth.EthToWei(nodeFee), userDepositBalance, startBalance, endBalance); err != nil {
        return nil, fmt.Errorf("Could not get minipool node reward amount: %w", err)
    }
    return *nodeAmount, nil
}


// Estimate the gas of SubmitMinipoolWithdrawable
func EstimateSubmitMinipoolWithdrawableGas(rp *rocketpool.RocketPool, minipoolAddress common.Address, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
    rocketMinipoolStatus, err := getRocketMinipoolStatus(rp)
    if err != nil {
        return rocketpool.GasInfo{}, err
    }
    return rocketMinipoolStatus.GetTransactionGasInfo(opts, "submitMinipoolWithdrawable", minipoolAddress)
}


// Submit a minipool withdrawable event
func SubmitMinipoolWithdrawable(rp *rocketpool.RocketPool, minipoolAddress common.Address, opts *bind.TransactOpts) (common.Hash, error) {
    rocketMinipoolStatus, err := getRocketMinipoolStatus(rp)
    if err != nil {
        return common.Hash{}, err
    }
    hash, err := rocketMinipoolStatus.Transact(opts, "submitMinipoolWithdrawable", minipoolAddress)
    if err != nil {
        return common.Hash{}, fmt.Errorf("Could not submit minipool withdrawable event: %w", err)
    }
    return hash, nil
}


// Get contracts
var rocketMinipoolStatusLock sync.Mutex
func getRocketMinipoolStatus(rp *rocketpool.RocketPool) (*rocketpool.Contract, error) {
    rocketMinipoolStatusLock.Lock()
    defer rocketMinipoolStatusLock.Unlock()
    return rp.GetContract("rocketMinipoolStatus")
}

