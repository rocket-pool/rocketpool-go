package validator

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/go-ssz"
	"github.com/rocket-pool/node-manager-core/beacon"

	"github.com/rocket-pool/rocketpool-go/v2/types"

	"github.com/rocket-pool/rocketpool-go/v2/tests"
)

// Deposit settings
const depositAmount = 16000000000 // gwei

// Deposit data
type depositData struct {
	PublicKey             []byte `ssz-size:"48"`
	WithdrawalCredentials []byte `ssz-size:"32"`
	Amount                uint64
	Signature             []byte `ssz-size:"96"`
}

// Get the validator pubkey
func GetValidatorPubkey(pubkey int) (beacon.ValidatorPubkey, error) {
	if pubkey == 1 {
		return types.HexToValidatorPubkey(tests.ValidatorPubkey)
	} else if pubkey == 2 {
		return types.HexToValidatorPubkey(tests.ValidatorPubkey2)
	} else if pubkey == 3 {
		return types.HexToValidatorPubkey(tests.ValidatorPubkey3)
	} else {
		return beacon.ValidatorPubkey{}, fmt.Errorf("Invalid pubkey index %d", pubkey)
	}
}

// Get the validator deposit signature
func GetValidatorSignature(pubkey int) (beacon.ValidatorSignature, error) {
	if pubkey == 1 {
		return types.HexToValidatorSignature(tests.ValidatorSignature)
	} else if pubkey == 2 {
		return types.HexToValidatorSignature(tests.ValidatorSignature2)
	} else if pubkey == 3 {
		return types.HexToValidatorSignature(tests.ValidatorSignature3)
	} else {
		return beacon.ValidatorSignature{}, fmt.Errorf("Invalid pubkey index %d", pubkey)
	}
}

// Get the validator deposit depositDataRoot
func GetDepositDataRoot(validatorPubkey beacon.ValidatorPubkey, withdrawalCredentials common.Hash, validatorSignature beacon.ValidatorSignature) (common.Hash, error) {
	return ssz.HashTreeRoot(depositData{
		PublicKey:             validatorPubkey.Bytes(),
		WithdrawalCredentials: withdrawalCredentials[:],
		Amount:                depositAmount,
		Signature:             validatorSignature.Bytes(),
	})
}
