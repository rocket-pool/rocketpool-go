package utils

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
)

type FilterQuery struct {
	BlockHash *common.Hash
	FromBlock *big.Int
	ToBlock   *big.Int
	Topics    [][]common.Hash
}

// Gets the logs for a particular log request, breaking the calls into batches if necessary
func GetLogs(rp *rocketpool.RocketPool, addressFilter []common.Address, topicFilter [][]common.Hash, intervalSize, fromBlock, toBlock *big.Int, blockHash *common.Hash) ([]types.Log, error) {
	var logs []types.Log

	// Get the block that Rocket Pool was deployed on as the lower bound if one wasn't specified
	if fromBlock == nil {
		fromBlock = rp.DeployBlock
	}

	if intervalSize == nil {
		// Handle unlimited intervals with a single call
		logs, err := rp.Client.FilterLogs(context.Background(), ethereum.FilterQuery{
			Addresses: addressFilter,
			Topics:    topicFilter,
			FromBlock: fromBlock,
			ToBlock:   toBlock,
			BlockHash: blockHash,
		})
		if err != nil {
			return nil, err
		}
		return logs, nil
	} else {
		// Get the latest block
		if toBlock == nil {
			latestBlock, err := rp.Client.BlockNumber(context.Background())
			if err != nil {
				return nil, err
			}
			toBlock = big.NewInt(0)
			toBlock.SetUint64(latestBlock)
		}

		// Set the start and end, clamping on the latest block
		intervalSize := big.NewInt(0).Sub(intervalSize, big.NewInt(1))
		start := big.NewInt(0).Set(fromBlock)
		end := big.NewInt(0).Add(start, intervalSize)
		if end.Cmp(toBlock) == 1 {
			end.Set(toBlock)
		}
		for {
			// Get the logs using the current interval
			newLogs, err := rp.Client.FilterLogs(context.Background(), ethereum.FilterQuery{
				Addresses: addressFilter,
				Topics:    topicFilter,
				FromBlock: start,
				ToBlock:   end,
				BlockHash: blockHash,
			})
			if err != nil {
				return nil, err
			}

			// Append the logs to the total list
			logs = append(logs, newLogs...)

			// Return once we've finished iterating
			if end.Cmp(toBlock) == 0 {
				return logs, nil
			}

			// Update to the next interval (end+1 : that + interval - 1)
			start.Add(end, big.NewInt(1))
			end.Add(start, intervalSize)
			if end.Cmp(toBlock) == 1 {
				end.Set(toBlock)
			}
		}
	}
}
