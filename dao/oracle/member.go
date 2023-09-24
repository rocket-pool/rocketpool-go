package oracle

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ===============
// === Structs ===
// ===============

// Binding for Oracle DAO members
type OracleDaoMember struct {
	// The address of this member
	Address common.Address

	// True if this member exists (is part of the Oracle DAO)
	Exists *core.SimpleField[bool]

	// The member's ID
	ID *core.SimpleField[string]

	// The member's URL
	Url *core.SimpleField[string]

	// The time the member was invited to the Oracle DAO
	InvitedTime *core.FormattedUint256Field[time.Time]

	// The time the member joined the Oracle DAO
	JoinedTime *core.FormattedUint256Field[time.Time]

	// The time the member's address was replaced
	ReplacedTime *core.FormattedUint256Field[time.Time]

	// The time the member voluntarily left the Oracle DAO
	LeftTime *core.FormattedUint256Field[time.Time]

	// The time the member last made a proposal
	LastProposalTime *core.FormattedUint256Field[time.Time]

	// The member's RPL bond amount
	RplBondAmount *core.SimpleField[*big.Int]

	// The member's replacement address, if a replace proposal is pending
	ReplacementAddress *core.SimpleField[common.Address]

	// True if the member has an active challenge raised against it
	IsChallenged *core.SimpleField[bool]

	// === Internal fields ===
	dnt *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new OracleDaoMember instance
func NewOracleDaoMember(rp *rocketpool.RocketPool, address common.Address) (*OracleDaoMember, error) {
	// Create the contract
	dnt, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrusted)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted contract: %w", err)
	}

	return &OracleDaoMember{
		Address:            address,
		Exists:             core.NewSimpleField[bool](dnt, "getMemberIsValid", address),
		ID:                 core.NewSimpleField[string](dnt, "getMemberID", address),
		Url:                core.NewSimpleField[string](dnt, "getMemberUrl", address),
		InvitedTime:        core.NewFormattedUint256Field[time.Time](dnt, "getMemberProposalExecutedTime", "invited", address),
		JoinedTime:         core.NewFormattedUint256Field[time.Time](dnt, "getMemberJoinedTime", address),
		ReplacedTime:       core.NewFormattedUint256Field[time.Time](dnt, "getMemberProposalExecutedTime", "replace", address),
		LeftTime:           core.NewFormattedUint256Field[time.Time](dnt, "getMemberProposalExecutedTime", "leave", address),
		LastProposalTime:   core.NewFormattedUint256Field[time.Time](dnt, "getMemberLastProposalTime", address),
		RplBondAmount:      core.NewSimpleField[*big.Int](dnt, "getMemberRPLBondAmount", address),
		ReplacementAddress: core.NewSimpleField[common.Address](dnt, "getMemberReplacedAddress", "new", address),
		IsChallenged:       core.NewSimpleField[bool](dnt, "getMemberIsChallenged", address),

		dnt: dnt,
	}, nil
}
