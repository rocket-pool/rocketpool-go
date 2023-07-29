package trustednode

import (
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/settings/protocol"
	"gonum.org/v1/gonum/mathext"
)

// ===============
// === Structs ===
// ===============

// The results of the trusted node participation calculation
type TrustedNodeParticipation struct {
	StartBlock          uint64
	UpdateFrequency     uint64
	UpdateCount         uint64
	Probability         float64
	ExpectedSubmissions float64
	ActualSubmissions   map[common.Address]float64
	Participation       map[common.Address][]bool
}

// =============
// === Utils ===
// =============

// Calculates the participation rate of every trusted node on price submission since the last block that member count changed
func CalculateTrustedNodePricesParticipation(tn *DaoNodeTrusted, intervalSize *big.Int, opts *bind.CallOpts) (*TrustedNodeParticipation, error) {
	// Get the update frequency
	updatePricesFrequency, err := protocol.GetSubmitPricesFrequency(tn.rp, opts)
	if err != nil {
		return nil, err
	}
	// Get the current block
	currentBlock, err := tn.rp.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	currentBlockNumber := currentBlock.Number.Uint64()
	// Get the block of the most recent member join (limiting to 50 intervals)
	minBlock := (currentBlockNumber/updatePricesFrequency - 50) * updatePricesFrequency
	latestMemberCountChangedBlock, err := getLatestMemberCountChangedBlock(rp, minBlock, intervalSize, opts)
	if err != nil {
		return nil, err
	}
	// Get the number of current members
	memberCount, err := trustednode.GetMemberCount(rp, nil)
	if err != nil {
		return nil, err
	}
	// Start block is the first interval after the latest join
	startBlock := (latestMemberCountChangedBlock/updatePricesFrequency + 1) * updatePricesFrequency
	// The number of members that have to submit each interval
	consensus := math.Floor(float64(memberCount)/2 + 1)
	// Check if any intervals have passed
	intervalsPassed := uint64(0)
	if currentBlockNumber > startBlock {
		// The number of intervals passed
		intervalsPassed = (currentBlockNumber-startBlock)/updatePricesFrequency + 1
	}
	// How many submissions would we expect per member given a random submission
	expected := float64(intervalsPassed) * consensus / float64(memberCount)
	// Get trusted members
	members, err := trustednode.GetMembers(rp, nil)
	if err != nil {
		return nil, err
	}
	// Construct the epoch map
	participationTable := make(map[common.Address][]bool)
	// Iterate members and sum chi-square
	submissions := make(map[common.Address]float64)
	chi := float64(0)
	for _, member := range members {
		participationTable[member.Address] = make([]bool, intervalsPassed)
		actual := 0
		if intervalsPassed > 0 {
			blocks, err := GetPricesSubmissions(rp, member.Address, startBlock, intervalSize, opts)
			if err != nil {
				return nil, err
			}
			actual = len(*blocks)
			delta := float64(actual) - expected
			chi += (delta * delta) / expected
			// Add to participation table
			for _, block := range *blocks {
				// Ignore out of step updates
				if block%updatePricesFrequency == 0 {
					index := block/updatePricesFrequency - startBlock/updatePricesFrequency
					participationTable[member.Address][index] = true
				}
			}
		}
		// Save actual submission
		submissions[member.Address] = float64(actual)
	}
	// Calculate inverse cumulative density function with members-1 DoF
	probability := float64(1)
	if intervalsPassed > 0 {
		probability = 1 - mathext.GammaIncReg(float64(len(members)-1)/2, chi/2)
	}
	// Construct return value
	participation := TrustedNodeParticipation{
		Probability:         probability,
		ExpectedSubmissions: expected,
		ActualSubmissions:   submissions,
		StartBlock:          startBlock,
		UpdateFrequency:     updatePricesFrequency,
		UpdateCount:         intervalsPassed,
		Participation:       participationTable,
	}
	return &participation, nil
}

// Calculates the participation rate of every trusted node on balance submission since the last block that member count changed
func CalculateTrustedNodeBalancesParticipation(rp *rocketpool.RocketPool, intervalSize *big.Int, opts *bind.CallOpts) (*TrustedNodeParticipation, error) {
	// Get the update frequency
	updateBalancesFrequency, err := protocol.GetSubmitBalancesFrequency(rp, opts)
	if err != nil {
		return nil, err
	}
	// Get the current block
	currentBlock, err := rp.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	currentBlockNumber := currentBlock.Number.Uint64()
	// Get the block of the most recent member join (limiting to 50 intervals)
	minBlock := (currentBlockNumber/updateBalancesFrequency - 50) * updateBalancesFrequency
	latestMemberCountChangedBlock, err := getLatestMemberCountChangedBlock(rp, minBlock, intervalSize, opts)
	if err != nil {
		return nil, err
	}
	// Get the number of current members
	memberCount, err := trustednode.GetMemberCount(rp, nil)
	if err != nil {
		return nil, err
	}
	// Start block is the first interval after the latest join
	startBlock := (latestMemberCountChangedBlock/updateBalancesFrequency + 1) * updateBalancesFrequency
	// The number of members that have to submit each interval
	consensus := math.Floor(float64(memberCount)/2 + 1)
	// Check if any intervals have passed
	intervalsPassed := uint64(0)
	if currentBlockNumber > startBlock {
		// The number of intervals passed
		intervalsPassed = (currentBlockNumber-startBlock)/updateBalancesFrequency + 1
	}
	// How many submissions would we expect per member given a random submission
	expected := float64(intervalsPassed) * consensus / float64(memberCount)
	// Get trusted members
	members, err := trustednode.GetMembers(rp, nil)
	if err != nil {
		return nil, err
	}
	// Construct the epoch map
	participationTable := make(map[common.Address][]bool)
	// Iterate members and sum chi-square
	submissions := make(map[common.Address]float64)
	chi := float64(0)
	for _, member := range members {
		participationTable[member.Address] = make([]bool, intervalsPassed)
		actual := 0
		if intervalsPassed > 0 {
			blocks, err := GetBalancesSubmissions(rp, member.Address, startBlock, intervalSize, opts)
			if err != nil {
				return nil, err
			}
			actual = len(*blocks)
			delta := float64(actual) - expected
			chi += (delta * delta) / expected
			// Add to participation table
			for _, block := range *blocks {
				// Ignore out of step updates
				if block%updateBalancesFrequency == 0 {
					index := block/updateBalancesFrequency - startBlock/updateBalancesFrequency
					participationTable[member.Address][index] = true
				}
			}
		}
		// Save actual submission
		submissions[member.Address] = float64(actual)
	}
	// Calculate inverse cumulative density function with members-1 DoF
	probability := float64(1)
	if intervalsPassed > 0 {
		probability = 1 - mathext.GammaIncReg(float64(len(members)-1)/2, chi/2)
	}
	// Construct return value
	participation := TrustedNodeParticipation{
		Probability:         probability,
		ExpectedSubmissions: expected,
		ActualSubmissions:   submissions,
		StartBlock:          startBlock,
		UpdateFrequency:     updateBalancesFrequency,
		UpdateCount:         intervalsPassed,
		Participation:       participationTable,
	}
	return &participation, nil
}

// Returns a mapping of members and whether they have submitted balances this interval or not
func GetTrustedNodeLatestBalancesParticipation(rp *rocketpool.RocketPool, intervalSize *big.Int, opts *bind.CallOpts) (map[common.Address]bool, error) {
	// Get the update frequency
	updateBalancesFrequency, err := protocol.GetSubmitBalancesFrequency(rp, opts)
	if err != nil {
		return nil, err
	}
	// Get the current block
	currentBlock, err := rp.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	currentBlockNumber := currentBlock.Number.Uint64()
	// Get trusted members
	members, err := trustednode.GetMembers(rp, nil)
	if err != nil {
		return nil, err
	}
	// Get submission within the current interval
	fromBlock := currentBlockNumber / updateBalancesFrequency * updateBalancesFrequency
	submissions, err := GetLatestBalancesSubmissions(rp, fromBlock, intervalSize, opts)
	if err != nil {
		return nil, err
	}
	// Build and return result table
	participationTable := make(map[common.Address]bool)
	for _, member := range members {
		participationTable[member.Address] = false
	}
	for _, submission := range submissions {
		participationTable[submission] = true
	}
	return participationTable, nil
}

// Returns a mapping of members and whether they have submitted prices this interval or not
func GetTrustedNodeLatestPricesParticipation(rp *rocketpool.RocketPool, intervalSize *big.Int, opts *bind.CallOpts) (map[common.Address]bool, error) {
	// Get the update frequency
	updatePricesFrequency, err := protocol.GetSubmitPricesFrequency(rp, opts)
	if err != nil {
		return nil, err
	}
	// Get the current block
	currentBlock, err := rp.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	currentBlockNumber := currentBlock.Number.Uint64()
	// Get trusted members
	members, err := trustednode.GetMembers(rp, nil)
	if err != nil {
		return nil, err
	}
	// Get submission within the current interval
	fromBlock := currentBlockNumber / updatePricesFrequency * updatePricesFrequency
	submissions, err := GetLatestPricesSubmissions(rp, fromBlock, intervalSize, opts)
	if err != nil {
		return nil, err
	}
	// Build and return result table
	participationTable := make(map[common.Address]bool)
	for _, member := range members {
		participationTable[member.Address] = false
	}
	for _, submission := range submissions {
		participationTable[submission] = true
	}
	return participationTable, nil
}
