package minipool

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils"
)

const (
	eventScanInterval uint64 = 10000
)

// ===============
// === Structs ===
// ===============

// Basic binding for version-agnostic RocketMinipool contracts
type MinipoolCommon struct {
	// The address of the minipool contract
	Address common.Address

	// The version of the minipool
	Version uint8

	// The address of the node that owns this minipool
	NodeAddress *core.SimpleField[common.Address]

	// The minipool's status
	Status *core.FormattedUint8Field[types.MinipoolStatus]

	// The block that the minipool's status last changed
	StatusBlock *core.FormattedUint256Field[uint64]

	// The time that the minipool's status last changed
	StatusTime *core.FormattedUint256Field[time.Time]

	// True if the minipool has been finalised
	IsFinalised *core.SimpleField[bool]

	// The minipool's commission rate
	NodeFee *core.FormattedUint256Field[float64]

	// The balance the node has deposited to the minipool
	NodeDepositBalance *core.SimpleField[*big.Int]

	// The amount of ETH ready to be refunded to the node
	NodeRefundBalance *core.SimpleField[*big.Int]

	// True if the node deposit has been assigned to the minipool
	NodeDepositAssigned *core.SimpleField[bool]

	// The balance the pool stakers have deposited to the minipool
	UserDepositBalance *core.SimpleField[*big.Int]

	// True if the pool staker deposits has been assigned to the minipool
	UserDepositAssigned *core.SimpleField[bool]

	// The time at which the pool stakers were assigned to the minipool
	UserDepositAssignedTime *core.FormattedUint256Field[time.Time]

	// True if the "use latest delegate" flag is enabled
	IsUseLatestDelegateEnabled *core.SimpleField[bool]

	// The address of the current delegate the minipool has recorded
	DelegateAddress *core.SimpleField[common.Address]

	// The address of the previous delegate the minipool will use after a rollback
	PreviousDelegateAddress *core.SimpleField[common.Address]

	// The address of the delegate the minipool will use (may be different than DelegateAddress if UseLatestDelegate is enabled)
	EffectiveDelegateAddress *core.SimpleField[common.Address]

	// The minipool's penalty count
	PenaltyCount *core.FormattedUint256Field[uint64]

	// True if a minipool exists (i.e. there is a minipool with this contract address)
	Exists *core.SimpleField[bool]

	// The pubkey of the validator on the Beacon Chain managed by this minipool
	Pubkey *core.SimpleField[types.ValidatorPubkey]

	// The minipool's 0x01-based withdrawal credentials
	WithdrawalCredentials *core.SimpleField[common.Hash]

	// True if the minipool's RPL has been slashed
	RplSlashed *core.SimpleField[bool]

	// The minipool's deposit type
	DepositType *core.FormattedUint8Field[types.MinipoolDeposit]

	// The queue position of the minipool (-1 means not in the queue, otherwise 0-indexed)
	QueuePosition *core.FormattedUint256Field[int64]

	// === Internal fields ===
	contract *core.Contract
	rp       *rocketpool.RocketPool
	mpMgr    *core.Contract
	mpQueue  *core.Contract
	mpStatus *core.Contract
}

// The data from a minipool's MinipoolPrestaked event
type MinipoolPrestakeEvent struct {
	Pubkey                []byte   `abi:"validatorPubkey"`
	Signature             []byte   `abi:"validatorSignature"`
	DepositDataRoot       [32]byte `abi:"depositDataRoot"`
	Amount                *big.Int `abi:"amount"`
	WithdrawalCredentials []byte   `abi:"withdrawalCredentials"`
	Time                  *big.Int `abi:"time"`
}

// Formatted MinipoolPrestaked event data
type PrestakeData struct {
	Pubkey                types.ValidatorPubkey    `json:"pubkey"`
	WithdrawalCredentials common.Hash              `json:"withdrawalCredentials"`
	Amount                *big.Int                 `json:"amount"`
	Signature             types.ValidatorSignature `json:"signature"`
	DepositDataRoot       common.Hash              `json:"depositDataRoot"`
	Time                  time.Time                `json:"time"`
}

// ====================
// === Constructors ===
// ====================

// Create a minipool common binding from an explicit version number
func newMinipoolCommonFromVersion(rp *rocketpool.RocketPool, contract *core.Contract, version uint8) (*MinipoolCommon, error) {
	mpMgr, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolManager)
	if err != nil {
		return nil, fmt.Errorf("error creating minipool manager: %w", err)
	}

	mpQueue, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolQueue)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool queue contract: %w", err)
	}

	mpStatus, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolStatus)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool status contract: %w", err)
	}

	address := *contract.Address
	penaltyCountKey := crypto.Keccak256Hash([]byte("network.penalties.penalty"), address.Bytes())
	return &MinipoolCommon{
		Address: address,
		Version: version,

		// Minipool
		PenaltyCount:               core.NewFormattedUint256Field[uint64](rp.Storage.Contract, "getUint", penaltyCountKey),
		Status:                     core.NewFormattedUint8Field[types.MinipoolStatus](contract, "getStatus"),
		StatusBlock:                core.NewFormattedUint256Field[uint64](contract, "getStatusBlock"),
		StatusTime:                 core.NewFormattedUint256Field[time.Time](contract, "getStatusTime"),
		IsFinalised:                core.NewSimpleField[bool](contract, "getFinalised"),
		NodeAddress:                core.NewSimpleField[common.Address](contract, "getNodeAddress"),
		NodeFee:                    core.NewFormattedUint256Field[float64](contract, "getNodeFee"),
		NodeDepositBalance:         core.NewSimpleField[*big.Int](contract, "getNodeDepositBalance"),
		NodeRefundBalance:          core.NewSimpleField[*big.Int](contract, "getNodeRefundBalance"),
		NodeDepositAssigned:        core.NewSimpleField[bool](contract, "getNodeDepositAssigned"),
		UserDepositBalance:         core.NewSimpleField[*big.Int](contract, "getUserDepositBalance"),
		UserDepositAssigned:        core.NewSimpleField[bool](contract, "getUserDepositAssigned"),
		UserDepositAssignedTime:    core.NewFormattedUint256Field[time.Time](contract, "getUserDepositAssignedTime"),
		IsUseLatestDelegateEnabled: core.NewSimpleField[bool](contract, "getUseLatestDelegate"),
		DelegateAddress:            core.NewSimpleField[common.Address](contract, "getDelegate"),
		PreviousDelegateAddress:    core.NewSimpleField[common.Address](contract, "getPreviousDelegate"),
		EffectiveDelegateAddress:   core.NewSimpleField[common.Address](contract, "getEffectiveDelegate"),

		// MinipoolManager
		Exists:                core.NewSimpleField[bool](mpMgr, "getMinipoolExists", address),
		Pubkey:                core.NewSimpleField[types.ValidatorPubkey](mpMgr, "getMinipoolPubkey", address),
		WithdrawalCredentials: core.NewSimpleField[common.Hash](mpMgr, "getMinipoolWithdrawalCredentials", address),
		RplSlashed:            core.NewSimpleField[bool](mpMgr, "getMinipoolRPLSlashed", address),
		DepositType:           core.NewFormattedUint8Field[types.MinipoolDeposit](mpMgr, "getMinipoolDepositType", address),

		// MinipoolQueue
		QueuePosition: core.NewFormattedUint256Field[int64](mpQueue, "getMinipoolPosition", address),

		rp:       rp,
		contract: contract,
		mpMgr:    mpMgr,
		mpQueue:  mpQueue,
		mpStatus: mpStatus,
	}, nil
}

// =============
// === Calls ===
// =============

// Gets the common binding for all minipool types
func (c *MinipoolCommon) Common() *MinipoolCommon {
	return c
}

// Gets the underlying minipool's contract
func (c *MinipoolCommon) GetContract() *core.Contract {
	return c.contract
}

// ====================
// === Transactions ===
// ====================

// === Minipool ===

// Get info for refunding node ETH from the minipool
func (c *MinipoolCommon) Refund(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "refund", opts)
}

// Get info for progressing the prelaunch minipool to staking
func (c *MinipoolCommon) Stake(validatorSignature types.ValidatorSignature, depositDataRoot common.Hash, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "stake", opts, validatorSignature[:], depositDataRoot)
}

// Get info for dissolving the initialized or prelaunch minipool
func (c *MinipoolCommon) Dissolve(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "dissolve", opts)
}

// Get info for withdrawing node balances from the dissolved minipool and closing it
func (c *MinipoolCommon) Close(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "close", opts)
}

// Get info for finalising a minipool to get the RPL stake back
func (c *MinipoolCommon) Finalise(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "finalise", opts)
}

// Get info for upgrading this minipool to the latest network delegate contract
func (c *MinipoolCommon) DelegateUpgrade(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "delegateUpgrade", opts)
}

// Get info for rolling back to the previous delegate contract
func (c *MinipoolCommon) DelegateRollback(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "delegateRollback", opts)
}

// Get info for setting the UseLatestDelegate flag (if set to true, will automatically use the latest delegate contract)
func (c *MinipoolCommon) SetUseLatestDelegate(setting bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "setUseLatestDelegate", opts, setting)
}

// Get info for voting to scrub a minipool
func (c *MinipoolCommon) VoteScrub(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "voteScrub", opts)
}

// === MinipoolStatus ===

// Get info for submitting a minipool withdrawable event
func (c *MinipoolCommon) SubmitMinipoolWithdrawable(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.mpStatus, "submitMinipoolWithdrawable", opts, c.Address)
}

// =============
// === Utils ===
// =============

// Given a validator balance, calculates how much belongs to the node (taking into consideration rewards and penalties)
func (c *MinipoolCommon) CalculateNodeShare(mc *batch.MultiCaller, share_Out **big.Int, balance *big.Int) {
	core.AddCall(mc, c.contract, share_Out, "calculateNodeShare", balance)
}

// Given a validator balance, calculates how much belongs to rETH pool stakers (taking into consideration rewards and penalties)
func (c *MinipoolCommon) CalculateUserShare(mc *batch.MultiCaller, share_Out **big.Int, balance *big.Int) {
	core.AddCall(mc, c.contract, share_Out, "calculateUserShare", balance)
}

// Get the data from this minipool's MinipoolPrestaked event
func (c *MinipoolCommon) GetPrestakeEvent(intervalSize *big.Int, opts *bind.CallOpts) (PrestakeData, error) {

	addressFilter := []common.Address{c.Address}
	topicFilter := [][]common.Hash{{c.contract.ABI.Events["MinipoolPrestaked"].ID}}

	// Grab the latest block number
	currentBlock, err := c.rp.Client.BlockNumber(context.Background())
	if err != nil {
		return PrestakeData{}, fmt.Errorf("error getting current block %s: %w", c.Address.Hex(), err)
	}

	// Grab the lowest block number worth querying from (should never have to go back this far in practice)
	deployBlockHash := crypto.Keccak256Hash([]byte("deploy.block"))
	var fromBlockBig *big.Int
	err = c.rp.Query(func(mc *batch.MultiCaller) error {
		c.rp.Storage.GetUint(mc, &fromBlockBig, deployBlockHash)
		return nil
	}, opts)
	if err != nil {
		return PrestakeData{}, fmt.Errorf("error getting deploy block %s: %w", c.Address.Hex(), err)
	}

	fromBlock := fromBlockBig.Uint64()
	var log gethtypes.Log
	found := false

	// Backwards scan through blocks to find the event
	for i := currentBlock; i >= fromBlock; i -= eventScanInterval {
		from := i - eventScanInterval + 1
		if from < fromBlock {
			from = fromBlock
		}

		fromBig := big.NewInt(0).SetUint64(from)
		toBig := big.NewInt(0).SetUint64(i)

		logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, fromBig, toBig, nil)
		if err != nil {
			return PrestakeData{}, fmt.Errorf("error getting prestake logs for minipool %s: %w", c.Address.Hex(), err)
		}

		if len(logs) > 0 {
			log = logs[0]
			found = true
			break
		}
	}

	if !found {
		// This should never happen
		return PrestakeData{}, fmt.Errorf("error finding prestake log for minipool %s", c.Address.Hex())
	}

	// Decode the event
	prestakeEvent := new(MinipoolPrestakeEvent)
	c.contract.Contract.UnpackLog(prestakeEvent, "MinipoolPrestaked", log)
	if err != nil {
		return PrestakeData{}, fmt.Errorf("error unpacking prestake data: %w", err)
	}

	// Convert the event to a more useable struct
	prestakeData := PrestakeData{
		Pubkey:                types.BytesToValidatorPubkey(prestakeEvent.Pubkey),
		WithdrawalCredentials: common.BytesToHash(prestakeEvent.WithdrawalCredentials),
		Amount:                prestakeEvent.Amount,
		Signature:             types.BytesToValidatorSignature(prestakeEvent.Signature),
		DepositDataRoot:       prestakeEvent.DepositDataRoot,
		Time:                  time.Unix(prestakeEvent.Time.Int64(), 0),
	}
	return prestakeData, nil
}
