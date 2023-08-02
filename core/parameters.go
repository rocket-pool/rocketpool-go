package core

import (
	"math/big"
	"time"

	"github.com/rocket-pool/rocketpool-go/types"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// A parameter represented as a uint256 (a *big.Int) in the contracts, but with more useful meaning when
// properly formatted
type Parameter[fType FormattedType] struct {
	// The raw value stored in the contracts
	RawValue *big.Int `json:"rawValue"`

	// The formatted value with a more useful type
	formattedValue *fType `json:"-"`
}

// Get the formatted value of the parameter
func (p *Parameter[fType]) Formatted() fType {
	// Return the cached value
	if p.formattedValue != nil {
		return *p.formattedValue
	}

	// Switch on the parameter type and convert it
	var formatted fType
	switch f := any(&formatted).(type) { // Go can't switch on type parameters yet so we have to do this nonsense
	case *time.Time:
		*f = time.Unix(p.RawValue.Int64(), 0)
	case *uint64:
		*f = p.RawValue.Uint64()
	case *int64:
		*f = p.RawValue.Int64()
	case *float64:
		*f = eth.WeiToEth(p.RawValue)
	case *time.Duration:
		*f = time.Duration(p.RawValue.Int64()) * time.Second
	case *types.MinipoolStatus:
		*f = types.MinipoolStatus(p.RawValue.Uint64())
	}

	// Cache and return
	p.formattedValue = &formatted
	return formatted
}

// A parameter represented as a uint8 in the contracts, but with more useful meaning when
// properly formatted
type Uint8Parameter[fType FormattedUint8Type] struct {
	// The raw value stored in the contracts
	RawValue uint8 `json:"rawValue"`

	// The formatted value with a more useful type
	formattedValue *fType `json:"-"`
}

// Get the formatted value of the parameter
func (p *Uint8Parameter[fType]) Formatted() fType {
	// Return the cached value
	if p.formattedValue != nil {
		return *p.formattedValue
	}

	// Switch on the parameter type and convert it
	var formatted fType
	switch f := any(&formatted).(type) { // Go can't switch on type parameters yet so we have to do this nonsense
	case *types.MinipoolStatus:
		*f = types.MinipoolStatus(p.RawValue)
	case *types.MinipoolDeposit:
		*f = types.MinipoolDeposit(p.RawValue)
	}

	// Cache and return
	p.formattedValue = &formatted
	return formatted
}
