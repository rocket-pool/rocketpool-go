package trustednode

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	// Settings
	memberAddressBatchSize = 50
	memberDetailsBatchSize = 20
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAONodeTrusted
type DaoNodeTrusted struct {
	Details  DaoNodeTrustedDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for DaoNodeTrusted
type DaoNodeTrustedDetails struct {
	MemberCount        core.Parameter[uint64] `json:"memberCount"`
	MinimumMemberCount core.Parameter[uint64] `json:"minimumMemberCount"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoNodeTrusted contract binding
func NewDaoNodeTrusted(rp *rocketpool.RocketPool) (*DaoNodeTrusted, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrusted)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted contract: %w", err)
	}

	return &DaoNodeTrusted{
		Details:  DaoNodeTrustedDetails{},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the member count
func (c *DaoNodeTrusted) GetMemberCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.MemberCount.RawValue, "getMemberCount")
}

// Get the minimum member count
func (c *DaoNodeTrusted) GetMinimumMemberCount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.contract, &c.Details.MinimumMemberCount.RawValue, "getMemberMinRequired")
}

// Get all basic details
func (c *DaoNodeTrusted) GetAllDetails(mc *multicall.MultiCaller) {
	c.GetMemberCount(mc)
	c.GetMinimumMemberCount(mc)
}

// Get the list of Oracle DAO member addresses
func (c *DaoNodeTrusted) GetMemberAddresses(memberCount uint64, opts *bind.CallOpts) ([]*common.Address, error) {
	// Run the multicall query for each address
	addresses, err := rocketpool.BatchQuery[common.Address](c.rp,
		memberCount,
		memberAddressBatchSize,
		func(mc *multicall.MultiCaller, index uint64) (*common.Address, error) {
			address := new(common.Address)
			multicall.AddCall(mc, c.contract, address, "getMemberAt", big.NewInt(int64(index)))
			return address, nil
		},
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting Oracle DAO member addresses: %w", err)
	}

	// Return
	return addresses, nil
}

// ====================
// === Transactions ===
// ====================

// Bootstrap a bool setting
func (c *DaoNodeTrusted) BootstrapBool(contractName string, settingPath string, value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "bootstrapSettingBool", opts, contractName, settingPath, value)
}

// Bootstrap a uint setting
func (c *DaoNodeTrusted) BootstrapUint(contractName string, settingPath string, value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "bootstrapSettingUint", opts, contractName, settingPath, value)
}

// Bootstrap a member into the Oracle DAO
func (c *DaoNodeTrusted) BootstrapMember(id string, url string, nodeAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "bootstrapMember", opts, id, url, nodeAddress)
}

// Bootstrap a contract upgrade
func (c *DaoNodeTrusted) BootstrapUpgrade(upgradeType string, contractName string, contractAbi string, contractAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	compressedAbi, err := core.EncodeAbiStr(contractAbi)
	if err != nil {
		return nil, fmt.Errorf("error compressing ABI: %w", err)
	}
	return core.NewTransactionInfo(c.contract, "bootstrapUpgrade", opts, upgradeType, contractName, compressedAbi, contractAddress)
}

// ===================
// === Sub-Getters ===
// ===================

// Get a member's details
func (c *DaoNodeTrusted) GetMemberAt(index uint64, address common.Address, opts *bind.CallOpts) (*OracleDaoMember, error) {
	// Create the member and get details via a multicall query
	member := NewOracleDaoMember(c, index, address)
	err := c.rp.Query(func(mc *multicall.MultiCaller) {
		member.GetAllDetails(mc)
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting Oracle DAO member %d: %w", index, err)
	}

	// Return
	return member, nil
}

// Get the details for all members
func (c *DaoNodeTrusted) GetAllMembers(addresses []*common.Address, opts *bind.CallOpts) ([]*OracleDaoMember, error) {
	// Run the multicall query for each lot
	members, err := rocketpool.BatchQuery[OracleDaoMember](c.rp,
		uint64(len(addresses)),
		memberDetailsBatchSize,
		func(mc *multicall.MultiCaller, index uint64) (*OracleDaoMember, error) {
			member := NewOracleDaoMember(c, index, *addresses[index])
			member.GetAllDetails(mc)
			return member, nil
		},
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting all Oracle DAO member details: %w", err)
	}

	// Return
	return members, nil
}
