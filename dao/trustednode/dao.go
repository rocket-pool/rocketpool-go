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

// ====================
// === Transactions ===
// ====================

// Bootstrap a bool setting
func (c *DaoNodeTrusted) BootstrapBool(contractName rocketpool.ContractName, settingPath string, value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "bootstrapSettingBool", opts, contractName, settingPath, value)
}

// Bootstrap a uint setting
func (c *DaoNodeTrusted) BootstrapUint(contractName rocketpool.ContractName, settingPath string, value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "bootstrapSettingUint", opts, contractName, settingPath, value)
}

// Bootstrap a member into the Oracle DAO
func (c *DaoNodeTrusted) BootstrapMember(id string, url string, nodeAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.contract, "bootstrapMember", opts, id, url, nodeAddress)
}

// Bootstrap a contract upgrade
func (c *DaoNodeTrusted) BootstrapUpgrade(upgradeType string, contractName rocketpool.ContractName, contractAbi string, contractAddress common.Address, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	compressedAbi, err := core.EncodeAbiStr(contractAbi)
	if err != nil {
		return nil, fmt.Errorf("error compressing ABI: %w", err)
	}
	return core.NewTransactionInfo(c.contract, "bootstrapUpgrade", opts, upgradeType, contractName, compressedAbi, contractAddress)
}

// =================
// === Addresses ===
// =================

// Get an Oracle DAO member address by index
func (c *DaoNodeTrusted) GetMemberAddress(mc *multicall.MultiCaller, address_Out *common.Address, index uint64) {
	multicall.AddCall(mc, c.contract, address_Out, "getMemberAt", big.NewInt(int64(index)))
}

// Get the list of Oracle DAO member addresses.
// Use GetMemberCount() for the memberCount parameter.
func (c *DaoNodeTrusted) GetMemberAddresses(memberCount uint64, opts *bind.CallOpts) ([]common.Address, error) {
	addresses := make([]common.Address, memberCount)

	// Run the multicall query for each address
	err := c.rp.BatchQuery(int(memberCount), c.rp.AddressBatchSize,
		func(mc *multicall.MultiCaller, index int) error {
			c.GetMemberAddress(mc, &addresses[index], uint64(index))
			return nil
		},
		opts,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting Oracle DAO member addresses: %w", err)
	}

	// Return
	return addresses, nil
}
