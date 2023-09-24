package proposals

import (
	batch "github.com/rocket-pool/batch-query"
)

// ==================
// === Interfaces ===
// ==================

type IProposal interface {
	// Get all of the proposal's details
	QueryAllDetails(mc *batch.MultiCaller)

	// Get all of the details common across each type of proposal
	GetProposalCommon() *ProposalCommon
}
