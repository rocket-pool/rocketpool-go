package rewards

import (
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
)

const (
	rewardsSnapshotSubmittedNodeKey string = "rewards.snapshot.submitted.node.key"
)

// Info for a rewards snapshot event
type RewardsEvent struct {
	Index             *big.Int
	ExecutionBlock    *big.Int
	ConsensusBlock    *big.Int
	MerkleRoot        common.Hash
	MerkleTreeCID     string
	IntervalsPassed   *big.Int
	TreasuryRPL       *big.Int
	TrustedNodeRPL    []*big.Int
	NodeRPL           []*big.Int
	NodeETH           []*big.Int
	UserETH           *big.Int
	IntervalStartTime time.Time
	IntervalEndTime   time.Time
	SubmissionTime    time.Time
}

// Struct for submitting the rewards for a checkpoint
type RewardSubmission struct {
	RewardIndex     *big.Int   `json:"rewardIndex"`
	ExecutionBlock  *big.Int   `json:"executionBlock"`
	ConsensusBlock  *big.Int   `json:"consensusBlock"`
	MerkleRoot      [32]byte   `json:"merkleRoot"`
	MerkleTreeCID   string     `json:"merkleTreeCID"`
	IntervalsPassed *big.Int   `json:"intervalsPassed"`
	TreasuryRPL     *big.Int   `json:"treasuryRPL"`
	TrustedNodeRPL  []*big.Int `json:"trustedNodeRPL"`
	NodeRPL         []*big.Int `json:"nodeRPL"`
	NodeETH         []*big.Int `json:"nodeETH"`
	UserETH         *big.Int   `json:"userETH"`
}

// Internal struct - this is the structure of what gets returned by the RewardSnapshot event
type rewardSnapshot struct {
	RewardIndex       *big.Int         `json:"rewardIndex"`
	Submission        RewardSubmission `json:"submission"`
	IntervalStartTime *big.Int         `json:"intervalStartTime"`
	IntervalEndTime   *big.Int         `json:"intervalEndTime"`
	Time              *big.Int         `json:"time"`
}

// Get the index of the active rewards period
func GetRewardIndex(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*big.Int, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, opts)
	if err != nil {
		return nil, err
	}
	index := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, index, "getRewardIndex"); err != nil {
		return nil, fmt.Errorf("error getting current reward index: %w", err)
	}
	return *index, nil
}

// Get the timestamp that the current rewards interval started
func GetClaimIntervalTimeStart(rp *rocketpool.RocketPool, opts *bind.CallOpts) (time.Time, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, opts)
	if err != nil {
		return time.Time{}, err
	}
	unixTime := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, unixTime, "getClaimIntervalTimeStart"); err != nil {
		return time.Time{}, fmt.Errorf("error getting claim interval time start: %w", err)
	}
	return time.Unix((*unixTime).Int64(), 0), nil
}

// Get the number of seconds in a claim interval
func GetClaimIntervalTime(rp *rocketpool.RocketPool, opts *bind.CallOpts) (time.Duration, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, opts)
	if err != nil {
		return 0, err
	}
	unixTime := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, unixTime, "getClaimIntervalTime"); err != nil {
		return 0, fmt.Errorf("error getting claim interval time: %w", err)
	}
	return time.Duration((*unixTime).Int64()) * time.Second, nil
}

// Get the amount of RPL rewards that will be provided to node operators
func GetPendingRPLRewards(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*big.Int, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, opts)
	if err != nil {
		return nil, err
	}
	rewards := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, rewards, "getPendingRPLRewards"); err != nil {
		return nil, fmt.Errorf("error getting pending RPL rewards: %w", err)
	}
	return *rewards, nil
}

// Get the amount of ETH rewards that will be provided to node operators
func GetPendingETHRewards(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*big.Int, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, opts)
	if err != nil {
		return nil, err
	}
	rewards := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, rewards, "getPendingETHRewards"); err != nil {
		return nil, fmt.Errorf("error getting pending ETH rewards: %w", err)
	}
	return *rewards, nil
}

// Check whether or not the given address has submitted for the given rewards interval
func GetTrustedNodeSubmitted(rp *rocketpool.RocketPool, nodeAddress common.Address, rewardsIndex uint64, opts *bind.CallOpts) (bool, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, opts)
	if err != nil {
		return false, err
	}

	indexBig := big.NewInt(0).SetUint64(rewardsIndex)
	hasSubmitted := new(bool)
	if err := rocketRewardsPool.Call(opts, hasSubmitted, "getTrustedNodeSubmitted", nodeAddress, indexBig); err != nil {
		return false, fmt.Errorf("error getting trusted node submission status: %w", err)
	}
	return *hasSubmitted, nil
}

// Check whether or not the given address has submitted specific rewards info
func GetTrustedNodeSubmittedSpecificRewards(rp *rocketpool.RocketPool, nodeAddress common.Address, submission RewardSubmission, opts *bind.CallOpts) (bool, error) {
	// NOTE: this doesn't have a view yet so we have to construct it manually, and RLP
	stringTy, _ := abi.NewType("string", "string", nil)
	addressTy, _ := abi.NewType("address", "address", nil)

	submissionTy, _ := abi.NewType("tuple", "struct RewardSubmission", []abi.ArgumentMarshaling{
		{Name: "rewardIndex", Type: "uint256"},
		{Name: "executionBlock", Type: "uint256"},
		{Name: "consensusBlock", Type: "uint256"},
		{Name: "merkleRoot", Type: "bytes32"},
		{Name: "merkleTreeCID", Type: "string"},
		{Name: "intervalsPassed", Type: "uint256"},
		{Name: "treasuryRPL", Type: "uint256"},
		{Name: "trustedNodeRPL", Type: "uint256[]"},
		{Name: "nodeRPL", Type: "uint256[]"},
		{Name: "nodeETH", Type: "uint256[]"},
		{Name: "userETH", Type: "uint256"},
	})

	args := abi.Arguments{
		{Type: stringTy, Name: "key"},
		{Type: addressTy, Name: "trustedNodeAddress"},
		{Type: submissionTy, Name: "submission"},
	}

	bytes, err := args.Pack(rewardsSnapshotSubmittedNodeKey, nodeAddress, &submission)
	if err != nil {
		return false, fmt.Errorf("error encoding submission data into ABI format: %w", err)
	}

	key := crypto.Keccak256Hash(bytes)
	result, err := rp.RocketStorage.GetBool(opts, key)
	if err != nil {
		return false, fmt.Errorf("error checking if trusted node submitted specific rewards: %w", err)
	}
	return result, nil
}

// Estimate the gas for submiting a Merkle Tree-based snapshot for a rewards interval
func EstimateSubmitRewardSnapshotGas(rp *rocketpool.RocketPool, submission RewardSubmission, opts *bind.TransactOpts) (rocketpool.GasInfo, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, nil)
	if err != nil {
		return rocketpool.GasInfo{}, err
	}
	return rocketRewardsPool.GetTransactionGasInfo(opts, "submitRewardSnapshot", submission)
}

// Submit a Merkle Tree-based snapshot for a rewards interval
func SubmitRewardSnapshot(rp *rocketpool.RocketPool, submission RewardSubmission, opts *bind.TransactOpts) (common.Hash, error) {
	rocketRewardsPool, err := getRocketRewardsPool(rp, nil)
	if err != nil {
		return common.Hash{}, err
	}
	tx, err := rocketRewardsPool.Transact(opts, "submitRewardSnapshot", submission)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error submitting rewards snapshot: %w", err)
	}
	return tx.Hash(), nil
}

// Get the event info for a rewards snapshot using the Atlas getter
func GetRewardsEvent(rp *rocketpool.RocketPool, index uint64, rocketRewardsPoolAddresses []common.Address, opts *bind.CallOpts) (bool, RewardsEvent, error) {
	// Get contracts
	rocketRewardsPool, err := getRocketRewardsPool(rp, opts)
	if err != nil {
		return false, RewardsEvent{}, err
	}

	// Get the block that the event was emitted on
	indexBig := big.NewInt(0).SetUint64(index)
	blockWrapper := new(*big.Int)
	if err := rocketRewardsPool.Call(opts, blockWrapper, "getClaimIntervalExecutionBlock", indexBig); err != nil {
		return false, RewardsEvent{}, fmt.Errorf("error getting the event block for interval %d: %w", index, err)
	}
	block := *blockWrapper

	// Create the list of addresses to check
	currentAddress := *rocketRewardsPool.Address
	if rocketRewardsPoolAddresses == nil {
		rocketRewardsPoolAddresses = []common.Address{currentAddress}
	} else {
		found := false
		for _, address := range rocketRewardsPoolAddresses {
			if address == currentAddress {
				found = true
				break
			}
		}
		if !found {
			rocketRewardsPoolAddresses = append(rocketRewardsPoolAddresses, currentAddress)
		}
	}

	// Construct a filter query for relevant logs
	rewardsSnapshotEvent := rocketRewardsPool.ABI.Events["RewardSnapshot"]
	indexBytes := [32]byte{}
	indexBig.FillBytes(indexBytes[:])
	addressFilter := rocketRewardsPoolAddresses
	topicFilter := [][]common.Hash{{rewardsSnapshotEvent.ID}, {indexBytes}}

	// Get the event logs
	logs, err := eth.GetLogs(rp, addressFilter, topicFilter, big.NewInt(1), block, block, nil)
	if err != nil {
		return false, RewardsEvent{}, err
	}
	if len(logs) == 0 {
		return false, RewardsEvent{}, nil
	}

	// Get the log info values
	values, err := rewardsSnapshotEvent.Inputs.Unpack(logs[0].Data)
	if err != nil {
		return false, RewardsEvent{}, fmt.Errorf("error unpacking rewards snapshot event data: %w", err)
	}

	// Convert to a native struct
	var snapshot rewardSnapshot
	err = rewardsSnapshotEvent.Inputs.Copy(&snapshot, values)
	if err != nil {
		return false, RewardsEvent{}, fmt.Errorf("error converting rewards snapshot event data to struct: %w", err)
	}

	// Get the decoded data
	submission := snapshot.Submission
	eventData := RewardsEvent{
		Index:             indexBig,
		ExecutionBlock:    submission.ExecutionBlock,
		ConsensusBlock:    submission.ConsensusBlock,
		IntervalsPassed:   submission.IntervalsPassed,
		TreasuryRPL:       submission.TreasuryRPL,
		TrustedNodeRPL:    submission.TrustedNodeRPL,
		NodeRPL:           submission.NodeRPL,
		NodeETH:           submission.NodeETH,
		UserETH:           submission.UserETH,
		MerkleRoot:        submission.MerkleRoot,
		MerkleTreeCID:     submission.MerkleTreeCID,
		IntervalStartTime: time.Unix(snapshot.IntervalStartTime.Int64(), 0),
		IntervalEndTime:   time.Unix(snapshot.IntervalEndTime.Int64(), 0),
		SubmissionTime:    time.Unix(snapshot.Time.Int64(), 0),
	}

	// Convert v1.1.0-rc1 events to modern ones
	if eventData.UserETH == nil {
		eventData.UserETH = big.NewInt(0)
	}

	return true, eventData, nil
}

// Get contracts
var rocketRewardsPoolLock sync.Mutex

func getRocketRewardsPool(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*rocketpool.Contract, error) {
	rocketRewardsPoolLock.Lock()
	defer rocketRewardsPoolLock.Unlock()
	return rp.GetContract("rocketRewardsPool", opts)
}
