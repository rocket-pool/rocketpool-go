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
}

// Interface for all parameter types
type IParameter interface {
	GetRawValue() *big.Int
}

// Get the parameter's raw value
func (p Parameter[fType]) GetRawValue() *big.Int {
	return p.RawValue
}

// Get the formatted value of the parameter
func (p *Parameter[fType]) Formatted() fType {
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
	}
	return formatted
}

// Sets the parameter's value using the well-defined type
func (p *Parameter[fType]) Set(value fType) {
	// Switch on the parameter type and convert it
	switch f := any(&value).(type) { // Go can't switch on type parameters yet so we have to do this nonsense
	case *time.Time:
		p.RawValue = big.NewInt(f.Unix())
	case *uint64:
		p.RawValue = big.NewInt(0).SetUint64(*f)
	case *int64:
		p.RawValue = big.NewInt(*f)
	case *float64:
		p.RawValue = eth.EthToWei(*f)
	case *time.Duration:
		p.RawValue = big.NewInt(int64(f.Seconds()))
	}
}

// A parameter represented as a uint8 in the contracts, but with more useful meaning when
// properly formatted
type Uint8Parameter[fType FormattedUint8Type] struct {
	// The raw value stored in the contracts
	RawValue uint8 `json:"rawValue"`
}

// Interface for all uint8 parameter types
type IUint8Parameter interface {
	GetRawValue() uint8
}

// Get the parameter's raw value
func (p Uint8Parameter[fType]) GetRawValue() uint8 {
	return p.RawValue
}

// Get the formatted value of the parameter
func (p *Uint8Parameter[fType]) Formatted() fType {
	// Switch on the parameter type and convert it
	var formatted fType
	switch f := any(&formatted).(type) { // Go can't switch on type parameters yet so we have to do this nonsense
	case *types.MinipoolStatus:
		*f = types.MinipoolStatus(p.RawValue)
	case *types.MinipoolDeposit:
		*f = types.MinipoolDeposit(p.RawValue)
	}
	return formatted
}

// Sets the parameter's value using the well-defined type
func (p *Uint8Parameter[fType]) Set(value fType) {
	// Switch on the parameter type and convert it
	switch f := any(&value).(type) { // Go can't switch on type parameters yet so we have to do this nonsense
	case *types.MinipoolStatus:
		p.RawValue = uint8(*f)
	case *types.MinipoolDeposit:
		p.RawValue = uint8(*f)
	}
}
