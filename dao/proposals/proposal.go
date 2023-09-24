package proposals

import (
	batch "github.com/rocket-pool/batch-query"
)

// ==================
// === Interfaces ===
// ==================

type IProposal interface {
	// Get all of the proposal's details
	QueryAllFields(mc *batch.MultiCaller)

	// Get the binding common across each type of proposal
	Common() *ProposalCommon
}
