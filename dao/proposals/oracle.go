package proposals

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for Oracle DAO proposals
type OracleDaoProposal struct {
	*ProposalCommon

	// === Internal fields ===
	rp   *rocketpool.RocketPool
	dntp *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new OracleDaoProposal contract binding
func newOracleDaoProposal(rp *rocketpool.RocketPool, base *ProposalCommon) (*OracleDaoProposal, error) {
	// Create the contract
	dntp, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrustedProposals)
	if err != nil {
		return nil, fmt.Errorf("error getting Oracle DAO proposals contract: %w", err)
	}

	return &OracleDaoProposal{
		ProposalCommon: base,
		rp:             rp,
		dntp:           dntp,
	}, nil
}

// Get a proposal as an Oracle DAO propsal
func GetProposalAsOracle(proposal IProposal) (*OracleDaoProposal, bool) {
	castedProp, ok := proposal.(*OracleDaoProposal)
	if ok {
		return castedProp, true
	}
	return nil, false
}

// =============
// === Calls ===
// =============

// Get the basic details
func (c *OracleDaoProposal) QueryAllDetails(mc *batch.MultiCaller) {
	c.ProposalCommon.QueryAllDetails(mc)
}

// Get the common fields
func (c *OracleDaoProposal) GetProposalCommon() *ProposalCommon {
	return c.ProposalCommon
}

// ====================
// === Transactions ===
// ====================

// Get info for cancelling a proposal
func (c *OracleDaoProposal) Cancel(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dntp, "cancel", opts, c.idBig)
}

// Get info for voting on a proposal
func (c *OracleDaoProposal) VoteOn(support bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dntp, "vote", opts, c.idBig, support)
}

// Get info for executing a proposal
func (c *OracleDaoProposal) Execute(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dntp, "execute", opts, c.idBig)
}

// Get the proposal's payload as a string
func (c *OracleDaoProposal) GetPayloadAsString() (string, error) {
	return getPayloadAsStringImpl(c.rp, c.dntp, c.Payload.Get())
}
