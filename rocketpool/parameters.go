package rocketpool

import (
	"math/big"
	"time"

	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// A parameter represented as a uint256 (a *big.Int) in the contracts, but with more useful meaning when
// properly formatted
type Parameter[FormattedType time.Time | uint64 | float64] struct {
	// The raw value stored in the contracts
	RawValue *big.Int `json:"rawValue"`

	// The formatted value with a more useful type
	formattedValue *FormattedType `json:"-"`
}

// Get the formatted value of the parameter
func (p *Parameter[FormattedType]) Formatted() FormattedType {
	// Return the cached value
	if p.formattedValue != nil {
		return *p.formattedValue
	}

	// Switch on the parameter type and convert it
	var formatted FormattedType
	switch f := any(&formatted).(type) { // Go can't switch on type parameters yet so we have to do this nonsense
	case *time.Time:
		*f = time.Unix(p.RawValue.Int64(), 0)
	case *uint64:
		*f = p.RawValue.Uint64()
	case *float64:
		*f = eth.WeiToEth(p.RawValue)
	}

	// Cache and return
	p.formattedValue = &formatted
	return formatted
}
