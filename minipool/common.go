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
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
	rptypes "github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	EventScanInterval uint64 = 10000
)

// ===============
// === Structs ===
// ===============

// Basic binding for version-agnostic RocketMinipool contracts
type MinipoolCommon struct {
	Details  MinipoolCommonDetails
	Contract *core.Contract
	rp       *rocketpool.RocketPool
	mgr      *MinipoolManager
}

// Basic details about a minipool, version-agnostic
type MinipoolCommonDetails struct {
	// Core parameters
	Address                    common.Address                            `json:"address"`
	Version                    uint8                                     `json:"version"`
	NodeAddress                common.Address                            `json:"nodeAddress"`
	Status                     core.Uint8Parameter[types.MinipoolStatus] `json:"status"`
	StatusBlock                core.Parameter[uint64]                    `json:"statusBlock"`
	StatusTime                 core.Parameter[time.Time]                 `json:"statusTime"`
	IsFinalised                bool                                      `json:"isFinalized"`
	NodeFee                    core.Parameter[float64]                   `json:"nodeFee"`
	NodeDepositBalance         *big.Int                                  `json:"nodeDepositBalance"`
	NodeRefundBalance          *big.Int                                  `json:"nodeRefundBalance"`
	NodeDepositAssigned        bool                                      `json:"nodeDepositAssigned"`
	UserDepositBalance         *big.Int                                  `json:"userDepositBalance"`
	UserDepositAssigned        bool                                      `json:"userDepositAssigned"`
	UserDepositAssignedTime    core.Parameter[time.Time]                 `json:"userDepositAssignedTime"`
	IsUseLatestDelegateEnabled bool                                      `json:"IsUseLatestDelegateEnabled"`
	DelegateAddress            common.Address                            `json:"delegateAddress"`
	PreviousDelegateAddress    common.Address                            `json:"previousDelegateAddress"`
	EffectiveDelegateAddress   common.Address                            `json:"effectiveDelegateAddress"`
	PenaltyCount               core.Parameter[uint64]                    `json:"penaltyCount"`

	// MinipoolManager
	Exists                bool                                       `json:"exists"`
	Pubkey                types.ValidatorPubkey                      `json:"pubkey"`
	WithdrawalCredentials common.Hash                                `json:"withdrawalCredentials"`
	RplSlashed            bool                                       `json:"rplSlashed"`
	DepositType           core.Uint8Parameter[types.MinipoolDeposit] `json:"depositType"`

	// BondReducer
	IsBondReduceCancelled        bool                      `json:"isBondReduceCancelled"`
	ReduceBondTime               core.Parameter[time.Time] `json:"reduceBondTime"`
	ReduceBondValue              *big.Int                  `json:"reduceBondValue"`
	LastBondReductionTime        core.Parameter[time.Time] `json:"lastBondReductionTime"`
	LastBondReductionPrevValue   *big.Int                  `json:"lastBondReductionPrevValue"`
	LastBondReductionPrevNodeFee core.Parameter[float64]   `json:"lastBondReductionPrevNodeFee"`
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
func NewMinipoolCommonFromVersion(rp *rocketpool.RocketPool, contract *core.Contract, version uint8) (*MinipoolCommon, error) {
	mgr, err := NewMinipoolManager(rp) // TODO: get the latest instance instead of making a new one for memory reasons - maybe have RP register singletons
	if err != nil {
		return nil, fmt.Errorf("error creating minipool manager: %w", err)
	}

	return &MinipoolCommon{
		Details: MinipoolCommonDetails{
			Address: *contract.Address,
			Version: version,
		},
		rp:       rp,
		Contract: contract,
		mgr:      mgr,
	}, nil
}

// =============
// === Calls ===
// =============

// === Minipool ===

// Get the minipool's penalty count
func (c *MinipoolCommon) GetPenaltyCount(mc *multicall.MultiCaller) {
	// This isn't in the manager, it's in RocketStorage
	key := crypto.Keccak256Hash([]byte("network.penalties.penalty"), c.Details.Address.Bytes())
	c.rp.Storage.GetUint(mc, &c.Details.PenaltyCount.RawValue, key)
}

// Get the minipool's status
func (c *MinipoolCommon) GetStatus(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.Status.RawValue, "getStatus")
}

// Get the block that the minipool's status last changed
func (c *MinipoolCommon) GetStatusBlock(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.StatusBlock.RawValue, "getStatusBlock")
}

// Get the time that the minipool's status last changed
func (c *MinipoolCommon) GetStatusTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.StatusTime.RawValue, "getStatusTime")
}

// Check if the minipool has been finalised
func (c *MinipoolCommon) GetFinalised(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.IsFinalised, "getFinalised")
}

// Get the minipool's node address
func (c *MinipoolCommon) GetNodeAddress(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.NodeAddress, "getNodeAddress")
}

// Get the minipool's commission rate
func (c *MinipoolCommon) GetNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.NodeFee.RawValue, "getNodeFee")
}

// Get the balance the node has deposited to the minipool
func (c *MinipoolCommon) GetNodeDepositBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.NodeDepositBalance, "getNodeDepositBalance")
}

// Get the amount of ETH ready to be refunded to the node
func (c *MinipoolCommon) GetNodeRefundBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.NodeRefundBalance, "getNodeRefundBalance")
}

// Check if the node deposit has been assigned to the minipool
func (c *MinipoolCommon) GetNodeDepositAssigned(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.NodeDepositAssigned, "getNodeDepositAssigned")
}

// Get the balance the pool stakers have deposited to the minipool
func (c *MinipoolCommon) GetUserDepositBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.UserDepositBalance, "getUserDepositBalance")
}

// Check if the pool staker deposits has been assigned to the minipool
func (c *MinipoolCommon) GetUserDepositAssigned(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.UserDepositAssigned, "getUserDepositAssigned")
}

// Get the time at which the pool stakers were assigned to the minipool
func (c *MinipoolCommon) GetUserDepositAssignedTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.UserDepositAssignedTime.RawValue, "getUserDepositAssignedTime")
}

// Check if the "use latest delegate" flag is enabled
func (c *MinipoolCommon) GetUseLatestDelegate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.IsUseLatestDelegateEnabled, "getUseLatestDelegate")
}

// Get the address of the current delegate the minipool has recorded
func (c *MinipoolCommon) GetDelegate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.DelegateAddress, "getDelegate")
}

// Get the address of the previous delegate the minipool will use after a rollback
func (c *MinipoolCommon) GetPreviousDelegate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.PreviousDelegateAddress, "getPreviousDelegate")
}

// Get the address of the delegate the minipool will use (may be different than DelegateAddress if UseLatestDelegate is enabled)
func (c *MinipoolCommon) GetEffectiveDelegate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.Contract, &c.Details.EffectiveDelegateAddress, "getEffectiveDelegate")
}

// === Minipool Manager ===

// Check if a minipool exists
func (c *MinipoolCommon) GetExists(mc *multicall.MultiCaller) {
	// TODO: Is this really necessary?
	multicall.AddCall(mc, c.mgr.Contract, &c.Details.Exists, "getMinipoolExists", c.Details.Address)
}

// Get the minipool's pubkey
func (c *MinipoolCommon) GetPubkey(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.Contract, &c.Details.Pubkey, "getMinipoolPubkey", c.Details.Address)
}

// Get the minipool's 0x01-based withdrawal credentials
func (c *MinipoolCommon) GetWithdrawalCredentials(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.Contract, &c.Details.WithdrawalCredentials, "getMinipoolWithdrawalCredentials", c.Details.Address)
}

// Check if the minipool's RPL has been slashed
func (c *MinipoolCommon) GetRplSlashed(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.Contract, &c.Details.WithdrawalCredentials, "getMinipoolRPLSlashed", c.Details.Address)
}

// Get the minipool's deposit type
func (c *MinipoolCommon) GetDepositType(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.mgr.Contract, &c.Details.DepositType.RawValue, "getMinipoolDepositType", c.Details.Address)
}

// Get the basic details
func (c *MinipoolCommon) QueryAllDetails(mc *multicall.MultiCaller) {
	c.GetPenaltyCount(mc)
	c.GetStatus(mc)
	c.GetStatusBlock(mc)
	c.GetStatusTime(mc)
	c.GetFinalised(mc)
	c.GetNodeAddress(mc)
	c.GetNodeFee(mc)
	c.GetNodeDepositBalance(mc)
	c.GetNodeRefundBalance(mc)
	c.GetNodeDepositAssigned(mc)
	c.GetUserDepositBalance(mc)
	c.GetUserDepositAssigned(mc)
	c.GetUserDepositAssignedTime(mc)
	c.GetUseLatestDelegate(mc)
	c.GetDelegate(mc)
	c.GetPreviousDelegate(mc)
	c.GetEffectiveDelegate(mc)
	c.GetExists(mc)
	c.GetPubkey(mc)
	c.GetWithdrawalCredentials(mc)
	c.GetRplSlashed(mc)
	c.GetDepositType(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for refunding node ETH from the minipool
func (c *MinipoolCommon) Refund(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "refund", opts)
}

// Get info for progressing the prelaunch minipool to staking
func (c *MinipoolCommon) Stake(validatorSignature rptypes.ValidatorSignature, depositDataRoot common.Hash, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "stake", opts, validatorSignature[:], depositDataRoot)
}

// Get info for dissolving the initialized or prelaunch minipool
func (c *MinipoolCommon) Dissolve(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "dissolve", opts)
}

// Get info for withdrawing node balances from the dissolved minipool and closing it
func (c *MinipoolCommon) Close(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "close", opts)
}

// Get info for finalising a minipool to get the RPL stake back
func (c *MinipoolCommon) Finalise(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "finalise", opts)
}

// Get info for upgrading this minipool to the latest network delegate contract
func (c *MinipoolCommon) DelegateUpgrade(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "delegateUpgrade", opts)
}

// Get info for rolling back to the previous delegate contract
func (c *MinipoolCommon) DelegateRollback(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "delegateRollback", opts)
}

// Get info for setting the UseLatestDelegate flag (if set to true, will automatically use the latest delegate contract)
func (c *MinipoolCommon) SetUseLatestDelegate(setting bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "setUseLatestDelegate", opts, setting)
}

// Get info for voting to scrub a minipool
func (c *MinipoolCommon) VoteScrub(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "voteScrub", opts)
}

// =============
// === Utils ===
// =============

// Get the common (version-agnostic) details of the minipool
func (c *MinipoolCommon) GetCommonDetails() *MinipoolCommonDetails {
	return &c.Details
}

// Given a validator balance, calculates how much belongs to the node (taking into consideration rewards and penalties)
func (c *MinipoolCommon) CalculateNodeShare(mc *multicall.MultiCaller, share_Out **big.Int, balance *big.Int) {
	multicall.AddCall(mc, c.Contract, share_Out, "calculateNodeShare", balance)
}

// Given a validator balance, calculates how much belongs to rETH pool stakers (taking into consideration rewards and penalties)
func (c *MinipoolCommon) CalculateUserShare(mc *multicall.MultiCaller, share_Out **big.Int, balance *big.Int) {
	multicall.AddCall(mc, c.Contract, share_Out, "calculateUserShare", balance)
}

// Get the data from this minipool's MinipoolPrestaked event
func (c *MinipoolCommon) GetPrestakeEvent(intervalSize *big.Int, opts *bind.CallOpts) (PrestakeData, error) {

	addressFilter := []common.Address{c.Details.Address}
	topicFilter := [][]common.Hash{{c.Contract.ABI.Events["MinipoolPrestaked"].ID}}

	// Grab the latest block number
	currentBlock, err := c.rp.Client.BlockNumber(context.Background())
	if err != nil {
		return PrestakeData{}, fmt.Errorf("Error getting current block %s: %w", c.Details.Address.Hex(), err)
	}

	// Grab the lowest block number worth querying from (should never have to go back this far in practice)
	deployBlockHash := crypto.Keccak256Hash([]byte("deploy.block"))
	var fromBlockBig *big.Int
	err = c.rp.Query(func(mc *multicall.MultiCaller) error {
		c.rp.Storage.GetUint(mc, &fromBlockBig, deployBlockHash)
		return nil
	}, opts)
	if err != nil {
		return PrestakeData{}, fmt.Errorf("Error getting deploy block %s: %w", c.Details.Address.Hex(), err)
	}

	fromBlock := fromBlockBig.Uint64()
	var log gethtypes.Log
	found := false

	// Backwards scan through blocks to find the event
	for i := currentBlock; i >= fromBlock; i -= EventScanInterval {
		from := i - EventScanInterval + 1
		if from < fromBlock {
			from = fromBlock
		}

		fromBig := big.NewInt(0).SetUint64(from)
		toBig := big.NewInt(0).SetUint64(i)

		logs, err := utils.GetLogs(c.rp, addressFilter, topicFilter, intervalSize, fromBig, toBig, nil)
		if err != nil {
			return PrestakeData{}, fmt.Errorf("Error getting prestake logs for minipool %s: %w", c.Details.Address.Hex(), err)
		}

		if len(logs) > 0 {
			log = logs[0]
			found = true
			break
		}
	}

	if !found {
		// This should never happen
		return PrestakeData{}, fmt.Errorf("Error finding prestake log for minipool %s", c.Details.Address.Hex())
	}

	// Decode the event
	prestakeEvent := new(MinipoolPrestakeEvent)
	c.Contract.Contract.UnpackLog(prestakeEvent, "MinipoolPrestaked", log)
	if err != nil {
		return PrestakeData{}, fmt.Errorf("Error unpacking prestake data: %w", err)
	}

	// Convert the event to a more useable struct
	prestakeData := PrestakeData{
		Pubkey:                rptypes.BytesToValidatorPubkey(prestakeEvent.Pubkey),
		WithdrawalCredentials: common.BytesToHash(prestakeEvent.WithdrawalCredentials),
		Amount:                prestakeEvent.Amount,
		Signature:             rptypes.BytesToValidatorSignature(prestakeEvent.Signature),
		DepositDataRoot:       prestakeEvent.DepositDataRoot,
		Time:                  time.Unix(prestakeEvent.Time.Int64(), 0),
	}
	return prestakeData, nil
}
