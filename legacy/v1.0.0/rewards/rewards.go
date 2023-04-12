package rewards

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

// The event for RPL rewards claims
type RPLTokensClaimed struct {
	ClaimingContract common.Address `json:"claimingContract"`
	ClaimingAddress  common.Address `json:"claimingAddress"`
	Amount           *big.Int       `json:"amount"`
	TimeRaw          *big.Int       `json:"time"`
	Time             time.Time      `json:"-"`
}

// Get whether a claims contract is enabled
func getEnabled(claimsContract *rocketpool.Contract, claimsName string, opts *bind.CallOpts) (bool, error) {
	enabled := new(bool)
	if err := claimsContract.Call(opts, enabled, "getEnabled"); err != nil {
		return false, fmt.Errorf("Could not get %s claims contract enabled status: %w", claimsName, err)
	}
	return *enabled, nil
}

// Get whether a claimer can make a claim
// Use to check whether a claimer is able to make claims at all
func getClaimPossible(claimsContract *rocketpool.Contract, claimsName string, claimerAddress common.Address, opts *bind.CallOpts) (bool, error) {
	claimPossible := new(bool)
	if err := claimsContract.Call(opts, claimPossible, "getClaimPossible", claimerAddress); err != nil {
		return false, fmt.Errorf("Could not get %s claim possible status for %s: %w", claimsName, claimerAddress.Hex(), err)
	}
	return *claimPossible, nil
}

// Get the percentage of rewards available to a claimer
func getClaimRewardsPerc(claimsContract *rocketpool.Contract, claimsName string, claimerAddress common.Address, opts *bind.CallOpts) (float64, error) {
	claimRewardsPerc := new(*big.Int)
	if err := claimsContract.Call(opts, claimRewardsPerc, "getClaimRewardsPerc", claimerAddress); err != nil {
		return 0, fmt.Errorf("Could not get %s claim rewards percent for %s: %w", claimsName, claimerAddress.Hex(), err)
	}
	return eth.WeiToEth(*claimRewardsPerc), nil
}

// Get the total amount of rewards available to a claimer
// Use to check whether a claimer is able to make a claim for the current interval (returns zero if unable)
func getClaimRewardsAmount(claimsContract *rocketpool.Contract, claimsName string, claimerAddress common.Address, opts *bind.CallOpts) (*big.Int, error) {
	claimRewardsAmount := new(*big.Int)
	if err := claimsContract.Call(opts, claimRewardsAmount, "getClaimRewardsAmount", claimerAddress); err != nil {
		return nil, fmt.Errorf("Could not get %s claim rewards amount for %s: %w", claimsName, claimerAddress.Hex(), err)
	}
	return *claimRewardsAmount, nil
}

// Get the time that the user registered as a claimer
func getClaimingContractUserRegisteredTime(rp *rocketpool.RocketPool, claimsContract string, claimerAddress common.Address, opts *bind.CallOpts, legacyRocketRewardsPoolAddress *common.Address) (time.Time, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, legacyRocketRewardsPoolAddress, opts)
	if err != nil {
		return time.Time{}, err
	}
	claimTime := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, claimTime, "getClaimingContractUserRegisteredTime", claimsContract, claimerAddress); err != nil {
		return time.Time{}, fmt.Errorf("Could not get claims registration time on contract %s for %s: %w", claimsContract, claimerAddress.Hex(), err)
	}
	return time.Unix((*claimTime).Int64(), 0), nil
}

// Get the total amount claimed in the current interval by the given claiming contract
func getClaimingContractTotalClaimed(rp *rocketpool.RocketPool, claimsContract string, opts *bind.CallOpts, legacyRocketRewardsPoolAddress *common.Address) (*big.Int, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, legacyRocketRewardsPoolAddress, opts)
	if err != nil {
		return nil, err
	}
	totalClaimed := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, totalClaimed, "getClaimingContractTotalClaimed", claimsContract); err != nil {
		return nil, fmt.Errorf("Could not get total claimed for %s: %w", claimsContract, err)
	}
	return *totalClaimed, nil
}

// Estimate the gas of claim
func estimateClaimGas(claimsContract *rocketpool.Contract, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	return claimsContract.GetTransactionGasInfo(opts, "claim")
}

// Claim rewards
func claim(claimsContract *rocketpool.Contract, claimsName string, opts *bind.TransactOpts) (common.Hash, error) {
	tx, err := claimsContract.Transact(opts, "claim")
	if err != nil {
		return common.Hash{}, fmt.Errorf("Could not claim %s rewards: %w", claimsName, err)
	}
	return tx.Hash(), nil
}

// Get the timestamp that the current rewards interval started
func GetClaimIntervalTimeStart(rp *rocketpool.RocketPool, opts *bind.CallOpts, legacyRocketRewardsPoolAddress *common.Address) (time.Time, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, legacyRocketRewardsPoolAddress, opts)
	if err != nil {
		return time.Time{}, err
	}
	unixTime := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, unixTime, "getClaimIntervalTimeStart"); err != nil {
		return time.Time{}, fmt.Errorf("Could not get claim interval time start: %w", err)
	}
	return time.Unix((*unixTime).Int64(), 0), nil
}

// Get the number of seconds in a claim interval
func GetClaimIntervalTime(rp *rocketpool.RocketPool, opts *bind.CallOpts, legacyRocketRewardsPoolAddress *common.Address) (time.Duration, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, legacyRocketRewardsPoolAddress, opts)
	if err != nil {
		return 0, err
	}
	unixTime := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, unixTime, "getClaimIntervalTime"); err != nil {
		return 0, fmt.Errorf("Could not get claim interval time: %w", err)
	}
	return time.Duration((*unixTime).Int64()) * time.Second, nil
}

// Get the percent of checkpoint rewards that goes to node operators
func GetNodeOperatorRewardsPercent(rp *rocketpool.RocketPool, opts *bind.CallOpts, legacyRocketRewardsPoolAddress *common.Address) (float64, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, legacyRocketRewardsPoolAddress, opts)
	if err != nil {
		return 0, err
	}
	perc := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, perc, "getClaimingContractPerc", "rocketClaimNode"); err != nil {
		return 0, fmt.Errorf("Could not get node operator rewards percent: %w", err)
	}
	return eth.WeiToEth(*perc), nil
}

// Get the percent of checkpoint rewards that goes to ODAO members
func GetTrustedNodeOperatorRewardsPercent(rp *rocketpool.RocketPool, opts *bind.CallOpts, legacyRocketRewardsPoolAddress *common.Address) (float64, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, legacyRocketRewardsPoolAddress, opts)
	if err != nil {
		return 0, err
	}
	perc := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, perc, "getClaimingContractPerc", "rocketClaimTrustedNode"); err != nil {
		return 0, fmt.Errorf("Could not get trusted node operator rewards percent: %w", err)
	}
	return eth.WeiToEth(*perc), nil
}

// Get all RPLTokensClaimed events for the provided log range
func GetRPLTokenClaimsForNode(rp *rocketpool.RocketPool, nodeAddress common.Address, startBlock *big.Int, intervalSize *big.Int, legacyRocketRewardsPoolAddress *common.Address, opts *bind.CallOpts) ([]RPLTokensClaimed, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, legacyRocketRewardsPoolAddress, opts)
	if err != nil {
		return nil, err
	}
	eventSig := rocketRewardsPool.ABI.Events["RPLTokensClaimed"]

	// Construct a filter query for relevant logs
	addressFilter := []common.Address{*legacyRocketRewardsPoolAddress}
	topicFilter := [][]common.Hash{{eventSig.ID}}
	endBlock := big.NewInt(0).Add(startBlock, intervalSize)
	endBlock.Sub(endBlock, big.NewInt(1))

	// Get the event logs
	logs, err := eth.GetLogs(rp, addressFilter, topicFilter, intervalSize, startBlock, endBlock, nil)
	if err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return []RPLTokensClaimed{}, nil
	}

	// Record the relevant events
	events := []RPLTokensClaimed{}
	for _, log := range logs {
		events, err := eventSig.Inputs.Unpack(log.Data)
		if err != nil {
			return nil, fmt.Errorf("error decoding RPLTokensClaimed data from logs for block %d: %w", log.BlockNumber, err)
		}
		for _, event := range events {
			convertedEvent := event.(RPLTokensClaimed)
			if convertedEvent.ClaimingAddress == nodeAddress {
				convertedEvent.Time = time.Unix(convertedEvent.TimeRaw.Int64(), 0)
				events = append(events, convertedEvent)
			}
		}
	}

	return events, nil
}

// Get contracts
var rocketRewardsPoolLock sync.Mutex

func getRocketRewardsPool(rp *rocketpool.RocketPool, address *common.Address, opts *bind.CallOpts) (*rocketpool.Contract, error) {
	rocketRewardsPoolLock.Lock()
	defer rocketRewardsPoolLock.Unlock()
	if address == nil {
		return rp.VersionManager.V1_0_0.GetContract("rocketRewardsPool", opts)
	} else {
		return rp.VersionManager.V1_0_0.GetContractWithAddress("rocketRewardsPool", *address)
	}
}
