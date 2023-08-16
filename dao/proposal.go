package dao

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"

	strutils "github.com/rocket-pool/rocketpool-go/utils/strings"
)

// ===============
// === Structs ===
// ===============

// Binding for proposals
type Proposal struct {
	Details ProposalDetails
	rp      *rocketpool.RocketPool
	mgr     *core.Contract
}

// Details for proposals
type ProposalDetails struct {
	ID              core.Parameter[uint64]                   `json:"id"`
	DAO             string                                   `json:"dao"`
	ProposerAddress common.Address                           `json:"proposerAddress"`
	Message         string                                   `json:"message"`
	CreatedTime     core.Parameter[time.Time]                `json:"createdTime"`
	StartTime       core.Parameter[time.Time]                `json:"startTime"`
	EndTime         core.Parameter[time.Time]                `json:"endTime"`
	ExpiryTime      core.Parameter[time.Time]                `json:"expiryTime"`
	VotesRequired   core.Parameter[float64]                  `json:"votesRequired"`
	VotesFor        core.Parameter[float64]                  `json:"votesFor"`
	VotesAgainst    core.Parameter[float64]                  `json:"votesAgainst"`
	MemberVoted     bool                                     `json:"memberVoted"`
	MemberSupported bool                                     `json:"memberSupported"`
	IsCancelled     bool                                     `json:"isCancelled"`
	IsExecuted      bool                                     `json:"isExecuted"`
	Payload         []byte                                   `json:"payload"`
	State           core.Uint8Parameter[types.ProposalState] `json:"state"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProposal contract binding
func NewProposal(rp *rocketpool.RocketPool, id uint64) (*Proposal, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAOProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO proposal contract: %w", err)
	}

	return &Proposal{
		Details: ProposalDetails{
			ID: core.Parameter[uint64]{
				RawValue: big.NewInt(0).SetUint64(id),
			},
		},
		rp:  rp,
		mgr: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get which DAO the proposal is for
func (c *Proposal) GetDAO(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.DAO, "getDAO", c.Details.ID.RawValue)
}

// Get the address of the node that created the proposal
func (c *Proposal) GetProposerAddress(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.ProposerAddress, "getProposer", c.Details.ID.RawValue)
}

// Get the message provided with the proposal
func (c *Proposal) GetMessage(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.Message, "getMessage", c.Details.ID.RawValue)
}

// Get the time the proposal was created
func (c *Proposal) GetCreatedTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.CreatedTime.RawValue, "getCreated", c.Details.ID.RawValue)
}

// Get the time the voting window on the proposal started
func (c *Proposal) GetStartTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.StartTime.RawValue, "getStart", c.Details.ID.RawValue)
}

// Get the time the voting window on the proposal ended
func (c *Proposal) GetEndTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.EndTime.RawValue, "getEnd", c.Details.ID.RawValue)
}

// Get the time the proposal expires
func (c *Proposal) GetExpiryTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.ExpiryTime.RawValue, "getExpires", c.Details.ID.RawValue)
}

// Get the number of votes required for the proposal to pass
func (c *Proposal) GetVotesRequired(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.VotesRequired.RawValue, "getVotesRequired", c.Details.ID.RawValue)
}

// Get the number of votes in favor of the proposal
func (c *Proposal) GetVotesFor(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.VotesFor.RawValue, "getVotesFor", c.Details.ID.RawValue)
}

// Get the number of votes against the proposal
func (c *Proposal) GetVotesAgainst(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.VotesAgainst.RawValue, "getVotesAgainst", c.Details.ID.RawValue)
}

// Check if the proposal has been cancelled
func (c *Proposal) GetIsCancelled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.IsCancelled, "getCancelled", c.Details.ID.RawValue)
}

// Check if the proposal has been executed
func (c *Proposal) GetIsExecuted(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.IsExecuted, "getExecuted", c.Details.ID.RawValue)
}

// Get the proposal's payload
func (c *Proposal) GetPayload(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.Payload, "getPayload", c.Details.ID.RawValue)
}

// Get the proposal's state
func (c *Proposal) GetState(mc *batch.MultiCaller) {
	core.AddCall(mc, c.mgr, &c.Details.State.RawValue, "getState", c.Details.ID.RawValue)
}

// Get all of the proposal's details
func (c *Proposal) GetAllDetails(mc *batch.MultiCaller) {
	c.GetDAO(mc)
	c.GetProposerAddress(mc)
	c.GetMessage(mc)
	c.GetCreatedTime(mc)
	c.GetStartTime(mc)
	c.GetEndTime(mc)
	c.GetExpiryTime(mc)
	c.GetVotesRequired(mc)
	c.GetVotesFor(mc)
	c.GetVotesAgainst(mc)
	c.GetIsCancelled(mc)
	c.GetIsExecuted(mc)
	c.GetPayload(mc)
	c.GetState(mc)
}

// Check if a node has voted on the proposal
func (c *Proposal) GetMemberHasVoted(mc *batch.MultiCaller, out *bool, address common.Address) {
	core.AddCall(mc, c.mgr, out, "getReceiptHasVoted", c.Details.ID.RawValue, address)
}

// Check if a node has voted in favor of the proposal
func (c *Proposal) GetMemberSupported(mc *batch.MultiCaller, out *bool, address common.Address) {
	core.AddCall(mc, c.mgr, out, "getReceiptSupported", c.Details.ID.RawValue, address)
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
