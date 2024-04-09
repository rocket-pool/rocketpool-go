package security

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/v2/core"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for security council members
type SecurityCouncilMember struct {
	// The address of this member
	Address common.Address

	// True if this member exists (is part of the security council)
	Exists *core.SimpleField[bool]

	// The member's ID
	ID *core.SimpleField[string]

	// The time the member was invited to the Oracle DAO
	InvitedTime *core.FormattedUint256Field[time.Time]

	// The time the member joined the Oracle DAO
	JoinedTime *core.FormattedUint256Field[time.Time]

	// The time the member voluntarily left the Oracle DAO
	LeftTime *core.FormattedUint256Field[time.Time]

	// === Internal fields ===
	ds *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new SecurityCouncilMember instance
func NewSecurityCouncilMember(rp *rocketpool.RocketPool, address common.Address) (*SecurityCouncilMember, error) {
	// Create the contract
	ds, err := rp.GetContract(rocketpool.ContractName_RocketDAOSecurity)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO security contract: %w", err)
	}

	return &SecurityCouncilMember{
		Address:     address,
		Exists:      core.NewSimpleField[bool](ds, "getMemberIsValid", address),
		ID:          core.NewSimpleField[string](ds, "getMemberID", address),
		InvitedTime: core.NewFormattedUint256Field[time.Time](ds, "getMemberProposalExecutedTime", "invited", address),
		JoinedTime:  core.NewFormattedUint256Field[time.Time](ds, "getMemberJoinedTime", address),
		LeftTime:    core.NewFormattedUint256Field[time.Time](ds, "getMemberProposalExecutedTime", "leave", address),

		ds: ds,
	}, nil
}
