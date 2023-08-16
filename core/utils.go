package core

import batch "github.com/rocket-pool/batch-query"

// This is a helper for adding calls to multicall that has strongly-typed output and can take in RP contracts
func AddCall[OutType CallReturnType](mc *batch.MultiCaller, contract *Contract, output *OutType, method string, args ...any) {
	mc.AddCall(*contract.Address, contract.ABI, output, method, args...)
}
