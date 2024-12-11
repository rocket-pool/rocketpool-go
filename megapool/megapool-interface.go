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
	GetAssignedValue(opts *bind.CallOpts) (uint64, error)
	GetDebt(opts *bind.CallOpts) (uint64, error)
	GetRefundValue(opts *bind.CallOpts) (uint64, error)
	GetNodeCapital(opts *bind.CallOpts) (uint64, error)
	GetNodeBond(opts *bind.CallOpts) (uint64, error)
	GetUserCapital(opts *bind.CallOpts) (uint64, error)
	// CalculateRewards (not yet implemented)
	GetPendingRewards(opts *bind.CallOpts) (uint64, error)
	GetNodeAddress(opts *bind.CallOpts) (common.Address, error)
	// The functions below require gas estimators
	NewValidator(bondAmount *big.Int, useExpressTicket bool, validatorPubkey rptypes.ValidatorPubkey, validatorSignature rptypes.ValidatorSignature, opts *bind.TransactOpts) (common.Hash, error)
	Dequeue(validatorId uint32, opts *bind.TransactOpts) (common.Hash, error)
	AssignFunds(validatorId uint32, opts *bind.TransactOpts) (common.Hash, error)
	DissolveValidator(validatorId uint32, opts *bind.TransactOpts) (common.Hash, error)
	RepayDebt(opts *bind.TransactOpts) (common.Hash, error)
	GetWithdrawalCredentials(opts *bind.CallOpts) ([]byte, error)
	RequestUnstakeRPL(opts *bind.TransactOpts) (common.Hash, error)
	EstimateStakeGas(validatorId uint32, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, opts *bind.TransactOpts) (rocketpool.GasInfo, error)
	Stake(validatorId uint32, validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, validatorProof validatorProof, opts *bind.TransactOpts) (common.Hash, error)
}
