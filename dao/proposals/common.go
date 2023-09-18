package proposals

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
)

// ===============
// === Structs ===
// ===============

// Binding for proposals
type ProposalCommon struct {
	Details  ProposalCommonDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for proposals
type ProposalCommonDetails struct {
	ID              core.Parameter[uint64]                   `json:"id"`
	ProposerAddress common.Address                           `json:"proposerAddress"`
	Message         string                                   `json:"message"`
	CreatedTime     core.Parameter[time.Time]                `json:"createdTime"`
	StartTime       core.Parameter[time.Time]                `json:"startTime"`
	EndTime         core.Parameter[time.Time]                `json:"endTime"`
	ExpiryTime      core.Parameter[time.Time]                `json:"expiryTime"`
	VotesRequired   core.Parameter[float64]                  `json:"votesRequired"`
	VotesFor        core.Parameter[float64]                  `json:"votesFor"`
	VotesAgainst    core.Parameter[float64]                  `json:"votesAgainst"`
	MemberVoted     bool                                     `json:"memberVoted"`
	MemberSupported bool                                     `json:"memberSupported"`
	IsCancelled     bool                                     `json:"isCancelled"`
	IsExecuted      bool                                     `json:"isExecuted"`
	Payload         []byte                                   `json:"payload"`
	State           core.Uint8Parameter[types.ProposalState] `json:"state"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProposalCommon contract binding
func newProposalCommon(rp *rocketpool.RocketPool, id uint64) (*ProposalCommon, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAOProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO proposal contract: %w", err)
	}

	return &ProposalCommon{
		Details: ProposalCommonDetails{
			ID: core.Parameter[uint64]{
				RawValue: big.NewInt(0).SetUint64(id),
			},
		},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the address of the node that created the proposal
func (c *ProposalCommon) GetProposerAddress(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.ProposerAddress, "getProposer", c.Details.ID.RawValue)
}

// Get the message provided with the proposal
func (c *ProposalCommon) GetMessage(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.Message, "getMessage", c.Details.ID.RawValue)
}

// Get the time the proposal was created
func (c *ProposalCommon) GetCreatedTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.CreatedTime.RawValue, "getCreated", c.Details.ID.RawValue)
}

// Get the time the voting window on the proposal started
func (c *ProposalCommon) GetStartTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.StartTime.RawValue, "getStart", c.Details.ID.RawValue)
}

// Get the time the voting window on the proposal ended
func (c *ProposalCommon) GetEndTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.EndTime.RawValue, "getEnd", c.Details.ID.RawValue)
}

// Get the time the proposal expires
func (c *ProposalCommon) GetExpiryTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.ExpiryTime.RawValue, "getExpires", c.Details.ID.RawValue)
}

// Get the number of votes required for the proposal to pass
func (c *ProposalCommon) GetVotesRequired(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.VotesRequired.RawValue, "getVotesRequired", c.Details.ID.RawValue)
}

// Get the number of votes in favor of the proposal
func (c *ProposalCommon) GetVotesFor(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.VotesFor.RawValue, "getVotesFor", c.Details.ID.RawValue)
}

// Get the number of votes against the proposal
func (c *ProposalCommon) GetVotesAgainst(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.VotesAgainst.RawValue, "getVotesAgainst", c.Details.ID.RawValue)
}

// Check if the proposal has been cancelled
func (c *ProposalCommon) GetIsCancelled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.IsCancelled, "getCancelled", c.Details.ID.RawValue)
}

// Check if the proposal has been executed
func (c *ProposalCommon) GetIsExecuted(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.IsExecuted, "getExecuted", c.Details.ID.RawValue)
}

// Get the proposal's payload
func (c *ProposalCommon) GetPayload(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.Payload, "getPayload", c.Details.ID.RawValue)
}

// Get the proposal's state
func (c *ProposalCommon) GetState(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.State.RawValue, "getState", c.Details.ID.RawValue)
}

// Get all of the proposal's details
func (c *ProposalCommon) QueryAllDetails(mc *batch.MultiCaller) {
	c.GetProposerAddress(mc)
	c.GetMessage(mc)
	c.GetCreatedTime(mc)
	c.GetStartTime(mc)
	c.GetEndTime(mc)
	c.GetExpiryTime(mc)
	c.GetVotesRequired(mc)
	c.GetVotesFor(mc)
	c.GetVotesAgainst(mc)
	c.GetIsCancelled(mc)
	c.GetIsExecuted(mc)
	c.GetPayload(mc)
	c.GetState(mc)
}

// Check if a node has voted on the proposal
func (c *ProposalCommon) GetMemberHasVoted(mc *batch.MultiCaller, out *bool, address common.Address) {
	core.AddCall(mc, c.contract, out, "getReceiptHasVoted", c.Details.ID.RawValue, address)
}

// Check if a node has voted in favor of the proposal
func (c *ProposalCommon) GetMemberSupported(mc *batch.MultiCaller, out *bool, address common.Address) {
	core.AddCall(mc, c.contract, out, "getReceiptSupported", c.Details.ID.RawValue, address)
}

// Get which DAO the proposal is for - reserved for internal use
func (c *ProposalCommon) getDAO(mc *batch.MultiCaller, dao_Out *string) {
	core.AddCall(mc, c.contract, dao_Out, "getDAO", c.Details.ID.RawValue)
}
