package rocketpool

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
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
func GetContractVersion(mc *multicall.MultiCaller, version_Out *uint8, address common.Address) error {
	if versionAbi == nil {
		// Parse ABI using the hardcoded string until the contract is deployed
		abiParsed, err := abi.JSON(strings.NewReader(rocketVersionInterfaceAbiString))
		if err != nil {
			return fmt.Errorf("Could not parse version interface JSON: %w", err)
		}
		versionAbi = &abiParsed
	}

	// Get the contract version
	contract := &core.Contract{
		Contract: bind.NewBoundContract(address, *versionAbi, nil, nil, nil),
	}
	version := new(uint8)
	multicall.AddCall(mc, contract, version, "version")
	return nil
}

// Get the rocketVersion contract binding at the given address
func GetRocketVersionContractForAddress(rp *RocketPool, address common.Address) (*core.Contract, error) {
	if versionAbi == nil {
		// Parse ABI using the hardcoded string until the contract is deployed
		abiParsed, err := abi.JSON(strings.NewReader(rocketVersionInterfaceAbiString))
		if err != nil {
			return nil, fmt.Errorf("Could not parse version interface JSON: %w", err)
		}
		versionAbi = &abiParsed
	}

	return &core.Contract{
		Contract: bind.NewBoundContract(address, *versionAbi, rp.Client, rp.Client, rp.Client),
		Address:  &address,
		ABI:      versionAbi,
		Client:   rp.Client,
	}, nil
}
