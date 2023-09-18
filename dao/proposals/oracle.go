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
	Details OracleDaoProposalDetails
	rp      *rocketpool.RocketPool
	dntp    *core.Contract
}

// Details for proposals
type OracleDaoProposalDetails struct {
	*ProposalCommonDetails
}

// ====================
// === Constructors ===
// ====================

// Creates a new OracleDaoProposal contract binding
func newOracleDaoProposal(rp *rocketpool.RocketPool, base *ProposalCommon) (*OracleDaoProposal, error) {
	// Create the dntp
	dntp, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrustedProposals)
	if err != nil {
		return nil, fmt.Errorf("error getting Oracle DAO proposals contract: %w", err)
	}

	return &OracleDaoProposal{
		ProposalCommon: base,
		Details: OracleDaoProposalDetails{
			ProposalCommonDetails: &base.Details,
		},
		rp:   rp,
		dntp: dntp,
	}, nil
}

// =============
// === Calls ===
// =============

func (c *OracleDaoProposal) GetProposalCommon() *ProposalCommon {
	return c.ProposalCommon
}

// Get the basic details
func (c *OracleDaoProposal) QueryAllDetails(mc *batch.MultiCaller) {
	c.ProposalCommon.QueryAllDetails(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for cancelling a proposal
func (c *OracleDaoProposal) Cancel(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dntp, "cancel", opts, c.Details.ID.RawValue)
}

// Get info for voting on a proposal
func (c *OracleDaoProposal) VoteOn(support bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dntp, "vote", opts, c.Details.ID.RawValue, support)
}

// Get info for executing a proposal
func (c *OracleDaoProposal) Execute(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.dntp, "execute", opts, c.Details.ID.RawValue)
}
