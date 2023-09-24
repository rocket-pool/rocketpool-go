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
	// The proposal's ID
	ID uint64

	// The address of the node that created the proposal
	ProposerAddress *core.SimpleField[common.Address]

	// The message provided with the proposal
	Message *core.SimpleField[string]

	// The time the proposal was created
	CreatedTime *core.FormattedUint256Field[time.Time]

	// The time the voting window on the proposal started
	StartTime *core.FormattedUint256Field[time.Time]

	// The time the voting window on the proposal ended
	EndTime *core.FormattedUint256Field[time.Time]

	// The time the proposal expires
	ExpiryTime *core.FormattedUint256Field[time.Time]

	// The number of votes required for the proposal to pass
	VotesRequired *core.FormattedUint256Field[float64]

	// The number of votes in favor of the proposal
	VotesFor *core.FormattedUint256Field[float64]

	// The number of votes against the proposal
	VotesAgainst *core.FormattedUint256Field[float64]

	// True if the proposal has been cancelled
	IsCancelled *core.SimpleField[bool]

	// True if the proposal has been executed
	IsExecuted *core.SimpleField[bool]

	// The proposal's payload
	Payload *core.SimpleField[[]byte]

	// The proposal's state
	State *core.FormattedUint8Field[types.ProposalState]

	// === Internal fields ===
	idBig *big.Int
	rp    *rocketpool.RocketPool
	dp    *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProposalCommon contract binding
func newProposalCommon(rp *rocketpool.RocketPool, id uint64) (*ProposalCommon, error) {
	// Create the contract
	dp, err := rp.GetContract(rocketpool.ContractName_RocketDAOProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO proposal contract: %w", err)
	}

	idBig := big.NewInt(0).SetUint64(id)
	return &ProposalCommon{
		ID:              id,
		ProposerAddress: core.NewSimpleField[common.Address](dp, "getProposer", idBig),
		Message:         core.NewSimpleField[string](dp, "getMessage", idBig),
		CreatedTime:     core.NewFormattedUint256Field[time.Time](dp, "getCreated", idBig),
		StartTime:       core.NewFormattedUint256Field[time.Time](dp, "getStart", idBig),
		EndTime:         core.NewFormattedUint256Field[time.Time](dp, "getEnd", idBig),
		ExpiryTime:      core.NewFormattedUint256Field[time.Time](dp, "getExpires", idBig),
		VotesRequired:   core.NewFormattedUint256Field[float64](dp, "getVotesRequired", idBig),
		VotesFor:        core.NewFormattedUint256Field[float64](dp, "getVotesFor", idBig),
		VotesAgainst:    core.NewFormattedUint256Field[float64](dp, "getVotesAgainst", idBig),
		IsCancelled:     core.NewSimpleField[bool](dp, "getCancelled", idBig),
		IsExecuted:      core.NewSimpleField[bool](dp, "getExecuted", idBig),
		Payload:         core.NewSimpleField[[]byte](dp, "getPayload", idBig),
		State:           core.NewFormattedUint8Field[types.ProposalState](dp, "getState", idBig),

		idBig: idBig,
		rp:    rp,
		dp:    dp,
	}, nil
}

// =============
// === Calls ===
// =============

// Get all of the proposal's details
func (c *ProposalCommon) QueryAllDetails(mc *batch.MultiCaller) {
	core.AddQueryablesToMulticall(mc,
		c.ProposerAddress,
		c.Message,
		c.CreatedTime,
		c.StartTime,
		c.EndTime,
		c.ExpiryTime,
		c.VotesRequired,
		c.VotesFor,
		c.VotesAgainst,
		c.IsCancelled,
		c.IsExecuted,
		c.Payload,
		c.State,
	)
}

// Check if a node has voted on the proposal
func (c *ProposalCommon) GetMemberHasVoted(mc *batch.MultiCaller, out *bool, address common.Address) {
	core.AddCall(mc, c.dp, out, "getReceiptHasVoted", c.idBig, address)
}

// Check if a node has voted in favor of the proposal
func (c *ProposalCommon) GetMemberSupported(mc *batch.MultiCaller, out *bool, address common.Address) {
	core.AddCall(mc, c.dp, out, "getReceiptSupported", c.idBig, address)
}

// Get which DAO the proposal is for - reserved for internal use
func (c *ProposalCommon) getDAO(mc *batch.MultiCaller, dao_Out *string) {
	core.AddCall(mc, c.dp, dao_Out, "getDAO", c.idBig)
}
