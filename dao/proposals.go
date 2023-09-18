package dao

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// Settings
const (
	proposalBatchSize int = 100
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
// NOTE: Proposals are 1-indexed
func (c *DaoProposal) GetProposalCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.ProposalCount.RawValue, "getTotal")
}

// =============
// === Utils ===
// =============

// Get all of the Protocol DAO proposals
// Returns: Protocol DAO proposals, Oracle DAO proposals, error
// NOTE: Proposals are 1-indexed
func (c *DaoProposal) GetProposals(rp *rocketpool.RocketPool, proposalCount uint64, includeDetails bool, opts *bind.CallOpts) ([]*Proposal, []*Proposal, error) {
	props := make([]*Proposal, proposalCount)
	for i := uint64(1); i <= proposalCount; i++ { // Proposals are 1-indexed
		prop, err := NewProposal(rp, i)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating DAO proposal %d", i)
		}
		props[i-1] = prop
	}

	err := rp.BatchQuery(int(proposalCount), proposalBatchSize, func(mc *batch.MultiCaller, index int) error {
		if includeDetails {
			props[index].GetAllDetails(mc)
		} else {
			props[index].GetDAO(mc)   // Needed for sorting
			props[index].GetState(mc) // Frequently needed for determining which proposals to query more
		}
		return nil
	}, opts)
	if err != nil {
		return nil, nil, err
	}

	pDaoProps := []*Proposal{}
	oDaoProps := []*Proposal{}
	for _, prop := range props {
		if prop.Details.DAO == string(rocketpool.ContractName_RocketDAOProtocolProposals) {
			pDaoProps = append(pDaoProps, prop)
		} else if prop.Details.DAO == string(rocketpool.ContractName_RocketDAONodeTrustedProposals) {
			// oDAO
			oDaoProps = append(oDaoProps, prop)
		} else {
			return nil, nil, fmt.Errorf("proposal %d has DAO [%s] which is not recognized", prop.Details.ID.Formatted(), prop.Details.DAO)
		}
	}

	return pDaoProps, oDaoProps, nil
}
