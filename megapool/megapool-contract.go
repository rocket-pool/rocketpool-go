package megapool

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	rptypes "github.com/rocket-pool/rocketpool-go/types"
)

type validatorProof struct {
	slot                  uint64
	validatorIndex        *big.Int
	pubkey                []byte
	withdrawalCredentials [32]byte
	witnesses             [][32]byte
}

type withdrawal struct {
	index                 *big.Int
	validatorIndex        *big.Int
	withdrawalCredentials [32]byte
	amountInGwei          *big.Int
}

type MegapoolV1 interface {
	Megapool
}

// Megapool contract
type megapoolV1 struct {
	Address    common.Address
	Version    uint8
	Contract   *rocketpool.Contract
	RocketPool *rocketpool.RocketPool
}

const (
	megapoolV1EncodedAbi string = "eJztWE1v2zAM/SuDz8EO3dZDb/0aUKwdsrTdDkVQ0BaTCbUlT6KSBkX/++jYie3Yjp02aXPoLYmopyeSeiRz9+RJFTuy3tFd8pHQKAhvZjF6R16gFRkI6NNABw9I16QNjPEiMRpBgF7PUxAlhvemaHAshEFreZlSHMh+eB72PEtAeOUIfBlKmvGq0iqGGfgh5jv4ZEvGBQzoPfeePGCjWaQd0xxBaLFXZi3wEYV3xDvmK6VLwJJNRjbCMcRahzX8em1gjr8ffDvMwSYQSgFM80LkeAurEt6S93pAkhHWIA2XBlcZ/d+Lk4+tlWOFhfNxgoo+/NbitzP859B9+G1Tv52rZr81a0mC/OWgIBmNN2GrwvEwz+7vTonEpdpRdkCtlFR0ZORUQFKrMrkcPQhBRu24dRLVBi0wNhgw0jbht+ZWkSb/fnKT1upwgsuM21GAJJcY6bvtRijHHyMtpJlvsuLrWmeVHmb9o6xSm0icbkTqDH3aFy7nj7E0kBidhNxC7Autn1oUupg1lCpi37XpeRGlE63EPrnoFGJJEO4LpT4qIdV4gFMwoiVwb8dqgGy0Xwpwa9HsWeyWWn+qnXpPedpaFSte6kKN9MqdAh3FWnHfVHuoPyMsqErs/B84y09Ll1Pea8mGYLMiFPFRN5VeL6HcDWWQNAyW8lr2OpRMylpB/LRxXjSDHLBJ4fD5ausufJRUbFe77ZLqV9Yhbbitz7JP8PACmnFSMPhxbsx10TFVNg4rA0A6VmcT/aK3zwb3S5hxin5eTdvKKyMXcyO0pZf/R9JfYWAK4alBwVkqeXBZLwDz9C+m1srLyB7iFgWgJDf3PifvcZQJVe0MtiZW987ieRrqG5kEoUOwy2qQK1Df+Q+ddKER4ZqlAcgZ7ApSkkKecbSVdAYEA62pPgyZscLpjht6g7yl0t2+ZlDskAvQlAdFXnPRu1VzWRj0L3c+dZUpitQlzQRtzVDwvsNgewbbHSdur61EJ1QPv+aYNtQrXubVhrLY8JdO8tdNN0WptAcbyEDREdMG8a3xRj3U3bAEpljVsArAVs2laKkKfaP1qBCkOPteqToLj2dV9i3fUofnvotXnqO7nahIjj9BY5Plt+3Ah/8BcMhg6g=="
)

// The decoded ABI for megapools
var megapoolV1Abi *abi.ABI

// Create new minipool contract
func NewMegaPoolV1(rp *rocketpool.RocketPool, address common.Address, opts *bind.CallOpts) (Megapool, error) {

	var contract *rocketpool.Contract
	var err error
	if megapoolV1Abi == nil {
		// Get contract
		contract, err = createMegapoolContractFromEncodedAbi(rp, address, megapoolV1EncodedAbi)
	} else {
		contract, err = createMegapoolContractFromAbi(rp, address, megapoolV1Abi)
	}
	if err != nil {
		return nil, err
	} else if megapoolV1Abi == nil {
		megapoolV1Abi = contract.ABI
	}

	// Create and return
	return &megapoolV1{
		Address:    address,
		Version:    1,
		Contract:   contract,
		RocketPool: rp,
	}, nil
}

// Get the contract
func (mp *megapoolV1) GetContract() *rocketpool.Contract {
	return mp.Contract
}

// Get the contract address
func (mp *megapoolV1) GetAddress() common.Address {
	return mp.Address
}

// Get the contract version
func (mp *megapoolV1) GetVersion() uint8 {
	return mp.Version
}

func (mp *megapoolV1) GetValidatorCount(opts *bind.CallOpts) (uint64, error) {
	validatorCount := new(*big.Int)
	if err := mp.Contract.Call(opts, validatorCount, "getValidatorCount"); err != nil {
		return 0, fmt.Errorf("error getting megapool %s validator count: %w", mp.Address.Hex(), err)
	}
	return (*validatorCount).Uint64(), nil
}

//TODO func GetValidatorInfo()

func (mp *megapoolV1) GetAssignedValue(opts *bind.CallOpts) (*big.Int, error) {
	assignedValue := new(*big.Int)
	if err := mp.Contract.Call(opts, assignedValue, "getAssignedValue"); err != nil {
		return nil, fmt.Errorf("error getting megapool %s assigned value: %w", mp.Address.Hex(), err)
	}
	return *assignedValue, nil
}

func (mp *megapoolV1) GetDebt(opts *bind.CallOpts) (*big.Int, error) {
	debt := new(*big.Int)
	if err := mp.Contract.Call(opts, debt, "getDebt"); err != nil {
		return nil, fmt.Errorf("error getting megapool %s debt: %w", mp.Address.Hex(), err)
	}
	return *debt, nil
}

func (mp *megapoolV1) GetRefundValue(opts *bind.CallOpts) (*big.Int, error) {
	refundValue := new(*big.Int)
	if err := mp.Contract.Call(opts, refundValue, "getRefundValue"); err != nil {
		return nil, fmt.Errorf("error getting megapool %s refund value: %w", mp.Address.Hex(), err)
	}
	return *refundValue, nil
}

func (mp *megapoolV1) GetNodeCapital(opts *bind.CallOpts) (*big.Int, error) {
	nodeCapital := new(*big.Int)
	if err := mp.Contract.Call(opts, nodeCapital, "getNodeCapital"); err != nil {
		return nil, fmt.Errorf("error getting megapool %s node capital: %w", mp.Address.Hex(), err)
	}
	return *nodeCapital, nil
}

func (mp *megapoolV1) GetNodeBond(opts *bind.CallOpts) (*big.Int, error) {
	nodeBond := new(*big.Int)
	if err := mp.Contract.Call(opts, nodeBond, "getNodeBond"); err != nil {
		return nil, fmt.Errorf("error getting megapool %s debt: %w", mp.Address.Hex(), err)
	}
	return *nodeBond, nil
}

func (mp *megapoolV1) GetUserCapital(opts *bind.CallOpts) (*big.Int, error) {
	userCapital := new(*big.Int)
	if err := mp.Contract.Call(opts, userCapital, "getUserCapital"); err != nil {
		return nil, fmt.Errorf("error getting megapool %s user capital: %w", mp.Address.Hex(), err)
	}
	return *userCapital, nil
}

//TODO _calculateRewards is currently a view in RocketMegapoolDelegate.sol

func (mp *megapoolV1) GetPendingRewards(opts *bind.CallOpts) (*big.Int, error) {
	pendingRewards := new(*big.Int)
	if err := mp.Contract.Call(opts, pendingRewards, "getPendingRewards"); err != nil {
		return nil, fmt.Errorf("error getting megapool %s pending rewards: %w", mp.Address.Hex(), err)
	}
	return *pendingRewards, nil
}

func (mp *megapoolV1) GetNodeAddress(opts *bind.CallOpts) (common.Address, error) {
	nodeAddress := new(common.Address)
	if err := mp.Contract.Call(opts, nodeAddress, "getNodeAddress"); err != nil {
		return common.Address{}, fmt.Errorf("error getting megapool %s node address: %w", mp.Address.Hex(), err)
	}
	return *nodeAddress, nil
}

// Estimate the gas required to create a new validator as part of a megapool
func (mp *megapoolV1) EstimateNewValidatorGas(validatorId uint32, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	return mp.Contract.GetTransactionGasInfo(opts, "newValidator", validatorId, validatorSignature[:], depositDataRoot)
}

// Create a new validator as part of a megapool
func (mp *megapoolV1) NewValidator(bondAmount *big.Int, useExpressTicket bool, validatorPubkey rptypes.ValidatorPubkey, validatorSignature rptypes.ValidatorSignature, opts *bind.TransactOpts) (common.Hash, error) {
	tx, err := mp.Contract.Transact(opts, "newValidator", bondAmount, useExpressTicket, validatorPubkey[:], validatorSignature[:])
	if err != nil {
		return common.Hash{}, fmt.Errorf("error creating new validator %s: %w", validatorPubkey.Hex(), err)
	}
	return tx.Hash(), nil
}

// Estimate the gas required to remove a validator from the deposit queue
func (mp *megapoolV1) EstimateDequeueGas(validatorId uint32, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	return mp.Contract.GetTransactionGasInfo(opts, "dequeue", validatorId)
}

// Remove a validator from the deposit queue
func (mp *megapoolV1) Dequeue(validatorId uint32, opts *bind.TransactOpts) (common.Hash, error) {
	tx, err := mp.Contract.Transact(opts, "dequeue", validatorId)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error dequeuing validator ID %d: %w", validatorId, err)
	}
	return tx.Hash(), nil
}

// Estimate the gas required to accept requested funds from the deposit pool
func (mp *megapoolV1) EstimateAssignFundsGas(validatorId uint32, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	return mp.Contract.GetTransactionGasInfo(opts, "assignFunds", validatorId)
}

// Accept requested funds from the deposit pool
func (mp *megapoolV1) AssignFunds(validatorId uint32, opts *bind.TransactOpts) (common.Hash, error) {
	tx, err := mp.Contract.Transact(opts, "assignFunds", validatorId)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error assigning funds to validator ID %d: %w", validatorId, err)
	}
	return tx.Hash(), nil
}

// Estimate the gas required to dissolve a validator that has not staked within the required period
func (mp *megapoolV1) EstimateDissolveValidatorGas(validatorId uint32, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	return mp.Contract.GetTransactionGasInfo(opts, "dissolveValidator", validatorId)
}

// Dissolve a validator that has not staked within the required period
func (mp *megapoolV1) DissolveValidator(validatorId uint32, opts *bind.TransactOpts) (common.Hash, error) {
	tx, err := mp.Contract.Transact(opts, "dissolveValidator", validatorId)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error dissolving validator ID %d: %w", validatorId, err)
	}
	return tx.Hash(), nil
}

// Estimate the gas required to repay megapool debt
func (mp *megapoolV1) EstimateRepayDebtGas(opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	return mp.Contract.GetTransactionGasInfo(opts, "repayDebt")
}

// Receive ETH, which is sent to the rETH contract, to repay a megapool debt
func (mp *megapoolV1) RepayDebt(opts *bind.TransactOpts) (common.Hash, error) {
	tx, err := mp.Contract.Transact(opts, "repayDebt")
	if err != nil {
		return common.Hash{}, fmt.Errorf("error repaying debt for megapool %s: %w", mp.Address.Hex(), err)
	}
	return tx.Hash(), nil
}

// Get the expected withdrawal credentials for any validator within this megapool
func (mp *megapoolV1) GetWithdrawalCredentials(opts *bind.CallOpts) (common.Hash, error) {
	withdrawalCredentials := new(common.Hash)
	if err := mp.Contract.Call(opts, withdrawalCredentials, "getWithdrawalCredentials"); err != nil {
		return common.Hash{}, fmt.Errorf("error getting megapool %s withdrawal credentials: %w", mp.Address.Hex(), err)
	}
	return *withdrawalCredentials, nil
}

// Estimate the gas required to Request RPL previously staked on this megapool to be unstaked
func (mp *megapoolV1) EstimateRequestUnstakeRPL(opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	return mp.Contract.GetTransactionGasInfo(opts, "requestUnstakeRPL")
}

// RequestUnstakeRPL is not yet implemented in RocketMegapoolDelegate.sol
// Request RPL previously staked on this megapool to be unstaked
func (mp *megapoolV1) RequestUnstakeRPL(opts *bind.TransactOpts) (common.Hash, error) {
	tx, err := mp.Contract.Transact(opts, "requestUnstakeRPL")
	if err != nil {
		return common.Hash{}, fmt.Errorf("error requesting unstake rpl for megapool %s: %w", mp.Address.Hex(), err)
	}
	return tx.Hash(), nil
}

// Estimate the gas of Stake
func (mp *megapoolV1) EstimateStakeGas(validatorId uint32, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	return mp.Contract.GetTransactionGasInfo(opts, "stake", validatorId, validatorSignature[:], depositDataRoot)
}

// Progress the prelaunch megapool to staking
func (mp *megapoolV1) Stake(validatorId uint32, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, validatorProof validatorProof, opts *bind.TransactOpts) (common.Hash, error) {
	tx, err := mp.Contract.Transact(opts, "stake", validatorSignature[:], depositDataRoot, validatorProof)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error staking megapool %s: %w", mp.Address.Hex(), err)
	}
	return tx.Hash(), nil
}

// Create a megapool contract directly from its ABI
func createMegapoolContractFromAbi(rp *rocketpool.RocketPool, address common.Address, abi *abi.ABI) (*rocketpool.Contract, error) {
	// Create and return
	return &rocketpool.Contract{
		Contract: bind.NewBoundContract(address, *abi, rp.Client, rp.Client, rp.Client),
		Address:  &address,
		ABI:      abi,
		Client:   rp.Client,
	}, nil
}

// Create a megapool contract directly from its ABI, encoded in string form
func createMegapoolContractFromEncodedAbi(rp *rocketpool.RocketPool, address common.Address, encodedAbi string) (*rocketpool.Contract, error) {
	// Decode ABI
	abi, err := rocketpool.DecodeAbi(encodedAbi)
	if err != nil {
		return nil, fmt.Errorf("error decoding megapool %s ABI: %w", address, err)
	}

	// Create and return
	return &rocketpool.Contract{
		Contract: bind.NewBoundContract(address, *abi, rp.Client, rp.Client, rp.Client),
		Address:  &address,
		ABI:      abi,
		Client:   rp.Client,
	}, nil
}