package trustednode

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// Config
const (
	NetworkEnabledPath string = "rewards.network.enabled"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketDAONodeTrustedSettingsRewards
type DaoNodeTrustedSettingsRewards struct {
	Details                         DaoNodeTrustedSettingsProposalsDetails
	rp                              *rocketpool.RocketPool
	contract                        *core.Contract
	daoNodeTrustedContract          *trustednode.DaoNodeTrusted
	daoNodeTrustedProposalsContract *trustednode.DaoNodeTrustedProposals
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoNodeTrustedSettingsRewards contract binding
func NewDaoNodeTrustedSettingsRewards(rp *rocketpool.RocketPool, daoNodeTrustedContract *trustednode.DaoNodeTrusted, daoNodeTrustedProposalsContract *trustednode.DaoNodeTrustedProposals, opts *bind.CallOpts) (*DaoNodeTrustedSettingsRewards, error) {
	// Create the contract
	contract, err := rp.GetContract("", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted settings rewards contract: %w", err)
	}

	return &DaoNodeTrustedSettingsRewards{
		rp:                              rp,
		contract:                        contract,
		daoNodeTrustedContract:          daoNodeTrustedContract,
		daoNodeTrustedProposalsContract: daoNodeTrustedProposalsContract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get whether or not the provided rewards network is enabled
func (c *DaoNodeTrustedSettingsProposals) GetNetworkEnabled(mc *multicall.MultiCaller, enabled_Out *bool, network *big.Int) {
	multicall.AddCall(mc, c.contract, enabled_Out, "getNetworkEnabled", network)
}
