package proposals

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/v2/core"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
)

// Settings
const (
	proposalBatchSize int = 100
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAOProposal
type DaoProposalManager struct {
	// The total number of DAO proposals
	// NOTE: Proposals are 1-indexed
	ProposalCount *core.FormattedUint256Field[uint64]

	// === Internal fields ===
	rp *rocketpool.RocketPool
	dp *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProposalManager contract binding
func NewDaoProposalManager(rp *rocketpool.RocketPool) (*DaoProposalManager, error) {
	// Create the contract
	dp, err := rp.GetContract(rocketpool.ContractName_RocketDAOProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO proposal manager contract: %w", err)
	}

	return &DaoProposalManager{
		ProposalCount: core.NewFormattedUint256Field[uint64](dp, "getTotal"),

		rp: rp,
		dp: dp,
	}, nil
}

// =============
// === Utils ===
// =============

// Create a proposal binding from an explicit DAO ID if you already know what it is
func (c *DaoProposalManager) NewProposalFromDao(id uint64, dao rocketpool.ContractName) (IProposal, error) {
	base, err := newProposalCommon(c.rp, id)
	if err != nil {
		return nil, fmt.Errorf("error creating common proposal binding: %w", err)
	}

	switch dao {
	case rocketpool.ContractName_RocketDAONodeTrustedProposals:
		return newOracleDaoProposal(c.rp, base)
	case rocketpool.ContractName_RocketDAOSecurityProposals:
		return newSecurityCouncilProposal(c.rp, base)
	case "":
		return nil, fmt.Errorf("proposal %d does not exist", id)
	default:
		return nil, fmt.Errorf("unexpected proposal DAO [%s]", dao)
	}
}

// Create a proposal binding by ID
func (c *DaoProposalManager) CreateProposalFromID(id uint64, opts *bind.CallOpts) (IProposal, error) {
	prop, err := newProposalCommon(c.rp, id)
	if err != nil {
		return nil, fmt.Errorf("error creating DAO proposal: %w", err)
	}

	var dao string
	err = c.rp.Query(func(mc *batch.MultiCaller) error {
		prop.getDAO(mc, &dao)
		return nil
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting proposal DAO: %w", err)
	}

	switch dao {
	case string(rocketpool.ContractName_RocketDAONodeTrustedProposals):
		return newOracleDaoProposal(c.rp, prop)
	case string(rocketpool.ContractName_RocketDAOSecurityProposals):
		return newSecurityCouncilProposal(c.rp, prop)
	case "":
		return nil, fmt.Errorf("proposal %d does not exist", id)
	default:
		return nil, fmt.Errorf("unexpected proposal DAO [%s]", dao)
	}
}

// Get all of the Oracle DAO and security council proposals
// NOTE: Proposals are 1-indexed
func (c *DaoProposalManager) GetProposals(proposalCount uint64, includeDetails bool, opts *bind.CallOpts) ([]*OracleDaoProposal, []*SecurityCouncilProposal, error) {
	// Create prop commons for each one
	props := make([]*ProposalCommon, proposalCount)
	for i := uint64(1); i <= proposalCount; i++ { // Proposals are 1-indexed
		prop, err := newProposalCommon(c.rp, i)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating DAO proposal %d: %w", i, err)
		}
		props[i-1] = prop
	}

	// Get the DAOs
	daos := make([]string, len(props))
	err := c.rp.BatchQuery(len(props), proposalBatchSize, func(mc *batch.MultiCaller, i int) error {
		props[i].getDAO(mc, &daos[i])
		return nil
	}, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting proposal DAOs: %w", err)
	}

	// Construct concrete bindings for each one
	oDaoProps := []*OracleDaoProposal{}
	secProps := []*SecurityCouncilProposal{}
	totalProps := []IProposal{}
	for i, prop := range props {
		switch daos[i] {
		case string(rocketpool.ContractName_RocketDAONodeTrustedProposals):
			odaoProp, err := newOracleDaoProposal(c.rp, prop)
			if err != nil {
				return nil, nil, fmt.Errorf("error creating Oracle DAO proposal binding for proposal %d: %w", prop.ID, err)
			}
			oDaoProps = append(oDaoProps, odaoProp)
			totalProps = append(totalProps, odaoProp)

		case string(rocketpool.ContractName_RocketDAOSecurityProposals):
			secProp, err := newSecurityCouncilProposal(c.rp, prop)
			if err != nil {
				return nil, nil, fmt.Errorf("error creating security council proposal binding for proposal %d: %w", prop.ID, err)
			}
			secProps = append(secProps, secProp)
			totalProps = append(totalProps, secProp)

		default:
			return nil, nil, fmt.Errorf("proposal %d has DAO [%s] which is not recognized", prop.ID, daos[i])
		}
	}

	// Get all details if requested
	if includeDetails {
		err = c.rp.BatchQuery(int(proposalCount), proposalBatchSize, func(mc *batch.MultiCaller, index int) error {
			totalProps[index].QueryAllFields(mc)
			return nil
		}, opts)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting proposal details: %w", err)
		}
	}

	return oDaoProps, secProps, nil
}
