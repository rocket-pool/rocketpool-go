package rocketpool

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/v2/core"
)

const (
	rocketVersionInterfaceAbiString string = `[
		{
		  "inputs": [],
		  "name": "version",
		  "outputs": [
			{
			  "internalType": "uint8",
			  "name": "",
			  "type": "uint8"
			}
		  ],
		  "stateMutability": "view",
		  "type": "function"
		}
	]`
)

var versionAbi *abi.ABI

// Get the version of the given contract
func GetContractVersion(mc *batch.MultiCaller, version_Out *uint8, address common.Address) error {
	if versionAbi == nil {
		// Parse ABI using the hardcoded string until the contract is deployed
		abiParsed, err := abi.JSON(strings.NewReader(rocketVersionInterfaceAbiString))
		if err != nil {
			return fmt.Errorf("error parsing version interface JSON: %w", err)
		}
		versionAbi = &abiParsed
	}

	// Get the contract version
	contract := &core.Contract{
		Contract: &eth.Contract{
			ContractImpl: bind.NewBoundContract(address, *versionAbi, nil, nil, nil),
			Address:      address,
			ABI:          versionAbi,
		},
		Version: 0,
	}
	core.AddCall(mc, contract, version_Out, "version")
	return nil
}

// Get the rocketVersion contract binding at the given address
func GetRocketVersionContractForAddress(rp *RocketPool, address common.Address) (*core.Contract, error) {
	if versionAbi == nil {
		// Parse ABI using the hardcoded string until the contract is deployed
		abiParsed, err := abi.JSON(strings.NewReader(rocketVersionInterfaceAbiString))
		if err != nil {
			return nil, fmt.Errorf("error parsing version interface JSON: %w", err)
		}
		versionAbi = &abiParsed
	}

	return &core.Contract{
		Contract: &eth.Contract{
			ContractImpl: bind.NewBoundContract(address, *versionAbi, rp.Client, rp.Client, rp.Client),
			Address:      address,
			ABI:          versionAbi,
		},
	}, nil
}
