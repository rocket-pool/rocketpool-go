package state

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/hashicorp/go-version"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

// Container for network contracts
type NetworkContracts struct {
	// Non-RP Utility
	BalanceBatcher     *batch.BalanceBatcher
	MulticallerAddress common.Address
	ElBlockNumber      *big.Int

	// Network version
	Version *version.Version

	// Redstone
	RocketDAONodeTrusted                 *core.Contract
	RocketDAONodeTrustedSettingsMinipool *core.Contract
	RocketDAOProtocolSettingsMinipool    *core.Contract
	RocketDAOProtocolSettingsNetwork     *core.Contract
	RocketDAOProtocolSettingsNode        *core.Contract
	RocketDepositPool                    *core.Contract
	RocketMinipoolManager                *core.Contract
	RocketMinipoolQueue                  *core.Contract
	RocketNetworkBalances                *core.Contract
	RocketNetworkFees                    *core.Contract
	RocketNetworkPrices                  *core.Contract
	RocketNodeDeposit                    *core.Contract
	RocketNodeDistributorFactory         *core.Contract
	RocketNodeManager                    *core.Contract
	RocketNodeStaking                    *core.Contract
	RocketRewardsPool                    *core.Contract
	RocketSmoothingPool                  *core.Contract
	RocketStorage                        *core.Contract
	RocketTokenRETH                      *core.Contract
	RocketTokenRPL                       *core.Contract
	RocketTokenRPLFixedSupply            *core.Contract

	// Atlas
	RocketMinipoolBondReducer *core.Contract
}

type contractArtifacts struct {
	name       string
	address    common.Address
	abiEncoded string
	contract   **core.Contract
}

// Get a new network contracts container
func NewNetworkContracts(rp *rocketpool.RocketPool, multicallerAddress common.Address, balanceBatcherAddress common.Address, opts *bind.CallOpts) (*NetworkContracts, error) {
	// Get the latest block number if it's not provided
	if opts == nil {
		latestElBlock, err := rp.Client.BlockNumber(context.Background())
		if err != nil {
			return nil, fmt.Errorf("error getting latest block number: %w", err)
		}
		opts = &bind.CallOpts{
			BlockNumber: big.NewInt(0).SetUint64(latestElBlock),
		}
	}

	// Create the contract binding
	contracts := &NetworkContracts{
		RocketStorage:      rp.Storage.Contract,
		ElBlockNumber:      opts.BlockNumber,
		MulticallerAddress: multicallerAddress,
	}

	// Create the balance batcher
	var err error
	contracts.BalanceBatcher, err = batch.NewBalanceBatcher(rp.Client, balanceBatcherAddress, 1000, 6)
	if err != nil {
		return nil, err
	}

	// Create the contract wrappers for Redstone
	wrappers := []contractArtifacts{
		{
			name:     "rocketDAONodeTrusted",
			contract: &contracts.RocketDAONodeTrusted,
		}, {
			name:     "rocketDAONodeTrustedSettingsMinipool",
			contract: &contracts.RocketDAONodeTrustedSettingsMinipool,
		}, {
			name:     "rocketDAOProtocolSettingsMinipool",
			contract: &contracts.RocketDAOProtocolSettingsMinipool,
		}, {
			name:     "rocketDAOProtocolSettingsNetwork",
			contract: &contracts.RocketDAOProtocolSettingsNetwork,
		}, {
			name:     "rocketDAOProtocolSettingsNode",
			contract: &contracts.RocketDAOProtocolSettingsNode,
		}, {
			name:     "rocketDepositPool",
			contract: &contracts.RocketDepositPool,
		}, {
			name:     "rocketMinipoolManager",
			contract: &contracts.RocketMinipoolManager,
		}, {
			name:     "rocketMinipoolQueue",
			contract: &contracts.RocketMinipoolQueue,
		}, {
			name:     "rocketNetworkBalances",
			contract: &contracts.RocketNetworkBalances,
		}, {
			name:     "rocketNetworkFees",
			contract: &contracts.RocketNetworkFees,
		}, {
			name:     "rocketNetworkPrices",
			contract: &contracts.RocketNetworkPrices,
		}, {
			name:     "rocketNodeDeposit",
			contract: &contracts.RocketNodeDeposit,
		}, {
			name:     "rocketNodeDistributorFactory",
			contract: &contracts.RocketNodeDistributorFactory,
		}, {
			name:     "rocketNodeManager",
			contract: &contracts.RocketNodeManager,
		}, {
			name:     "rocketNodeStaking",
			contract: &contracts.RocketNodeStaking,
		}, {
			name:     "rocketRewardsPool",
			contract: &contracts.RocketRewardsPool,
		}, {
			name:     "rocketSmoothingPool",
			contract: &contracts.RocketSmoothingPool,
		}, {
			name:     "rocketTokenRETH",
			contract: &contracts.RocketTokenRETH,
		}, {
			name:     "rocketTokenRPL",
			contract: &contracts.RocketTokenRPL,
		}, {
			name:     "rocketTokenRPLFixedSupply",
			contract: &contracts.RocketTokenRPLFixedSupply,
		},
	}

	// Atlas wrappers
	wrappers = append(wrappers, contractArtifacts{
		name:     "rocketMinipoolBondReducer",
		contract: &contracts.RocketMinipoolBondReducer,
	})

	// Create a multicaller
	mc, err := batch.NewMultiCaller(rp.Client, multicallerAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating multicaller: %w", err)
	}

	// Add the address and ABI getters to multicall
	for i, wrapper := range wrappers {
		// Add the address getter
		core.AddCall(mc, contracts.RocketStorage, &wrappers[i].address, "getAddress", [32]byte(crypto.Keccak256Hash([]byte("contract.address"), []byte(wrapper.name))))

		// Add the ABI getter
		core.AddCall(mc, contracts.RocketStorage, &wrappers[i].abiEncoded, "getString", [32]byte(crypto.Keccak256Hash([]byte("contract.abi"), []byte(wrapper.name))))
	}

	// Run the multi-getter
	_, err = mc.FlexibleCall(true, opts)
	if err != nil {
		return nil, fmt.Errorf("error executing multicall for contract retrieval: %w", err)
	}

	// Postprocess the contracts
	for i, wrapper := range wrappers {
		// Decode the ABI
		abi, err := core.DecodeAbi(wrapper.abiEncoded)
		if err != nil {
			return nil, fmt.Errorf("error decoding ABI for %s: %w", wrapper.name, err)
		}

		// Create the contract binding
		contract := &core.Contract{
			Contract: bind.NewBoundContract(wrapper.address, *abi, rp.Client, rp.Client, rp.Client),
			Address:  &wrappers[i].address,
			ABI:      abi,
			Client:   rp.Client,
		}

		// Set the contract in the main wrapper object
		*wrappers[i].contract = contract
	}

	err = contracts.getCurrentVersion(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting network contract version: %w", err)
	}

	return contracts, nil
}

// Get the current version of the network
func (c *NetworkContracts) getCurrentVersion(rp *rocketpool.RocketPool) error {
	opts := &bind.CallOpts{
		BlockNumber: c.ElBlockNumber,
	}

	// Get the contract versions
	var nodeStakingVersion uint8
	var nodeMgrVersion uint8
	err := rp.Query(func(mc *batch.MultiCaller) error {
		rocketpool.GetContractVersion(mc, &nodeStakingVersion, *c.RocketNodeStaking.Address)
		rocketpool.GetContractVersion(mc, &nodeMgrVersion, *c.RocketNodeManager.Address)
		return nil
	}, opts)
	if err != nil {
		return fmt.Errorf("error checking node staking version: %w", err)
	}

	// Check for v1.2
	if nodeStakingVersion > 3 {
		c.Version, err = version.NewSemver("1.2.0")
		return err
	}

	// Check for v1.1
	if err != nil {
		return fmt.Errorf("error checking node manager version: %w", err)
	}
	if nodeMgrVersion > 1 {
		c.Version, err = version.NewSemver("1.1.0")
		return err
	}

	// v1.0
	c.Version, err = version.NewSemver("1.0.0")
	return err
}
