package proposals

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/rocketpool"

	strutils "github.com/rocket-pool/rocketpool-go/utils/strings"
)

// ==================
// === Interfaces ===
// ==================

type IProposal interface {
	QueryAllDetails(mc *batch.MultiCaller)
	GetProposalCommon() *ProposalCommon
}

// ====================
// === Constructors ===
// ====================

// Create a minipool binding from an explicit version number
func NewProposalFromDao(rp *rocketpool.RocketPool, id uint64, dao rocketpool.ContractName) (IProposal, error) {
	base, err := newProposalCommon(rp, id)
	if err != nil {
		return nil, fmt.Errorf("error creating common proposal binding: %w", err)
	}

	switch dao {
	case rocketpool.ContractName_RocketDAOProtocolProposals:
		return newProtocolDaoProposal(rp, base)
	case rocketpool.ContractName_RocketDAONodeTrustedProposals:
		return newOracleDaoProposal(rp, base)
	default:
		return nil, fmt.Errorf("unexpected proposal DAO [%s]", dao)
	}
}

// =============
// === Utils ===
// =============

// Get the proposal's payload as a string
func GetPayloadAsString(rp *rocketpool.RocketPool, daoName string, payload []byte) (string, error) {
	// Get the ABI
	contract, err := rp.GetContract(rocketpool.ContractName(daoName))
	if err != nil {
		return "", fmt.Errorf("error getting contract [%s]: %w", daoName, err)
	}
	daoContractAbi := contract.ABI

	// Get proposal payload method
	method, err := daoContractAbi.MethodById(payload)
	if err != nil {
		return "", fmt.Errorf("Could not get proposal payload method: %w", err)
	}

	// Get proposal payload argument values
	args, err := method.Inputs.UnpackValues(payload[4:])
	if err != nil {
		return "", fmt.Errorf("Could not get proposal payload arguments: %w", err)
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
