package minipool

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
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
	*MinipoolQueueDetails
	rp *rocketpool.RocketPool
	mq *core.Contract
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
	mq, err := rp.GetContract(rocketpool.ContractName_RocketMinipoolQueue)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool queue contract: %w", err)
	}

	return &MinipoolQueue{
		MinipoolQueueDetails: &MinipoolQueueDetails{},
		rp:                   rp,
		mq:                   mq,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the total length of the minipool queue
func (c *MinipoolQueue) GetTotalLength(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mq, &c.TotalLength.RawValue, "getTotalLength")
}

// Get the total capacity of the minipool queue
func (c *MinipoolQueue) GetTotalCapacity(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mq, &c.TotalCapacity, "getTotalCapacity")
}

// Get the total effective capacity of the minipool queue (used in node demand calculation)
func (c *MinipoolQueue) GetEffectiveCapacity(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mq, &c.EffectiveCapacity, "getEffectiveCapacity")
}

// =============
// === Utils ===
// =============

// Get the minipool at the specified position in queue (0-indexed).
func (c *MinipoolQueue) GetQueueMinipoolAtPosition(mc *batch.MultiCaller, address_Out *common.Address, position uint64) {
	core.AddCall(mc, c.mq, address_Out, "getMinipoolAt", big.NewInt(int64(position)))
}
