package dao

import (
	"fmt"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// Settings
const (
	ProposalDAONamesBatchSize = 50
	ProposalDetailsBatchSize  = 10
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProposal
type DaoProposal struct {
	Details  DaoProposalDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for RocketDAOProposal
type DaoProposalDetails struct {
	ProposalCount core.Parameter[uint64] `json:"proposalCount"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProposal contract binding
func NewDaoProposal(rp *rocketpool.RocketPool) (*DaoProposal, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAOProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO proposal contract: %w", err)
	}

	return &DaoProposal{
		Details:  DaoProposalDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the total number of DAO proposals
func (c *DaoProposal) GetTotalRPLBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.ProposalCount.RawValue, "getTotal")
}
