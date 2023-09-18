package proposals

import (
	"fmt"

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
	mgr     *core.Contract
}

// Details for proposals
type OracleDaoProposalDetails struct {
}

// ====================
// === Constructors ===
// ====================

// Creates a new OracleDaoProposal contract binding
func newOracleDaoProposal(rp *rocketpool.RocketPool, base *ProposalCommon) (*OracleDaoProposal, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAOProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO proposal contract: %w", err)
	}

	return &OracleDaoProposal{
		ProposalCommon: base,
		Details:        OracleDaoProposalDetails{},
		rp:             rp,
		mgr:            contract,
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
