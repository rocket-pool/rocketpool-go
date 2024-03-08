package protocol

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
	strutils "github.com/rocket-pool/rocketpool-go/utils/strings"
)

// ===============
// === Structs ===
// ===============

// Binding for proposals
type ProtocolDaoProposal struct {
	// The proposal's ID
	ID uint64

	// The address of the node that created the proposal
	ProposerAddress *core.SimpleField[common.Address]

	// The block number that was used as the reference when
	// generating the voting tree for this proposal's pollard
	TargetBlock *core.FormattedUint256Field[uint32]

	// The message provided with the proposal
	Message *core.SimpleField[string]

	// The length of time from proposal creation where challenges can be responded to
	ChallengeWindow *core.FormattedUint256Field[time.Duration]

	// The time when nodes can start voting on the proposal (start of phase 1)
	VotingStartTime *core.FormattedUint256Field[time.Time]

	// The time that marks the end of phase 1 and the start of phase 2
	Phase1EndTime *core.FormattedUint256Field[time.Time]

	// The time that marks the end of phase 2
	Phase2EndTime *core.FormattedUint256Field[time.Time]

	// The time the proposal expires on, where it can no longer be executed if successful
	ExpiryTime *core.FormattedUint256Field[time.Time]

	// The time the proposal was created
	CreatedTime *core.FormattedUint256Field[time.Time]

	// The amount of voting power required for the proposal to be decided (the quorum)
	VotingPowerRequired *core.FormattedUint256Field[float64]

	// The amount of voting power that has voted in favor of the proposal
	VotingPowerFor *core.FormattedUint256Field[float64]

	// The amount of voting power that has voted against the proposal
	VotingPowerAgainst *core.FormattedUint256Field[float64]

	// The amount of voting power that has abstained from voting on the proposal
	VotingPowerAbstained *core.FormattedUint256Field[float64]

	// The amount of voting power that has voted to veto the proposal
	VotingPowerToVeto *core.FormattedUint256Field[float64]

	// Whether or not the proposal has been destroyed
	IsDestroyed *core.SimpleField[bool]

	// Whether or not the proposal has been finalized
	IsFinalized *core.SimpleField[bool]

	// Whether or not the proposal has been executed
	IsExecuted *core.SimpleField[bool]

	// Whether or not the proposal has been vetoed
	IsVetoed *core.SimpleField[bool]

	// The amount of voting power required to veto the proposal
	VetoQuorum *core.FormattedUint256Field[float64]

	// The proposal's payload
	Payload *core.SimpleField[[]byte]

	// The proposal's state
	State *core.FormattedUint8Field[types.ProtocolDaoProposalState]

	// The RPL bond locked by the proposer as part of submitting this proposal
	ProposalBond *core.SimpleField[*big.Int]

	// The RPL bond locked by a challenger as part of submitting a challenge against this proposal
	ChallengeBond *core.SimpleField[*big.Int]

	// The index of the tree node that challenged and not responded to, if this proposal was defeated
	DefeatIndex *core.FormattedUint256Field[uint64]

	// === Internal fields ===
	idBig *big.Int
	rp    *rocketpool.RocketPool
	dpp   *core.Contract
	dpps  *core.Contract
	dpv   *core.Contract
	txMgr *eth.TransactionManager
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProtocolDaoProposal contract binding
func NewProtocolDaoProposal(rp *rocketpool.RocketPool, id uint64) (*ProtocolDaoProposal, error) {
	// Create the contracts
	dpp, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocolProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO proposal contract: %w", err)
	}
	dpps, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocolProposals)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO proposals contract: %w", err)
	}
	dpv, err := rp.GetContract(rocketpool.ContractName_RocketDAOProtocolVerifier)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO verifier contract: %w", err)
	}

	idBig := big.NewInt(0).SetUint64(id)
	return &ProtocolDaoProposal{
		ID:                   id,
		ProposerAddress:      core.NewSimpleField[common.Address](dpp, "getProposer", idBig),
		TargetBlock:          core.NewFormattedUint256Field[uint32](dpp, "getProposalBlock", idBig),
		Message:              core.NewSimpleField[string](dpp, "getMessage", idBig),
		ChallengeWindow:      core.NewFormattedUint256Field[time.Duration](dpv, "getChallengePeriod", idBig),
		VotingStartTime:      core.NewFormattedUint256Field[time.Time](dpp, "getStart", idBig),
		Phase1EndTime:        core.NewFormattedUint256Field[time.Time](dpp, "getPhase1End", idBig),
		Phase2EndTime:        core.NewFormattedUint256Field[time.Time](dpp, "getPhase2End", idBig),
		ExpiryTime:           core.NewFormattedUint256Field[time.Time](dpp, "getExpires", idBig),
		CreatedTime:          core.NewFormattedUint256Field[time.Time](dpp, "getCreated", idBig),
		VotingPowerRequired:  core.NewFormattedUint256Field[float64](dpp, "getVotingPowerRequired", idBig),
		VotingPowerFor:       core.NewFormattedUint256Field[float64](dpp, "getVotingPowerFor", idBig),
		VotingPowerAgainst:   core.NewFormattedUint256Field[float64](dpp, "getVotingPowerAgainst", idBig),
		VotingPowerAbstained: core.NewFormattedUint256Field[float64](dpp, "getVotingPowerAbstained", idBig),
		VotingPowerToVeto:    core.NewFormattedUint256Field[float64](dpp, "getVotingPowerVeto", idBig),
		IsDestroyed:          core.NewSimpleField[bool](dpp, "getDestroyed", idBig),
		IsFinalized:          core.NewSimpleField[bool](dpp, "getFinalised", idBig),
		IsExecuted:           core.NewSimpleField[bool](dpp, "getExecuted", idBig),
		IsVetoed:             core.NewSimpleField[bool](dpp, "getVetoed", idBig),
		VetoQuorum:           core.NewFormattedUint256Field[float64](dpp, "getProposalVetoQuorum", idBig),
		Payload:              core.NewSimpleField[[]byte](dpp, "getPayload", idBig),
		State:                core.NewFormattedUint8Field[types.ProtocolDaoProposalState](dpp, "getState", idBig),
		ProposalBond:         core.NewSimpleField[*big.Int](dpv, "getProposalBond", idBig),
		ChallengeBond:        core.NewSimpleField[*big.Int](dpv, "getChallengeBond", idBig),
		DefeatIndex:          core.NewFormattedUint256Field[uint64](dpv, "getDefeatIndex", idBig),

		idBig: idBig,
		rp:    rp,
		dpp:   dpp,
		dpps:  dpps,
		dpv:   dpv,
		txMgr: rp.GetTransactionManager(),
	}, nil
}

// =============
// === Calls ===
// =============

// Get the option that the address voted on for the proposal, and whether or not it's voted yet
func (p *ProtocolDaoProposal) GetAddressVoteDirection(mc *batch.MultiCaller, address common.Address) func() types.VoteDirection {
	out := new(uint8)
	core.AddCall(mc, p.dpp, out, "getReceiptDirection", p.idBig, address)

	return func() types.VoteDirection {
		return types.VoteDirection(*out)
	}
}

// Get the tree node of the proposal's voting tree at the given index
func (p *ProtocolDaoProposal) GetTreeNode(mc *batch.MultiCaller, nodeIndex uint64) func() types.VotingTreeNode {
	type nodeRaw struct {
		Sum  *big.Int `json:"sum"`
		Hash [32]byte `json:"hash"`
	}
	out := new(nodeRaw)
	core.AddCallRaw(mc, p.dpv, out, "getNode", p.idBig, big.NewInt(int64(nodeIndex)))

	return func() types.VotingTreeNode {
		return types.VotingTreeNode{
			Sum:  out.Sum,
			Hash: common.BytesToHash(out.Hash[:]),
		}
	}
}

// Get the state of a challenge on a proposal and tree node index
func (p *ProtocolDaoProposal) GetChallengeState(mc *batch.MultiCaller, index uint64) func() types.ChallengeState {
	out := new(uint8)
	core.AddCall(mc, p.dpv, out, "getChallengeState", p.idBig, big.NewInt(int64(index)))

	return func() types.ChallengeState {
		return types.ChallengeState(*out)
	}
}

// ====================
// === Transactions ===
// ====================

// Get info for voting on a proposal
func (p *ProtocolDaoProposal) Vote(voteDirection types.VoteDirection, votingPower *big.Int, nodeIndex uint64, witness []types.VotingTreeNode, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return p.txMgr.CreateTransactionInfo(p.dpp.Contract, "vote", opts, p.idBig, voteDirection, votingPower, big.NewInt(int64(nodeIndex)), witness)
}

// Get info for overriding a delegate's vote during phase 2
func (p *ProtocolDaoProposal) OverrideVote(voteDirection types.VoteDirection, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return p.txMgr.CreateTransactionInfo(p.dpp.Contract, "overrideVote", opts, p.idBig, voteDirection)
}

// Get info for finalizing a vetoed proposal by burning the proposer's bond
func (p *ProtocolDaoProposal) Finalize(opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return p.txMgr.CreateTransactionInfo(p.dpp.Contract, "finalise", opts, p.idBig)
}

// Get info for executing a proposal
func (p *ProtocolDaoProposal) Execute(opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return p.txMgr.CreateTransactionInfo(p.dpp.Contract, "execute", opts, p.idBig)
}

// Get info for defeaing a proposal if the proposer fails to respond to a challenge within the challenge window, providing the node index that wasn't responded to
func (p *ProtocolDaoProposal) Defeat(index uint64, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return p.txMgr.CreateTransactionInfo(p.dpv.Contract, "defeatProposal", opts, p.idBig, big.NewInt(int64(index)))
}

// Get info for challenging the proposal at a specific tree node index, providing a Merkle proof of the node as well
func (p *ProtocolDaoProposal) CreateChallenge(index uint64, node types.VotingTreeNode, witness []types.VotingTreeNode, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return p.txMgr.CreateTransactionInfo(p.dpv.Contract, "createChallenge", opts, p.idBig, big.NewInt((int64(index))), node, witness)
}

// Get info for submitting the Merkle root for the proposal at the specific index in response to a challenge
func (p *ProtocolDaoProposal) SubmitRoot(index uint64, treeNodes []types.VotingTreeNode, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return p.txMgr.CreateTransactionInfo(p.dpv.Contract, "submitRoot", opts, p.idBig, big.NewInt((int64(index))), treeNodes)
}

// Get info for claiming any RPL bond refunds or rewards for a proposal, as a challenger
func (p *ProtocolDaoProposal) ClaimBondChallenger(indices []uint64, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	indicesBig := make([]*big.Int, len(indices))
	for i, index := range indices {
		indicesBig[i] = big.NewInt(int64(index))
	}
	return p.txMgr.CreateTransactionInfo(p.dpv.Contract, "claimBondChallenger", opts, p.idBig, indicesBig)
}

// Get info for claiming any RPL bond refunds or rewards for a proposal, as the proposer
func (p *ProtocolDaoProposal) ClaimBondProposer(indices []uint64, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	indicesBig := make([]*big.Int, len(indices))
	for i, index := range indices {
		indicesBig[i] = big.NewInt(int64(index))
	}
	return p.txMgr.CreateTransactionInfo(p.dpv.Contract, "claimBondProposer", opts, p.idBig, indicesBig)
}

// =============
// === Utils ===
// =============

// Get a proposal's payload as a human-readable string
func (p *ProtocolDaoProposal) GetProposalPayloadString() (string, error) {
	// Get proposal DAO contract ABI
	daoContractAbi := p.dpps.ABI

	// Get proposal payload method
	payload := p.Payload.Get()
	if len(payload) == 0 {
		return "", fmt.Errorf("payload has not been queried yet")
	}
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
