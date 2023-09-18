package minipool

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// ==================
// === Interfaces ===
// ==================

type IMinipool interface {
	QueryAllDetails(mc *batch.MultiCaller)
	GetMinipoolCommon() *MinipoolCommon
	GetContract() *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Create a minipool binding from an explicit version number
func NewMinipoolFromVersion(rp *rocketpool.RocketPool, address common.Address, version uint8) (IMinipool, error) {
	switch version {
	case 1, 2:
		return newMinipool_v2(rp, address)
	case 3:
		return newMinipool_v3(rp, address)
	default:
		return nil, fmt.Errorf("unexpected minipool contract version [%d]", version)
	}
}

// Create a minipool binding from its address
func CreateMinipoolFromAddress(rp *rocketpool.RocketPool, address common.Address, includeDetails bool, opts *bind.CallOpts) (IMinipool, error) {
	// Get the minipool version
	var version uint8
	results, err := rp.FlexQuery(func(mc *batch.MultiCaller) error {
		return rocketpool.GetContractVersion(mc, &version, address)
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error querying minipool version: %w", err)
	}
	if !results[0] {
		// If it failed, this is a contract on Prater from before version() existed so it's v1
		version = 1
	}

	// Get the minipool
	minipool, err := NewMinipoolFromVersion(rp, address, version)
	if err != nil {
		return nil, fmt.Errorf("error creating minipool: %w", err)
	}

	// Include the details if requested
	if includeDetails {
		err := rp.Query(func(mc *batch.MultiCaller) error {
			minipool.QueryAllDetails(mc)
			return nil
		}, opts)
		if err != nil {
			return nil, fmt.Errorf("error getting minipool %s details: %w", address.Hex(), err)
		}
	}

	return minipool, nil
}

// Create bindings for all minipools from the provided addresses in a standalone call.
// This will use an internal batched multicall invocation to build all of them quickly.
func CreateMinipoolsFromAddresses(rp *rocketpool.RocketPool, addresses []common.Address, includeDetails bool, opts *bind.CallOpts) ([]IMinipool, error) {
	minipoolCount := len(addresses)

	// Get the minipool versions
	versions := make([]uint8, minipoolCount)
	err := rp.FlexBatchQuery(int(minipoolCount), rp.ContractVersionBatchSize,
		func(mc *batch.MultiCaller, index int) error {
			return rocketpool.GetContractVersion(mc, &versions[index], addresses[index])
		},
		func(result bool, index int) error {
			if !result {
				// If it failed, this is a contract on Prater from before version() existed so it's v1
				versions[index] = 1
			}
			return nil
		}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting minipool versions: %w", err)
	}

	// Create the minipools
	minipools := make([]IMinipool, minipoolCount)
	for i := 0; i < int(minipoolCount); i++ {
		minipool, err := NewMinipoolFromVersion(rp, addresses[i], versions[i])
		if err != nil {
			return nil, fmt.Errorf("error creating minipool %d (%s): %w", i, addresses[i].Hex(), err)
		}
		minipools[i] = minipool
	}

	// Include the details if requested
	if includeDetails {
		err := rp.BatchQuery(int(minipoolCount), minipoolBatchSize, func(mc *batch.MultiCaller, index int) error {
			minipools[index].QueryAllDetails(mc)
			return nil
		}, opts)
		if err != nil {
			return nil, fmt.Errorf("error getting minipool details: %w", err)
		}
	}

	return minipools, nil
}
