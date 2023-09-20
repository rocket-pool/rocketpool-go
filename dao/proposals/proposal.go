package proposals

import (
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
)

// ==================
// === Interfaces ===
// ==================

type IProposal interface {
	// Get all of the proposal's details
	QueryAllDetails(mc *batch.MultiCaller)

	// Get all of the details common across each type of proposal
	GetCommonDetails() *ProposalCommonDetails

	// Get the address of the node that created the proposal
	GetProposerAddress(mc *batch.MultiCaller)

	// Get the message provided with the proposal
	GetMessage(mc *batch.MultiCaller)

	// Get the time the proposal was created
	GetCreatedTime(mc *batch.MultiCaller)

	// Get the time the voting window on the proposal started
	GetStartTime(mc *batch.MultiCaller)

	// Get the time the voting window on the proposal ended
	GetEndTime(mc *batch.MultiCaller)

	// Get the time the proposal expires
	GetExpiryTime(mc *batch.MultiCaller)

	// Get the number of votes required for the proposal to pass
	GetVotesRequired(mc *batch.MultiCaller)

	// Get the number of votes in favor of the proposal
	GetVotesFor(mc *batch.MultiCaller)

	// Get the number of votes against the proposal
	GetVotesAgainst(mc *batch.MultiCaller)

	// Check if the proposal has been cancelled
	GetIsCancelled(mc *batch.MultiCaller)

	// Check if the proposal has been executed
	GetIsExecuted(mc *batch.MultiCaller)

	// Get the proposal's payload
	GetPayload(mc *batch.MultiCaller)

	// Get the proposal's state
	GetState(mc *batch.MultiCaller)

	// Check if a node has voted on the proposal
	GetMemberHasVoted(mc *batch.MultiCaller, out *bool, address common.Address)

	// Check if a node has voted in favor of the proposal
	GetMemberSupported(mc *batch.MultiCaller, out *bool, address common.Address)

	// Get the proposal's payload as a string
	GetPayloadAsString() (string, error)

	// Get which DAO the proposal is for - reserved for internal use
	getDAO(mc *batch.MultiCaller, dao_Out *string)
}
