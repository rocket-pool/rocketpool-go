package minipool

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Minipool queue capacity
type QueueCapacity struct {
	Total     *big.Int
	Effective *big.Int
}

// Minipools queue status details
type QueueDetails struct {
	Position int64
}

// Binding for RocketMinipoolQueue
type MinipoolQueue struct {
	Details  MinipoolQueueDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for RocketMinipoolQueue
type MinipoolQueueDetails struct {
	TotalLength       core.Parameter[uint64] `json:"totalLength"`
	TotalCapacity     *big.Int               `json:"totalCapacity"`
	EffectiveCapacity *big.Int               `json:"effectiveCapacity"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new MinipoolQueue contract binding
func NewMinipoolQueue(rp *rocketpool.RocketPool) (*MinipoolQueue, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolQueue)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool queue contract: %w", err)
	}

	return &MinipoolQueue{
		Details:  MinipoolQueueDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the total length of the minipool queue
func (c *MinipoolQueue) GetTotalLength(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TotalLength.RawValue, "getTotalLength")
}

// Get the total capacity of the minipool queue
func (c *MinipoolQueue) GetTotalCapacity(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.TotalCapacity, "getTotalCapacity")
}

// Get the total effective capacity of the minipool queue (used in node demand calculation)
func (c *MinipoolQueue) GetEffectiveCapacity(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.EffectiveCapacity, "getEffectiveCapacity")
}

// =============
// === Utils ===
// =============

// Get the minipool at the specified position in queue (0-indexed).
func (c *MinipoolQueue) GetQueueMinipoolAtPosition(mc *multicall.MultiCaller, address_Out *common.Address, position uint64) {
	multicall.AddCall(mc, c.contract, address_Out, "getMinipoolAt", big.NewInt(int64(position)))
}
