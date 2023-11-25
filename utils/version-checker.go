package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/hashicorp/go-version"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

func GetCurrentVersion(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*version.Version, error) {
	nodeStaking, err := rp.GetContract(rocketpool.ContractName_RocketNodeStaking)
	if err != nil {
		return nil, fmt.Errorf("error getting node staking contract: %w", err)
	}
	nodeMgr, err := rp.GetContract(rocketpool.ContractName_RocketNodeManager)
	if err != nil {
		return nil, fmt.Errorf("error getting node manager contract: %w", err)
	}

	nodeStakingVersion := nodeStaking.Version
	nodeMgrVersion := nodeMgr.Version

	// Check for v1.3 (Houston)
	if nodeMgrVersion > 3 {
		return version.NewSemver("1.3.0")
	}

	// Check for v1.2 (Atlas)
	if nodeStakingVersion > 3 {
		return version.NewSemver("1.2.0")
	}

	// Check for v1.1 (Redstone)
	if nodeMgrVersion > 1 {
		return version.NewSemver("1.1.0")
	}

	// v1.0 (Classic)
	return version.NewSemver("1.0.0")
}
