package minipool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/types"
)

// ==================
// === Interfaces ===
// ==================

type IMinipool interface {
	// Get all of the minipool's details
	QueryAllDetails(mc *batch.MultiCaller)

	// Gets the underlying minipool's contract
	GetContract() *core.Contract

	// Gets the common details for all minipool types
	GetCommonDetails() *MinipoolCommonDetails

	// Get the minipool's penalty count
	GetPenaltyCount(mc *batch.MultiCaller)

	// Get the minipool's status
	GetStatus(mc *batch.MultiCaller)

	// Get the block that the minipool's status last changed
	GetStatusBlock(mc *batch.MultiCaller)

	// Get the time that the minipool's status last changed
	GetStatusTime(mc *batch.MultiCaller)

	// Check if the minipool has been finalised
	GetFinalised(mc *batch.MultiCaller)

	// Get the minipool's node address
	GetNodeAddress(mc *batch.MultiCaller)

	// Get the minipool's commission rate
	GetNodeFee(mc *batch.MultiCaller)

	// Get the balance the node has deposited to the minipool
	GetNodeDepositBalance(mc *batch.MultiCaller)

	// Get the amount of ETH ready to be refunded to the node
	GetNodeRefundBalance(mc *batch.MultiCaller)

	// Check if the node deposit has been assigned to the minipool
	GetNodeDepositAssigned(mc *batch.MultiCaller)

	// Get the balance the pool stakers have deposited to the minipool
	GetUserDepositBalance(mc *batch.MultiCaller)

	// Check if the pool staker deposits has been assigned to the minipool
	GetUserDepositAssigned(mc *batch.MultiCaller)

	// Get the time at which the pool stakers were assigned to the minipool
	GetUserDepositAssignedTime(mc *batch.MultiCaller)

	// Check if the "use latest delegate" flag is enabled
	GetUseLatestDelegate(mc *batch.MultiCaller)

	// Get the address of the current delegate the minipool has recorded
	GetDelegate(mc *batch.MultiCaller)

	// Get the address of the previous delegate the minipool will use after a rollback
	GetPreviousDelegate(mc *batch.MultiCaller)

	// Get the address of the delegate the minipool will use (may be different than DelegateAddress if UseLatestDelegate is enabled)
	GetEffectiveDelegate(mc *batch.MultiCaller)

	// Check if a minipool exists
	GetExists(mc *batch.MultiCaller)

	// Get the minipool's pubkey
	GetPubkey(mc *batch.MultiCaller)

	// Get the minipool's 0x01-based withdrawal credentials
	GetWithdrawalCredentials(mc *batch.MultiCaller)

	// Check if the minipool's RPL has been slashed
	GetRplSlashed(mc *batch.MultiCaller)

	// Get the minipool's deposit type
	GetDepositType(mc *batch.MultiCaller)

	// Get queue position of the minipool (-1 means not in the queue, otherwise 0-indexed).
	GetQueuePosition(mc *batch.MultiCaller)

	// Get info for refunding node ETH from the minipool
	Refund(opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Get info for progressing the prelaunch minipool to staking
	Stake(validatorSignature types.ValidatorSignature, depositDataRoot common.Hash, opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Get info for dissolving the initialized or prelaunch minipool
	Dissolve(opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Get info for withdrawing node balances from the dissolved minipool and closing it
	Close(opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Get info for finalising a minipool to get the RPL stake back
	Finalise(opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Get info for upgrading this minipool to the latest network delegate contract
	DelegateUpgrade(opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Get info for rolling back to the previous delegate contract
	DelegateRollback(opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Get info for setting the UseLatestDelegate flag (if set to true, will automatically use the latest delegate contract)
	SetUseLatestDelegate(setting bool, opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Get info for voting to scrub a minipool
	VoteScrub(opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Get info for submitting a minipool withdrawable event
	SubmitMinipoolWithdrawable(opts *bind.TransactOpts) (*core.TransactionInfo, error)

	// Given a validator balance, calculates how much belongs to the node (taking into consideration rewards and penalties)
	CalculateNodeShare(mc *batch.MultiCaller, share_Out **big.Int, balance *big.Int)

	// Given a validator balance, calculates how much belongs to rETH pool stakers (taking into consideration rewards and penalties)
	CalculateUserShare(mc *batch.MultiCaller, share_Out **big.Int, balance *big.Int)

	// Get the data from this minipool's MinipoolPrestaked event
	GetPrestakeEvent(intervalSize *big.Int, opts *bind.CallOpts) (PrestakeData, error)
}
