package proposals

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/v2/core"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for security council proposals
type SecurityCouncilProposal struct {
	*ProposalCommon

	// === Internal fields ===
	rp    *rocketpool.RocketPool
	dsp   *core.Contract
	txMgr *eth.TransactionManager
}

// ====================
// === Constructors ===
// ====================

// Creates a new SecurityCouncilProposal contract binding
func newSecurityCouncilProposal(rp *rocketpool.RocketPool, base *ProposalCommon) (*SecurityCouncilProposal, error) {
	// Create the contract
	dsp, err := rp.GetContract(rocketpool.ContractName_RocketDAOSecurityProposals)
	if err != nil {
		return nil, fmt.Errorf("error getting security council proposals contract: %w", err)
	}

	return &SecurityCouncilProposal{
		ProposalCommon: base,
		rp:             rp,
		dsp:            dsp,
		txMgr:          rp.GetTransactionManager(),
	}, nil
}

// Get a proposal as a security council propsal
func GetProposalAsSecurity(proposal IProposal) (*SecurityCouncilProposal, bool) {
	castedProp, ok := proposal.(*SecurityCouncilProposal)
	if ok {
		return castedProp, true
	}
	return nil, false
}

// =============
// === Calls ===
// =============

// Get the basic details
func (c *SecurityCouncilProposal) QueryAllFields(mc *batch.MultiCaller) {
	eth.QueryAllFields(c.ProposalCommon, mc)
}

// Get the common fields
func (c *SecurityCouncilProposal) Common() *ProposalCommon {
	return c.ProposalCommon
}

// Get the proposal's payload as a string
func (c *SecurityCouncilProposal) GetPayloadAsString() (string, error) {
	return getPayloadAsStringImpl(c.rp, c.dsp, c.Payload.Get())
}

// ====================
// === Transactions ===
// ====================

// Get info for cancelling a proposal
func (c *SecurityCouncilProposal) Cancel(opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.dsp.Contract, "cancel", opts, c.idBig)
}

// Get info for voting on a proposal
func (c *SecurityCouncilProposal) VoteOn(support bool, opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.dsp.Contract, "vote", opts, c.idBig, support)
}

// Get info for executing a proposal
func (c *SecurityCouncilProposal) Execute(opts *bind.TransactOpts) (*eth.TransactionInfo, error) {
	return c.txMgr.CreateTransactionInfo(c.dsp.Contract, "execute", opts, c.idBig)
}
