package trustednode

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

const (
	// Calls
	daoNodeTrusted_getMemberIsValid    string = "getMemberIsValid"
	daoNodeTrusted_getMemberID         string = "getMemberID"
	daoNodeTrusted_getMemberUrl        string = "getMemberUrl"
	daoNodeTrusted_getMemberJoinedTime string = "getMemberJoinedTime"

	// Transactions
	networkBalances_submitBalances string = "submitBalances"
)

// ===============
// === Structs ===
// ===============

// Binding for Oracle DAO members
type OracleDaoMember struct {
	Details OracleDaoMemberDetails
	index   *big.Int
	address common.Address
	mgr     *DaoNodeTrusted
}

// Details for Oracle DAO members
type OracleDaoMemberDetails struct {
	// Raw parameters
	IndexRaw                  *big.Int       `json:"indexRaw"`
	Address                   common.Address `json:"address"`
	Exists                    bool           `json:"exists"`
	ID                        string         `json:"id"`
	Url                       string         `json:"url"`
	JoinedTimeRaw             uint64         `json:"joinedTimeRaw"`
	LastProposalTimeRaw       uint64         `json:"lastProposalTimeRaw"`
	RPLBondAmount             *big.Int       `json:"rplBondAmount"`
	UnbondedValidatorCountRaw uint64         `json:"unbondedValidatorCountRaw"`

	// Formatted parameters
	Index                  uint64    `json:"index"`
	JoinedTime             time.Time `json:"joinedTime"`
	LastProposalTime       time.Time `json:"lastProposalTime"`
	UnbondedValidatorCount uint64    `json:"unbondedValidatorCount"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new OracleDaoMember instance
func NewOracleDaoMember(mgr *DaoNodeTrusted, index uint64, address common.Address) *OracleDaoMember {
	return &OracleDaoMember{
		Details: OracleDaoMemberDetails{},
		index:   big.NewInt(int64(index)),
		address: address,
		mgr:     mgr,
	}
}

// ===================
// === Raw Getters ===
// ===================

// Check whether or not the member exists
func (c *OracleDaoMember) GetMemberExists(opts *bind.CallOpts) (bool, error) {
	return rocketpool.Call[bool](c.mgr.contract, opts, daoNodeTrusted_getMemberIsValid, c.address)
}

// Get the member's ID
func (c *OracleDaoMember) GetMemberID(opts *bind.CallOpts) (string, error) {
	return rocketpool.Call[string](c.mgr.contract, opts, daoNodeTrusted_getMemberID, c.address)
}

// Get the member's URL
func (c *OracleDaoMember) GetMemberUrl(opts *bind.CallOpts) (string, error) {
	return rocketpool.Call[string](c.mgr.contract, opts, daoNodeTrusted_getMemberUrl, c.address)
}

// Get the time the member joined the Oracle DAO
func (c *OracleDaoMember) GetMemberJoinedTimeRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, daoNodeTrusted_getMemberJoinedTime, c.address)
}

// Get the time the member joined the Oracle DAO
func (c *OracleDaoMember) GetMemberJoinedTimeRaw(opts *bind.CallOpts) (*big.Int, error) {
	return rocketpool.Call[*big.Int](c.mgr.contract, opts, daoNodeTrusted_getMemberJoinedTime, c.address)
}

// =========================
// === Formatted Getters ===
// =========================

// Get the time the member joined the Oracle DAO
func (c *OracleDaoMember) GetMemberJoinedTime(opts *bind.CallOpts) (time.Time, error) {
	return ConvertToTime(c.GetMemberJoinedTimeRaw, opts)
}

func ConvertRawToTime(rawGetter func(opts *bind.CallOpts) (*big.Int, error), opts *bind.CallOpts) (time.Time, error) {
	raw, err := rawGetter(opts)
	if err != nil {
		return time.Time{}, err
	}
	return ConvertToTime(raw), nil
}

func ConvertToTime(raw *big.Int) time.Time {
	return time.Unix(raw.Int64(), 0)
}
