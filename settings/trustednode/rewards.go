package trustednode

import (
	"fmt"
	"math/big"

	"github.com/rocket-pool/rocketpool-go/core"
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
	Details  DaoNodeTrustedSettingsProposalsDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoNodeTrustedSettingsRewards contract binding
func NewDaoNodeTrustedSettingsRewards(rp *rocketpool.RocketPool) (*DaoNodeTrustedSettingsRewards, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAONodeTrustedSettingsRewards)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted settings rewards contract: %w", err)
	}

	return &DaoNodeTrustedSettingsRewards{
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get whether or not the provided rewards network is enabled
func (c *DaoNodeTrustedSettingsProposals) GetNetworkEnabled(mc *multicall.MultiCaller, enabled_Out *bool, network *big.Int) {
	multicall.AddCall(mc, c.contract, enabled_Out, "getNetworkEnabled", network)
}
