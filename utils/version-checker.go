package utils

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/hashicorp/go-version"
	"github.com/rocket-pool/rocketpool-go/node"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

func GetCurrentVersion(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*version.Version, error) {

	// TODO: refactor so it's no all atomic calls, it should use a state, or there should be a general "get me the versions for X contracts" function
	// Maybe rp.GetContract() should get the version too?

	// Check for v1.2
	nodeStaking, err := node.NewNodeStaking(rp, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node staking contract: %w", err)
	}
	err = rp.Query(func(mc *multicall.MultiCaller) {
		nodeStaking.GetVersion(mc)
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error checking node staking version: %w", err)
	}
	if nodeStaking.Details.Version > 3 {
		return version.NewSemver("1.2.0")
	}

	// Check for v1.1
	nodeMgrVersion, err := node.GetNodeManagerVersion(rp, opts)
	if err != nil {
		return nil, fmt.Errorf("error checking node manager version: %w", err)
	}
	if nodeMgrVersion > 1 {
		return version.NewSemver("1.1.0")
	}

	// v1.0
	return version.NewSemver("1.0.0")

}
