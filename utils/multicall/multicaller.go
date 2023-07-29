/*
* This code was derived from https://github.com/depocket/multicall-go
 */

package multicall

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"golang.org/x/sync/errgroup"
)

type Call struct {
	Method   string         `json:"method"`
	Target   common.Address `json:"target"`
	CallData []byte         `json:"call_data"`
	Contract *core.Contract
	output   interface{}
}

type CallResponse struct {
	Method        string
	Status        bool
	ReturnDataRaw []byte `json:"returnData"`
}

type Result struct {
	Success bool `json:"success"`
	Output  interface{}
}

func (call Call) GetMultiCall() MultiCall {
	return MultiCall{Target: call.Target, CallData: call.CallData}
}

type MultiCaller struct {
	Client          core.ExecutionClient
	ABI             abi.ABI
	ContractAddress common.Address
	calls           []Call
}

func NewMultiCaller(client core.ExecutionClient, multicallerAddress common.Address) (*MultiCaller, error) {
	mcAbi, err := abi.JSON(strings.NewReader(MulticallABI))
	if err != nil {
		return nil, err
	}

	return &MultiCaller{
		Client:          client,
		ABI:             mcAbi,
		ContractAddress: multicallerAddress,
		calls:           []Call{},
	}, nil
}

func AddCall[outType core.CallReturnType](mc *MultiCaller, contract *core.Contract, output *outType, method string, args ...interface{}) error {
	callData, err := contract.ABI.Pack(method, args...)
	if err != nil {
		return fmt.Errorf("error adding call [%s]: %w", method, err)
	}
	call := Call{
		Method:   method,
		Target:   *contract.Address,
		CallData: callData,
		Contract: contract,
		output:   output,
	}
	mc.calls = append(mc.calls, call)
	return nil
}

func (caller *MultiCaller) Execute(requireSuccess bool, opts *bind.CallOpts) ([]CallResponse, error) {
	var multiCalls = make([]MultiCall, 0, len(caller.calls))
	for _, call := range caller.calls {
		multiCalls = append(multiCalls, call.GetMultiCall())
	}
	callData, err := caller.ABI.Pack("tryAggregate", requireSuccess, multiCalls)
	if err != nil {
		return nil, err
	}

	resp, err := caller.Client.CallContract(context.Background(), ethereum.CallMsg{To: &caller.ContractAddress, Data: callData}, opts.BlockNumber)
	if err != nil {
		return nil, err
	}

	responses, err := caller.ABI.Unpack("tryAggregate", resp)

	if err != nil {
		return nil, err
	}

	results := make([]CallResponse, len(caller.calls))
	for i, response := range responses[0].([]struct {
		Success    bool   `json:"success"`
		ReturnData []byte `json:"returnData"`
	}) {
		results[i].Method = caller.calls[i].Method
		results[i].ReturnDataRaw = response.ReturnData
		results[i].Status = response.Success
	}
	return results, nil
}

func (caller *MultiCaller) FlexibleCall(requireSuccess bool, opts *bind.CallOpts) ([]Result, error) {
	res := make([]Result, len(caller.calls))
	results, err := caller.Execute(requireSuccess, opts)
	if err != nil {
		caller.calls = []Call{}
		return nil, err
	}
	for i, call := range caller.calls {
		callSuccess := results[i].Status
		if callSuccess {
			err := call.Contract.ABI.UnpackIntoInterface(call.output, call.Method, results[i].ReturnDataRaw)
			if err != nil {
				caller.calls = []Call{}
				return nil, err
			}
		}
		res[i].Success = callSuccess
		res[i].Output = call.output
	}
	caller.calls = []Call{}
	return res, err
}

// Run a single query using multicall
func MulticallQuery[ObjType any](client core.ExecutionClient, multicallAddress common.Address, queryAdder func(*MultiCaller) (*ObjType, error), postprocess func(*ObjType) error, opts *bind.CallOpts) (*ObjType, error) {
	// The query object
	var obj *ObjType

	// Create the multicaller
	mc, err := NewMultiCaller(client, multicallAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating multicaller: %w", err)
	}

	// Run the query adder
	if queryAdder != nil {
		obj, err = queryAdder(mc)
		if err != nil {
			return nil, fmt.Errorf("error running query adder: %w", err)
		}
	}

	// Execute the multicall
	_, err = mc.FlexibleCall(true, opts)
	if err != nil {
		return nil, fmt.Errorf("error executing multicall: %w", err)
	}

	// Postprocess
	if postprocess != nil {
		err = postprocess(obj)
		if err != nil {
			return nil, fmt.Errorf("error executing postprocessor: %w", err)
		}
	}

	return obj, nil
}

// Run a single query using multicall
func MulticallQuery2(client core.ExecutionClient, multicallAddress common.Address, queryAdder func(*MultiCaller), opts *bind.CallOpts) error {
	// Create the multicaller
	mc, err := NewMultiCaller(client, multicallAddress)
	if err != nil {
		return fmt.Errorf("error creating multicaller: %w", err)
	}

	// Run the query adder
	if queryAdder != nil {
		queryAdder(mc)
	}

	// Execute the multicall
	_, err = mc.FlexibleCall(true, opts)
	if err != nil {
		return fmt.Errorf("error executing multicall: %w", err)
	}

	return nil
}

// Run a batch query using multicall
func MulticallBatchQuery[ObjType any](client core.ExecutionClient, multicallAddress common.Address, count uint64, batchSize uint64, queryAdder func([]*ObjType, uint64, *MultiCaller) error, postprocess func(*ObjType) error, opts *bind.CallOpts) ([]*ObjType, error) {
	// Create the array of query objects
	objs := make([]*ObjType, count)

	// Sync
	var wg errgroup.Group
	wg.SetLimit(int(batchSize))

	// Run getters in batches
	for i := uint64(0); i < count; i += batchSize {
		i := i
		max := i + batchSize
		if max > count {
			max = count
		}

		// Load details
		wg.Go(func() error {
			mc, err := NewMultiCaller(client, multicallAddress)
			if err != nil {
				return err
			}
			for j := i; j < max; j++ {
				if queryAdder != nil {
					err = queryAdder(objs, j, mc)
					if err != nil {
						return fmt.Errorf("error running query adder: %w", err)
					}
				}
			}
			_, err = mc.FlexibleCall(true, opts)
			if err != nil {
				return fmt.Errorf("error executing multicall: %w", err)
			}
			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return nil, fmt.Errorf("error during multicall query: %w", err)
	}

	// Do some postprocessing
	for i := range objs {
		obj := objs[i]
		if postprocess != nil {
			err := postprocess(obj)
			if err != nil {
				return nil, fmt.Errorf("error executing postprocessor: %w", err)
			}
		}
	}

	// Return
	return objs, nil
}
