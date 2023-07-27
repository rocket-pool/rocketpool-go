package trustednode

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// Calls
	networkBalances_getBalancesBlock string = "getBalancesBlock"

	// Transactions
	networkBalances_submitBalances string = "submitBalances"
)

// ===============
// === Structs ===
// ===============

// Binding for Oracle DAO members
type OracleDaoMember struct {
	Index *big.Int
	mgr   *DaoNodeTrusted
}

// Details for Oracle DAO members
type OracleDaoMemberDetails struct {
	// Raw parameters
	Address                   common.Address `json:"address"`
	Exists                    bool           `json:"exists"`
	ID                        string         `json:"id"`
	Url                       string         `json:"url"`
	JoinedTimeRaw             uint64         `json:"joinedTimeRaw"`
	LastProposalTimeRaw       uint64         `json:"lastProposalTimeRaw"`
	RPLBondAmount             *big.Int       `json:"rplBondAmount"`
	UnbondedValidatorCountRaw uint64         `json:"unbondedValidatorCountRaw"`

	// Formatted parameters
	JoinedTime             time.Time `json:"joinedTime"`
	LastProposalTime       time.Time `json:"lastProposalTime"`
	UnbondedValidatorCount uint64    `json:"unbondedValidatorCount"`
}
