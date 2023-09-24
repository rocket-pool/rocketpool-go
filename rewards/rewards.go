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
	// The index of the active rewards period
	RewardIndex *core.FormattedUint256Field[uint64]

	// The timestamp that the current rewards interval started
	IntervalStart *core.FormattedUint256Field[time.Time]

	// The number of seconds in a claim interval
	IntervalDuration *core.FormattedUint256Field[time.Duration]

	// The percent of checkpoint rewards that goes to node operators
	NodeOperatorRewardsPercent *core.FormattedUint256Field[float64]

	// The percent of checkpoint rewards that goes to Ooracle DAO members
	OracleDaoRewardsPercent *core.FormattedUint256Field[float64]

	// The percent of checkpoint rewards that goes to the Protocol DAO
	ProtocolDaoRewardsPercent *core.FormattedUint256Field[float64]

	// The amount of RPL rewards that are currently pending distribution
	PendingRplRewards *core.SimpleField[*big.Int]

	// The amount of ETH rewards that are currently pending distribution
	PendingEthRewards *core.SimpleField[*big.Int]

	// === Internal fields ===
	rp          *rocketpool.RocketPool
	rewardsPool *core.Contract
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
		RewardIndex:                core.NewFormattedUint256Field[uint64](rewardsPool, "getRewardIndex"),
		IntervalStart:              core.NewFormattedUint256Field[time.Time](rewardsPool, "getClaimIntervalTimeStart"),
		IntervalDuration:           core.NewFormattedUint256Field[time.Duration](rewardsPool, "getClaimIntervalTime"),
		NodeOperatorRewardsPercent: core.NewFormattedUint256Field[float64](rewardsPool, "getClaimingContractPerc", "rocketClaimNode"),
		OracleDaoRewardsPercent:    core.NewFormattedUint256Field[float64](rewardsPool, "getClaimingContractPerc", "rocketClaimTrustedNode"),
		ProtocolDaoRewardsPercent:  core.NewFormattedUint256Field[float64](rewardsPool, "getClaimingContractPerc", "rocketClaimDAO"),
		PendingRplRewards:          core.NewSimpleField[*big.Int](rewardsPool, "getPendingRPLRewards"),
		PendingEthRewards:          core.NewSimpleField[*big.Int](rewardsPool, "getPendingETHRewards"),

		rp:          rp,
		rewardsPool: rewardsPool,
	}, nil
}

// =============
// === Calls ===
// =============

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
