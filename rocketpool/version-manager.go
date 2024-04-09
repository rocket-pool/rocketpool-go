package rocketpool

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-version"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/v2/core"
)

// Wrapper for legacy contract versions
type LegacyVersionWrapper interface {
	GetVersion() *version.Version
	GetVersionedContractName(contractName string) (string, bool)
	GetEncodedABI(contractName string) string
	GetContractWithAddress(contractName string, address common.Address) (*core.Contract, error)
}

type VersionManager struct {
	V1_0_0     LegacyVersionWrapper
	V1_1_0_RC1 LegacyVersionWrapper
	V1_1_0     LegacyVersionWrapper
	V1_2_0     LegacyVersionWrapper

	rp *RocketPool
}

func NewVersionManager(rp *RocketPool) *VersionManager {
	return &VersionManager{
		V1_0_0:     newLegacyVersionWrapper_v1_0_0(rp),
		V1_1_0_RC1: newLegacyVersionWrapper_v1_1_0_rc1(rp),
		V1_1_0:     newLegacyVersionWrapper_v1_1_0(rp),
		V1_2_0:     newLegacyVersionWrapper_v1_2_0(rp),
		rp:         rp,
	}
}

// Get the contract with the provided name, address, and version wrapper
func getLegacyContractWithAddress(rp *RocketPool, contractName string, address common.Address, m LegacyVersionWrapper) (*core.Contract, error) {

	abiEncoded := m.GetEncodedABI(contractName)
	abi, err := core.DecodeAbi(abiEncoded)
	if err != nil {
		return nil, fmt.Errorf("error decoding contract %s ABI: %w", contractName, err)
	}

	contract := &core.Contract{
		Contract: &eth.Contract{
			ContractImpl: bind.NewBoundContract(address, *abi, rp.Client, rp.Client, rp.Client),
			Address:      address,
			ABI:          abi,
		},
	}

	return contract, nil

}
