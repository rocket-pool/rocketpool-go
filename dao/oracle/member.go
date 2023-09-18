package oracle

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for Oracle DAO members
type OracleDaoMember struct {
	Details  OracleDaoMemberDetails
	contract *core.Contract
}

// Details for Oracle DAO members
type OracleDaoMemberDetails struct {
	Address                common.Address            `json:"address"`
	Exists                 bool                      `json:"exists"`
	ID                     string                    `json:"id"`
	Url                    string                    `json:"url"`
	InvitedTime            core.Parameter[time.Time] `json:"invitedTime"`
	JoinedTime             core.Parameter[time.Time] `json:"joinedTime"`
	ReplacedTime           core.Parameter[time.Time] `json:"replacedTime"`
	LeftTime               core.Parameter[time.Time] `json:"leftTime"`
	LastProposalTime       core.Parameter[time.Time] `json:"lastProposalTime"`
	RPLBondAmount          *big.Int                  `json:"rplBondAmount"`
	ReplacementAddress     common.Address            `json:"replacementAddress"`
	IsChallenged           bool                      `json:"isChallenged"`
	UnbondedValidatorCount core.Parameter[uint64]    `json:"unbondedValidatorCount"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new OracleDaoMember instance
func NewOracleDaoMember(rp *rocketpool.RocketPool, address common.Address) (*OracleDaoMember, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrusted)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted contract: %w", err)
	}

	return &OracleDaoMember{
		Details: OracleDaoMemberDetails{
			Address: address,
		},
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Check whether or not the member exists
func (c *OracleDaoMember) GetExists(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.Exists, "getMemberIsValid", c.Details.Address)
}

// Get the member's ID
func (c *OracleDaoMember) GetID(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.ID, "getMemberID", c.Details.Address)
}

// Get the member's URL
func (c *OracleDaoMember) GetUrl(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.Url, "getMemberUrl", c.Details.Address)
}

// Get the time the member was invited to the Oracle DAO
func (c *OracleDaoMember) GetInvitedTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.InvitedTime.RawValue, "getMemberProposalExecutedTime", "invited", c.Details.Address)
}

// Get the time the member joined the Oracle DAO
func (c *OracleDaoMember) GetJoinedTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.JoinedTime.RawValue, "getMemberJoinedTime", c.Details.Address)
}

// Get the time the member's address was replaced
func (c *OracleDaoMember) GetReplacedTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.ReplacedTime.RawValue, "getMemberProposalExecutedTime", "replace", c.Details.Address)
}

// Get the time the member voluntarily left the Oracle DAO
func (c *OracleDaoMember) GetLeftTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.LeftTime.RawValue, "getMemberProposalExecutedTime", "leave", c.Details.Address)
}

// Get the time the member last made a proposal
func (c *OracleDaoMember) GetLastProposalTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.LastProposalTime.RawValue, "getMemberLastProposalTime", c.Details.Address)
}

// Get the member's RPL bond amount
func (c *OracleDaoMember) GetRplBondAmount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.RPLBondAmount, "getMemberRPLBondAmount", c.Details.Address)
}

// Get the member's replacement address if a replace proposal is pending
func (c *OracleDaoMember) GetReplacementAddress(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.ReplacementAddress, "getMemberReplacedAddress", "new", c.Details.Address)
}

// Check if the member has been challenged
func (c *OracleDaoMember) GetIsChallenged(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.IsChallenged, "getMemberIsChallenged", c.Details.Address)
}

// Get the member's unbonded validator count (defunct; will never be above 0)
func (c *OracleDaoMember) GetUnbondedValidatorCount(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Details.UnbondedValidatorCount.RawValue, "getMemberUnbondedValidatorCount", c.Details.Address)
}

// Get all basic details
func (c *OracleDaoMember) GetAllDetails(mc *batch.MultiCaller) {
	c.GetExists(mc)
	c.GetID(mc)
	c.GetUrl(mc)
	c.GetInvitedTime(mc)
	c.GetJoinedTime(mc)
	c.GetReplacedTime(mc)
	c.GetLeftTime(mc)
	c.GetLastProposalTime(mc)
	c.GetRplBondAmount(mc)
	c.GetReplacementAddress(mc)
	c.GetIsChallenged(mc)
	c.GetUnbondedValidatorCount(mc)
}
