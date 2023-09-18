package proposals

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	strutils "github.com/rocket-pool/rocketpool-go/utils/strings"
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
	*DaoProposalManagerDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for DaoProposalManager
type DaoProposalManagerDetails struct {
	ProposalCount core.Parameter[uint64] `json:"proposalCount"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProposalManager contract binding
func NewDaoProposalManager(rp *rocketpool.RocketPool) (*DaoProposalManager, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAOProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO proposal manager contract: %w", err)
	}

	return &DaoProposalManager{
		DaoProposalManagerDetails: &DaoProposalManagerDetails{},
		rp:                        rp,
		contract:                  contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the total number of DAO proposals
// NOTE: Proposals are 1-indexed
func (c *DaoProposalManager) GetProposalCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.ProposalCount.RawValue, "getTotal")
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
	case rocketpool.ContractName_RocketDAOProtocolProposals:
		return newProtocolDaoProposal(c.rp, base)
	case rocketpool.ContractName_RocketDAONodeTrustedProposals:
		return newOracleDaoProposal(c.rp, base)
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
	case string(rocketpool.ContractName_RocketDAOProtocolProposals):
		return newProtocolDaoProposal(c.rp, prop)
	case string(rocketpool.ContractName_RocketDAONodeTrustedProposals):
		return newOracleDaoProposal(c.rp, prop)
	default:
		return nil, fmt.Errorf("unexpected proposal DAO [%s]", dao)
	}
}

// Get all of the Protocol DAO proposals
// NOTE: Proposals are 1-indexed
func (c *DaoProposalManager) GetProposals(proposalCount uint64, includeDetails bool, opts *bind.CallOpts) ([]*ProtocolDaoProposal, []*OracleDaoProposal, error) {
	// Create prop commons for each one
	props := make([]*proposalCommon, proposalCount)
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
	pDaoProps := []*ProtocolDaoProposal{}
	oDaoProps := []*OracleDaoProposal{}
	totalProps := []IProposal{}
	for i, prop := range props {
		switch daos[i] {
		case string(rocketpool.ContractName_RocketDAOProtocolProposals):
			pdaoProp, err := newProtocolDaoProposal(c.rp, prop)
			if err != nil {
				return nil, nil, fmt.Errorf("error creating Oracle DAO proposal binding for proposal %d: %w", prop.ID.Formatted(), err)
			}
			pDaoProps = append(pDaoProps, pdaoProp)
			totalProps = append(totalProps, pdaoProp)

		case string(rocketpool.ContractName_RocketDAONodeTrustedProposals):
			odaoProp, err := newOracleDaoProposal(c.rp, prop)
			if err != nil {
				return nil, nil, fmt.Errorf("error creating Oracle DAO proposal binding for proposal %d: %w", prop.ID.Formatted(), err)
			}
			oDaoProps = append(oDaoProps, odaoProp)
			totalProps = append(totalProps, odaoProp)

		default:
			return nil, nil, fmt.Errorf("proposal %d has DAO [%s] which is not recognized", prop.ID.Formatted(), daos[i])
		}
	}

	// Get all details if requested
	if includeDetails {
		err = c.rp.BatchQuery(int(proposalCount), proposalBatchSize, func(mc *batch.MultiCaller, index int) error {
			totalProps[index].QueryAllDetails(mc)
			return nil
		}, opts)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting ")
		}
	}

	return pDaoProps, oDaoProps, nil
}

// Get the proposal's payload as a string
func (c *DaoProposalManager) GetPayloadAsString(daoName rocketpool.ContractName, payload []byte) (string, error) {
	// Get the ABI
	contract, err := c.rp.GetContract(daoName)
	if err != nil {
		return "", fmt.Errorf("error getting contract [%s]: %w", daoName, err)
	}
	daoContractAbi := contract.ABI

	// Get proposal payload method
	method, err := daoContractAbi.MethodById(payload)
	if err != nil {
		return "", fmt.Errorf("error getting proposal payload method: %w", err)
	}

	// Get proposal payload argument values
	args, err := method.Inputs.UnpackValues(payload[4:])
	if err != nil {
		return "", fmt.Errorf("error getting proposal payload arguments: %w", err)
	}

	// Format argument values as strings
	argStrs := []string{}
	for ai, arg := range args {
		switch method.Inputs[ai].Type.T {
		case abi.AddressTy:
			argStrs = append(argStrs, arg.(common.Address).Hex())
		case abi.HashTy:
			argStrs = append(argStrs, arg.(common.Hash).Hex())
		case abi.FixedBytesTy:
			fallthrough
		case abi.BytesTy:
			argStrs = append(argStrs, hex.EncodeToString(arg.([]byte)))
		default:
			argStrs = append(argStrs, fmt.Sprintf("%v", arg))
		}
	}

	// Build & return payload string
	return strutils.Sanitize(fmt.Sprintf("%s(%s)", method.RawName, strings.Join(argStrs, ","))), nil
}
