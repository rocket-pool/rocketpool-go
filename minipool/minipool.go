package minipool

import (
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
)

// ==================
// === Interfaces ===
// ==================

type IMinipool interface {
	// Get all of the minipool's details
	QueryAllFields(mc *batch.MultiCaller)

	// Gets the underlying minipool's contract
	GetContract() *core.Contract

	// Gets the common binding for all minipool types
	Common() *MinipoolCommon
}
