package dao

import "math/big"

type IBoolSetting interface {
	GetRawValue() bool
	SetRawValue(bool)
}

type IUintSetting interface {
	GetRawValue() *big.Int
	SetRawValue(value *big.Int)
}
