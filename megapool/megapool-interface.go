package megapool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	rptypes "github.com/rocket-pool/rocketpool-go/types"
)

type Megapool interface {
	GetContract() *rocketpool.Contract
	GetAddress() common.Address
	GetVersion() uint8
	GetValidatorCount(opts *bind.CallOpts) (uint64, error)
	GetAssignedValue(opts *bind.CallOpts) (*big.Int, error)
	GetDebt(opts *bind.CallOpts) (*big.Int, error)
	GetRefundValue(opts *bind.CallOpts) (*big.Int, error)
	GetNodeCapital(opts *bind.CallOpts) (*big.Int, error)
	GetNodeBond(opts *bind.CallOpts) (*big.Int, error)
	GetUserCapital(opts *bind.CallOpts) (*big.Int, error)
	// CalculateRewards (not yet implemented)
	GetPendingRewards(opts *bind.CallOpts) (*big.Int, error)
	GetNodeAddress(opts *bind.CallOpts) (common.Address, error)
	EstimateNewValidatorGas(validatorId uint32, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, opts *bind.TransactOpts) (rocketpool.GasInfo, error)
	NewValidator(bondAmount *big.Int, useExpressTicket bool, validatorPubkey rptypes.ValidatorPubkey, validatorSignature rptypes.ValidatorSignature, opts *bind.TransactOpts) (common.Hash, error)
	EstimateDequeueGas(validatorId uint32, opts *bind.TransactOpts) (rocketpool.GasInfo, error)
	Dequeue(validatorId uint32, opts *bind.TransactOpts) (common.Hash, error)
	EstimateAssignFundsGas(validatorId uint32, opts *bind.TransactOpts) (rocketpool.GasInfo, error)
	AssignFunds(validatorId uint32, opts *bind.TransactOpts) (common.Hash, error)
	EstimateDissolveValidatorGas(validatorId uint32, opts *bind.TransactOpts) (rocketpool.GasInfo, error)
	DissolveValidator(validatorId uint32, opts *bind.TransactOpts) (common.Hash, error)
	EstimateRepayDebtGas(opts *bind.TransactOpts) (rocketpool.GasInfo, error)
	RepayDebt(opts *bind.TransactOpts) (common.Hash, error)
	GetWithdrawalCredentials(opts *bind.CallOpts) (common.Hash, error)
	EstimateRequestUnstakeRPL(opts *bind.TransactOpts) (rocketpool.GasInfo, error)
	RequestUnstakeRPL(opts *bind.TransactOpts) (common.Hash, error)
	EstimateStakeGas(validatorId uint32, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, opts *bind.TransactOpts) (rocketpool.GasInfo, error)
	Stake(validatorId uint32, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, validatorProof validatorProof, opts *bind.TransactOpts) (common.Hash, error)
}
