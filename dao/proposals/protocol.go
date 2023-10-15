package proposals

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for Protocol DAO proposals
type ProtocolDaoProposal struct {
	*ProposalCommon

	// The block that was used for voting power calculation in a proposal
	ProposalBlock *core.FormattedUint256Field[uint32]

	// The veto quorum required to veto a proposal
	VetoQuorum *core.SimpleField[*big.Int]

	// === Internal fields ===
	rp  *rocketpool.RocketPool
	dpp *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProtocolDaoProposal contract binding
func newProtocolDaoProposal(rp *rocketpool.RocketPool, base *ProposalCommon) (*ProtocolDaoProposal, error) {
	// Create the contract
	dpp, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocolProposals)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO proposal contract: %w", err)
	}

	return &ProtocolDaoProposal{
		ProposalCommon: base,
		ProposalBlock:  core.NewFormattedUint256Field[uint32](dpp, "getProposalBlock", base.idBig),
		VetoQuorum:     core.NewSimpleField[*big.Int](dpp, "getProposalVetoQuorum", base.idBig),

		rp:  rp,
		dpp: dpp,
	}, nil
}

// Get a proposal as an Protocol DAO propsal
func GetProposalAsProtocol(proposal IProposal) (*ProtocolDaoProposal, bool) {
	castedProp, ok := proposal.(*ProtocolDaoProposal)
	if ok {
		return castedProp, true
	}
	return nil, false
}

// =============
// === Calls ===
// =============

// Get the basic details
func (c *ProtocolDaoProposal) QueryAllFields(mc *batch.MultiCaller) {
	core.QueryAllFields(c.ProposalCommon, mc)
}

// Get the common fields
func (c *ProtocolDaoProposal) Common() *ProposalCommon {
	return c.ProposalCommon
}

// Get the proposal's payload as a string
func (c *ProtocolDaoProposal) GetPayloadAsString() (string, error) {
	return getPayloadAsStringImpl(c.rp, c.dpp, c.Payload.Get())
}

// ====================
// === Transactions ===
// ====================

// Get info for voting on a proposal
func (c *ProtocolDaoProposal) VoteOn(support bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dpp, "vote", opts, c.idBig, support)
}

// Get info for executing a proposal
func (c *ProtocolDaoProposal) Execute(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dpp, "execute", opts, c.idBig)
}
