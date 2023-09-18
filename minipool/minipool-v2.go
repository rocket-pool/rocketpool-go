package minipool

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

const (
	minipoolV2EncodedAbi string = "eJzdWd1v2jAQ/1cqnvvUaVPVt3ZdpUnrVEG7PVQVcpIDLIyN7HMYqva/7xwgHyRAKHGT7qkNXO5+97vzfZjn1x6TSi5nypre1YgJA+c9LucW6fH5lf6N4A9EvSvUNvkGQUsmHpdz6F31WBRpMKZ33pNs5j4YaTWjJyx+/fc8pyi1UdBk6fni85dMEyNEEjNdG4G36EJOf8qaXlKBbzgBfQtzZTiS3lQUYiAMzmSTJJFsaAt2TiFKqgiuGyTLGtBN6kOFTNwwwWRYFQRv4fzNcRJptmDiQauQ2PUfWFQfNfc3ZMm3U1SNJ1gi5BiKmeARQ6UfbDCFZWZtJVfDwV0KB3wsGVoNb9X56SLTGq1KwS1D1lcKt1SS5DtHdcvpRZraXzVEFCZOb73B7+OT5Z5LPleKjhQYZNNTjlTTkAahtkHg/5DPYBaAbuagH3QuceqXar4pOVuXGRAKJlpThHLpyaXE1NOcTm21V6kP2TsxaMOVK07KYs7DfS6VnCF1zk24t8gCLjgunWYOi0xyZGWIztAOHGPAwYapPUhA2tlZmpebF/ziuuNkn6+a3B5oASGqwpJ83ihFN0KF08MRKyRPdeI0BulxlZttISpZK9WW4c7iUvQmXxVaDvZ6ak4s1j8U67e8n4qfbjhOWd6DrhSK6hg0BOkO2szDEpx1NLIhvTPI+kCCUQeBrSm7Nobmzi6cwyeTbrAdoyuHrJN0bUC13B0K8B7d0pyW+QO1q+2mJQtFtnIsPqITDKNCRyl1hbUYve7WHpp4CuRU+klj8pwtWSDgaG+3Nq9hrQW2noYDG+v+DXV4eEXNuJJZwTpMVk2mMu02O0oetOukAzQZ40y3EcxM/KgercdxP9pDJgdu/W6ljnYxw02JjeaSBDC9Sly96MFIxA1qHliEdfO+ltGd1xQqWfRaRkstahjsvBHOp7kIrSAYbuIaTJju1PZ2ok9uAmnbp0I6GCViX/VKKF95HNN8lAxKXvM3VBI1C/Gsr8Kpu05Qmo3huxMasSTimxzQeYFjpqK256p6fBERVDdsSP6dfNdb8liJ6BYEjAnILoePUyhhcZLC4283N+b6SgigxTW5Amv0hvx/Zu1pPtYs8n+H/5F/pu5DCDyu5qjOvM2ECFxa1pTXK3M7+0Yxcp5mldyhCtjWrXLjG19hah7S+Idcjium52w+pFb+gxAYzB0bDzSMD1lq6QK4DpL3u1990BBzKhNdw/VtNAKSiaElYC//AGQZTdM="
)

// The decoded ABI for v2 minipools
var minipoolV2Abi *abi.ABI

// ===============
// === Structs ===
// ===============

type MinipoolV2 struct {
	*minipoolCommon
	*MinipoolV2Details
}

type MinipoolV2Details struct {
}

// ====================
// === Constructors ===
// ====================

// Create new minipool contract
func newMinipool_v2(rp *rocketpool.RocketPool, address common.Address) (*MinipoolV2, error) {
	var contract *core.Contract
	var err error
	if minipoolV2Abi == nil {
		// Get contract
		contract, err = rp.CreateMinipoolContractFromEncodedAbi(address, minipoolV2EncodedAbi)
	} else {
		contract, err = rp.CreateMinipoolContractFromAbi(address, minipoolV2Abi)
	}
	if err != nil {
		return nil, err
	} else if minipoolV2Abi == nil {
		minipoolV2Abi = contract.ABI
	}

	// Create the base binding
	base, err := newMinipoolCommonFromVersion(rp, contract, 2)
	if err != nil {
		return nil, fmt.Errorf("error creating minipool base: %w", err)
	}

	// Create and return
	return &MinipoolV2{
		minipoolCommon:    base,
		MinipoolV2Details: &MinipoolV2Details{},
	}, nil
}

// Get the minipool as a v2 minipool if it implements the required methods
func GetMinipoolAsV2(mp IMinipool) (*MinipoolV2, bool) {
	castedMp, ok := mp.(*MinipoolV2)
	if ok {
		return castedMp, true
	}
	return nil, false
}

// =============
// === Calls ===
// =============

// Query all of the minipool details
func (c *MinipoolV2) QueryAllDetails(mc *batch.MultiCaller) {
	c.minipoolCommon.QueryAllDetails(mc)
}

// ====================
// === Transactions ===
// ====================

// Get info for distributing the minipool's ETH balance to the node operator and rETH staking pool.
// !!! WARNING !!!
// DO NOT CALL THIS until the minipool's validator has exited from the Beacon Chain
// and the balance has been deposited into the minipool!
func (c *MinipoolV2) DistributeBalance(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "distributeBalance", opts)
}

// Get info for distributing the minipool's ETH balance to the node operator and rETH staking pool,
// then finalising the minipool
// !!! WARNING !!!
// DO NOT CALL THIS until the minipool's validator has exited from the Beacon Chain
// and the balance has been deposited into the minipool!
func (c *MinipoolV2) DistributeBalanceAndFinalise(opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.Contract, "distributeBalanceAndFinalise", opts)
}
