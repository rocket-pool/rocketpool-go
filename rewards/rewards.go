package rewards

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils"
)

const (
	rewardsSnapshotSubmittedNodeKey string = "rewards.snapshot.submitted.node.key"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketRewardsPool
type RewardsPool struct {
	*RewardsPoolDetails
	rp          *rocketpool.RocketPool
	rewardsPool *core.Contract
}

// Details for RocketRewardsPool
type RewardsPoolDetails struct {
	RewardIndex                core.Parameter[uint64]        `json:"rewardIndex"`
	IntervalStart              core.Parameter[time.Time]     `json:"intervalStart"`
	IntervalDuration           core.Parameter[time.Duration] `json:"intervalDuration"`
	NodeOperatorRewardsPercent core.Parameter[float64]       `json:"nodeOperatorRewardsPercent"`
	OracleDaoRewardsPercent    core.Parameter[float64]       `json:"oracleDaoRewardsPercent"`
	ProtocolDaoRewardsPercent  core.Parameter[float64]       `json:"protocolDaoRewardsPercent"`
	PendingRplRewards          *big.Int                      `json:"pendingRplRewards"`
	PendingEthRewards          *big.Int                      `json:"pendingEthRewards"`
}

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

// ====================
// === Constructors ===
// ====================

// Creates a new RewardsPool contract binding
func NewRewardsPool(rp *rocketpool.RocketPool) (*RewardsPool, error) {
	// Create the contract
	rewardsPool, err := rp.GetContract(rocketpool.ContractName_RocketRewardsPool)
	if err != nil {
		return nil, fmt.Errorf("error getting rewards pool contract: %w", err)
	}

	return &RewardsPool{
		RewardsPoolDetails: &RewardsPoolDetails{},
		rp:                 rp,
		rewardsPool:        rewardsPool,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the index of the active rewards period
func (c *RewardsPool) GetRewardIndex(mc *batch.MultiCaller) {
	core.AddCall(mc, c.rewardsPool, &c.RewardIndex.RawValue, "getRewardIndex")
}

// Get the timestamp that the current rewards interval started
func (c *RewardsPool) GetIntervalStart(mc *batch.MultiCaller) {
	core.AddCall(mc, c.rewardsPool, &c.IntervalStart.RawValue, "getClaimIntervalTimeStart")
}

// Get the number of seconds in a claim interval
func (c *RewardsPool) GetIntervalDuration(mc *batch.MultiCaller) {
	core.AddCall(mc, c.rewardsPool, &c.IntervalDuration.RawValue, "getClaimIntervalTime")
}

// Get the percent of checkpoint rewards that goes to node operators
func (c *RewardsPool) GetNodeOperatorRewardsPercent(mc *batch.MultiCaller) {
	core.AddCall(mc, c.rewardsPool, &c.NodeOperatorRewardsPercent.RawValue, "getClaimingContractPerc", "rocketClaimNode")
}

// Get the percent of checkpoint rewards that goes to Ooracle DAO members
func (c *RewardsPool) GetOracleDaoRewardsPercent(mc *batch.MultiCaller) {
	core.AddCall(mc, c.rewardsPool, &c.OracleDaoRewardsPercent.RawValue, "getClaimingContractPerc", "rocketClaimTrustedNode")
}

// Get the percent of checkpoint rewards that goes to the Protocol DAO
func (c *RewardsPool) GetProtocolDaoRewardsPercent(mc *batch.MultiCaller) {
	core.AddCall(mc, c.rewardsPool, &c.ProtocolDaoRewardsPercent.RawValue, "getClaimingContractPerc", "rocketClaimDAO")
}

// Get the amount of RPL rewards that are currently pending distribution
func (c *RewardsPool) GetPendingRplRewards(mc *batch.MultiCaller) {
	core.AddCall(mc, c.rewardsPool, &c.PendingRplRewards, "getPendingRPLRewards")
}

// Get the amount of ETH rewards that are currently pending distribution
func (c *RewardsPool) GetPendingEthRewards(mc *batch.MultiCaller) {
	core.AddCall(mc, c.rewardsPool, &c.PendingRplRewards, "getPendingETHRewards")
}

// Check whether or not the given address has submitted for the given rewards interval
func (c *RewardsPool) GetTrustedNodeSubmitted(mc *batch.MultiCaller, nodeAddress common.Address, rewardsIndex uint64, hasSubmitted_Out *bool, opts *bind.CallOpts) {
	indexBig := big.NewInt(0).SetUint64(rewardsIndex)
	core.AddCall(mc, c.rewardsPool, hasSubmitted_Out, "getTrustedNodeSubmitted", nodeAddress, indexBig)
}

// Check whether or not the given address has submitted specific rewards info
func (c *RewardsPool) GetTrustedNodeSubmittedSpecificRewards(mc *batch.MultiCaller, nodeAddress common.Address, submission RewardSubmission, hasSubmitted_Out *bool, opts *bind.CallOpts) error {
	// NOTE: this doesn't have a view yet so we have to construct it manually, and RLP encode it
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
		return fmt.Errorf("error encoding submission data into ABI format: %w", err)
	}

	key := crypto.Keccak256Hash(bytes)
	c.rp.Storage.GetBool(mc, hasSubmitted_Out, key)
	return nil
}

// ====================
// === Transactions ===
// ====================

// Get info for submitting a Merkle Tree-based snapshot for a rewards interval
func (c *RewardsPool) SubmitRewardSnapshot(submission RewardSubmission, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return core.NewTransactionInfo(c.rewardsPool, "submitRewardSnapshot", opts, submission)
}

// =============
// === Utils ===
// =============

// Get the event info for a rewards snapshot using the Atlas getter
func (c *RewardsPool) GetRewardsEvent(rp *rocketpool.RocketPool, index uint64, rocketRewardsPoolAddresses []common.Address, opts *bind.CallOpts) (bool, RewardsEvent, error) {
	// Get the block that the event was emitted on
	indexBig := big.NewInt(0).SetUint64(index)
	blockWrapper := new(*big.Int)
	if err := c.rewardsPool.Call(opts, blockWrapper, "getClaimIntervalExecutionBlock", indexBig); err != nil {
		return false, RewardsEvent{}, fmt.Errorf("error getting the event block for interval %d: %w", index, err)
	}
	block := *blockWrapper

	// Create the list of addresses to check
	currentAddress := *c.rewardsPool.Address
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
	indexBytes := [32]byte{}
	indexBig.FillBytes(indexBytes[:])
	addressFilter := rocketRewardsPoolAddresses
	topicFilter := [][]common.Hash{{c.rewardsPool.ABI.Events["RewardSnapshot"].ID}, {indexBytes}}

	// Get the event logs
	logs, err := utils.GetLogs(rp, addressFilter, topicFilter, big.NewInt(1), block, block, nil)
	if err != nil {
		return false, RewardsEvent{}, err
	}

	// Get the log info
	values := make(map[string]interface{})
	if len(logs) == 0 {
		return false, RewardsEvent{}, nil
	}
	err = c.rewardsPool.ABI.Events["RewardSnapshot"].Inputs.UnpackIntoMap(values, logs[0].Data)
	if err != nil {
		return false, RewardsEvent{}, err
	}

	// Get the decoded data
	submissionPrototype := RewardSubmission{}
	submissionType := reflect.TypeOf(submissionPrototype)
	submission := reflect.ValueOf(values["submission"]).Convert(submissionType).Interface().(RewardSubmission)
	eventIntervalStartTime := values["intervalStartTime"].(*big.Int)
	eventIntervalEndTime := values["intervalEndTime"].(*big.Int)
	submissionTime := values["time"].(*big.Int)
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
		MerkleRoot:        common.BytesToHash(submission.MerkleRoot[:]),
		MerkleTreeCID:     submission.MerkleTreeCID,
		IntervalStartTime: time.Unix(eventIntervalStartTime.Int64(), 0),
		IntervalEndTime:   time.Unix(eventIntervalEndTime.Int64(), 0),
		SubmissionTime:    time.Unix(submissionTime.Int64(), 0),
	}

	// Convert v1.1.0-rc1 events to modern ones
	if eventData.UserETH == nil {
		eventData.UserETH = big.NewInt(0)
	}

	return true, eventData, nil
}
