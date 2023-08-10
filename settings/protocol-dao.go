package settings

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/eth"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for Protocol DAO settings
type DaoProtocolSettings struct {
	Details           DaoProtocolSettingsDetails
	AuctionContract   *core.Contract
	DepositContract   *core.Contract
	InflationContract *core.Contract
	MinipoolContract  *core.Contract
	NetworkContract   *core.Contract
	NodeContract      *core.Contract
	RewardsContract   *core.Contract

	rp          *rocketpool.RocketPool
	daoProtocol *protocol.DaoProtocol
}

// Details for Protocol DAO settings
type DaoProtocolSettingsDetails struct {
	Auction struct {
		IsCreateLotEnabled    bool                    `json:"isCreateLotEnabled"`
		IsBidOnLotEnabled     bool                    `json:"isBidOnLotEnabled"`
		LotMinimumEthValue    *big.Int                `json:"lotMinimumEthValue"`
		LotMaximumEthValue    *big.Int                `json:"lotMaximumEthValue"`
		LotDuration           core.Parameter[uint64]  `json:"lotDuration"`
		LotStartingPriceRatio core.Parameter[float64] `json:"lotStartingPriceRatio"`
		LotReservePriceRatio  core.Parameter[float64] `json:"lotReservePriceRatio"`
	} `json:"auction"`

	Deposit struct {
		IsDepositingEnabled          bool                   `json:"isDepositingEnabled"`
		AreDepositAssignmentsEnabled bool                   `json:"areDepositAssignmentsEnabled"`
		MinimumDeposit               *big.Int               `json:"minimumDeposit"`
		MaximumDepositPoolSize       *big.Int               `json:"maximumDepositPoolSize"`
		MaximumAssignmentsPerDeposit core.Parameter[uint64] `json:"maximumAssignmentsPerDeposit"`
	} `json:"deposit"`

	Inflation struct {
		IntervalRate core.Parameter[float64]   `json:"intervalRate"`
		StartTime    core.Parameter[time.Time] `json:"startTime"`
	} `json:"inflation"`

	Minipool struct {
		LaunchBalance               *big.Int                      `json:"launchBalance"`
		FullDepositNodeAmount       *big.Int                      `json:"fullDepositNodeAmount"`
		HalfDepositNodeAmount       *big.Int                      `json:"halfDepositNodeAmount"`
		EmptyDepositNodeAmount      *big.Int                      `json:"emptyDepositNodeAmount"`
		FullDepositUserAmount       *big.Int                      `json:"fullDepositUserAmount"`
		HalfDepositUserAmount       *big.Int                      `json:"halfDepositUserAmount"`
		EmptyDepositUserAmount      *big.Int                      `json:"emptyDepositUserAmount"`
		IsSubmitWithdrawableEnabled bool                          `json:"isSubmitWithdrawableEnabled"`
		LaunchTimeout               core.Parameter[time.Duration] `json:"launchTimeout"`
		IsBondReductionEnabled      bool                          `json:"isBondReductionEnabled"`
	} `json:"minipool"`

	Network struct {
		OracleDaoConsensusThreshold core.Parameter[float64] `json:"oracleDaoConsensusThreshold"`
		IsSubmitBalancesEnabled     bool                    `json:"isSubmitBalancesEnabled"`
		SubmitBalancesFrequency     core.Parameter[uint64]  `json:"submitBalancesFrequency"`
		IsSubmitPricesEnabled       bool                    `json:"isSubmitPricesEnabled"`
		SubmitPricesFrequency       core.Parameter[uint64]  `json:"submitPricesFrequency"`
		MinimumNodeFee              core.Parameter[float64] `json:"minimumNodeFee"`
		TargetNodeFee               core.Parameter[float64] `json:"targetNodeFee"`
		MaximumNodeFee              core.Parameter[float64] `json:"maximumNodeFee"`
		NodeFeeDemandRange          *big.Int                `json:"nodeFeeDemandRange"`
		TargetRethCollateralRate    core.Parameter[float64] `json:"targetRethCollateralRate"`
	} `json:"network"`

	Node struct {
		IsRegistrationEnabled     bool                    `json:"isRegistrationEnabled"`
		IsDepositingEnabled       bool                    `json:"isDepositingEnabled"`
		AreVacantMinipoolsEnabled bool                    `json:"areVacantMinipoolsEnabled"`
		MinimumPerMinipoolStake   core.Parameter[float64] `json:"minimumPerMinipoolStake"`
		MaximumPerMinipoolStake   core.Parameter[float64] `json:"maximumPerMinipoolStake"`
	} `json:"node"`

	Rewards struct {
		PercentageTotal core.Parameter[float64]       `json:"percentageTotal"`
		IntervalTime    core.Parameter[time.Duration] `json:"intervalTime"`
	} `json:"rewards"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new DaoProtocolSettingsAuction contract binding
func NewDaoProtocolSettingsAuction(rp *rocketpool.RocketPool) (*DaoProtocolSettings, error) {
	daoProtocol, err := protocol.NewDaoProtocol(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO protocol binding: %w", err)
	}

	// Get the contracts
	contracts, err := rp.GetContracts([]rocketpool.ContractName{
		rocketpool.ContractName_RocketDAOProtocolSettingsAuction,
		rocketpool.ContractName_RocketDAOProtocolSettingsDeposit,
		rocketpool.ContractName_RocketDAOProtocolSettingsInflation,
		rocketpool.ContractName_RocketDAOProtocolSettingsMinipool,
		rocketpool.ContractName_RocketDAOProtocolSettingsNetwork,
		rocketpool.ContractName_RocketDAOProtocolSettingsNode,
		rocketpool.ContractName_RocketDAOProtocolSettingsRewards,
	}...)
	if err != nil {
		return nil, fmt.Errorf("error getting protocol DAO settings contracts: %w", err)
	}

	return &DaoProtocolSettings{
		Details:     DaoProtocolSettingsDetails{},
		rp:          rp,
		daoProtocol: daoProtocol,

		AuctionContract:   contracts[0],
		DepositContract:   contracts[1],
		InflationContract: contracts[2],
		MinipoolContract:  contracts[3],
		NetworkContract:   contracts[4],
		NodeContract:      contracts[5],
		RewardsContract:   contracts[6],
	}, nil
}

// =============
// === Calls ===
// =============

// === RocketDAOProtocolSettingsAuction ===

// Check if lot creation is currently enabled
func (c *DaoProtocolSettings) GetCreateAuctionLotEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.IsCreateLotEnabled, "getCreateLotEnabled")
}

// Check if lot bidding is currently enabled
func (c *DaoProtocolSettings) GetBidOnAuctionLotEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.IsBidOnLotEnabled, "getBidOnLotEnabled")
}

// Get the minimum lot size in ETH
func (c *DaoProtocolSettings) GetAuctionLotMinimumEthValue(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotMinimumEthValue, "getLotMinimumEthValue")
}

// Get the maximum lot size in ETH
func (c *DaoProtocolSettings) GetAuctionLotMaximumEthValue(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotMaximumEthValue, "getLotMaximumEthValue")
}

// Get the lot duration, in blocks
func (c *DaoProtocolSettings) GetAuctionLotDuration(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotDuration.RawValue, "getLotDuration")
}

// Get the lot starting price relative to current ETH price, as a fraction
func (c *DaoProtocolSettings) GetAuctionLotStartingPriceRatio(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotStartingPriceRatio.RawValue, "getStartingPriceRatio")
}

// Get the reserve price relative to current ETH price, as a fraction
func (c *DaoProtocolSettings) GetAuctionLotReservePriceRatio(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotReservePriceRatio.RawValue, "getReservePriceRatio")
}

// === RocketDAOProtocolSettingsDeposit ===

// Check if deposits are currently enabled
func (c *DaoProtocolSettings) GetPoolDepositEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.DepositContract, &c.Details.Deposit.IsDepositingEnabled, "getDepositEnabled")
}

// Check if deposit assignments are currently enabled
func (c *DaoProtocolSettings) GetAssignPoolDepositsEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.DepositContract, &c.Details.Deposit.AreDepositAssignmentsEnabled, "getAssignDepositsEnabled")
}

// Get the minimum deposit to the deposit pool
func (c *DaoProtocolSettings) GetMinimumPoolDeposit(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.DepositContract, &c.Details.Deposit.MinimumDeposit, "getMinimumDeposit")
}

// Get the maximum size of the deposit pool
func (c *DaoProtocolSettings) GetMaximumDepositPoolSize(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.DepositContract, &c.Details.Deposit.MaximumDepositPoolSize, "getMaximumDepositPoolSize")
}

// Get the maximum assignments per deposit transaction
func (c *DaoProtocolSettings) GetMaximumPoolDepositAssignments(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.DepositContract, &c.Details.Deposit.MaximumAssignmentsPerDeposit.RawValue, "getMaximumDepositAssignments")
}

// === RocketDAOProtocolSettingsInflation ===

// Get the RPL inflation rate per interval
func (c *DaoProtocolSettings) GetInflationIntervalRate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.InflationContract, &c.Details.Inflation.IntervalRate.RawValue, "getInflationIntervalRate")
}

// Get the RPL inflation start time
func (c *DaoProtocolSettings) GetInflationIntervalStartTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.InflationContract, &c.Details.Inflation.StartTime.RawValue, "getInflationIntervalStartTime")
}

// === RocketDAOProtocolSettingsMinipool ===

// Get the minipool launch balance
func (c *DaoProtocolSettings) GetMinipoolLaunchBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.LaunchBalance, "getLaunchBalance")
}

// Get the amount required from the node for a full deposit
func (c *DaoProtocolSettings) GetFullDepositNodeAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.FullDepositNodeAmount, "getFullDepositNodeAmount")
}

// Get the amount required from the node for a half deposit
func (c *DaoProtocolSettings) GetHalfDepositNodeAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.HalfDepositNodeAmount, "getHalfDepositNodeAmount")
}

// Get the amount required from the node for an empty deposit
func (c *DaoProtocolSettings) GetEmptyDepositNodeAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.EmptyDepositNodeAmount, "getEmptyDepositNodeAmount")
}

// Get the amount required from the pool stakers for a full deposit
func (c *DaoProtocolSettings) GetFullDepositUserAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.FullDepositUserAmount, "getFullDepositUserAmount")
}

// Get the amount required from the pool stakers for a half deposit
func (c *DaoProtocolSettings) GetHalfDepositUserAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.HalfDepositUserAmount, "getHalfDepositUserAmount")
}

// Get the amount required from the pool stakers for an empty deposit
func (c *DaoProtocolSettings) GetEmptyDepositUserAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.EmptyDepositUserAmount, "getEmptyDepositUserAmount")
}

// Check if minipool withdrawable event submissions are currently enabled
func (c *DaoProtocolSettings) GetSubmitWithdrawableEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.IsSubmitWithdrawableEnabled, "getSubmitWithdrawableEnabled")
}

// Get the timeout period, in seconds, for prelaunch minipools to launch
func (c *DaoProtocolSettings) GetMinipoolLaunchTimeout(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.LaunchTimeout.RawValue, "getLaunchTimeout")
}

// Check if minipool bond reductions are currently enabled
func (c *DaoProtocolSettings) GetBondReductionEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.IsBondReductionEnabled, "getBondReductionEnabled")
}

// === RocketDAOProtocolSettingsNetwork ===

// Get the threshold of Oracle DAO nodes that must reach consensus on oracle data to commit it
func (c *DaoProtocolSettings) GetOracleDaoConsensusThreshold(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.OracleDaoConsensusThreshold.RawValue, "getNodeConsensusThreshold")
}

// Check if network balance submission is enabled
func (c *DaoProtocolSettings) GetSubmitBalancesEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.IsSubmitBalancesEnabled, "getSubmitBalancesEnabled")
}

// Get the frequency, in blocks, at which network balances should be submitted by the Oracle DAO
func (c *DaoProtocolSettings) GetSubmitBalancesFrequency(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.SubmitBalancesFrequency.RawValue, "getSubmitBalancesFrequency")
}

// Check if network price submission is enabled
func (c *DaoProtocolSettings) GetSubmitPricesEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.IsSubmitPricesEnabled, "getSubmitPricesEnabled")
}

// Get the frequency, in blocks, at which network prices should be submitted by the Oracle DAO
func (c *DaoProtocolSettings) GetSubmitPricesFrequency(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.SubmitPricesFrequency.RawValue, "getSubmitPricesFrequency")
}

// Get the minimum node commission rate
func (c *DaoProtocolSettings) GetMinimumNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.MinimumNodeFee.RawValue, "getMinimumNodeFee")
}

// Get the target node commission rate
func (c *DaoProtocolSettings) GetTargetNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.TargetNodeFee.RawValue, "getTargetNodeFee")
}

// Get the maximum node commission rate
func (c *DaoProtocolSettings) GetMaximumNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.MaximumNodeFee.RawValue, "getMaximumNodeFee")
}

// Get the range of node demand values to base fee calculations on
func (c *DaoProtocolSettings) GetNodeFeeDemandRange(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.NodeFeeDemandRange, "getNodeFeeDemandRange")
}

// Get the target collateralization rate for the rETH NetworkContract as a fraction
func (c *DaoProtocolSettings) GetTargetRethCollateralRate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.TargetRethCollateralRate.RawValue, "getTargetRethCollateralRate")
}

// === RocketDAOProtocolSettingsNode ===

// Check if node registration is currently enabled
func (c *DaoProtocolSettings) GetNodeRegistrationEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NodeContract, &c.Details.Node.IsRegistrationEnabled, "getRegistrationEnabled")
}

// Check if node deposits are currently enabled
func (c *DaoProtocolSettings) GetNodeDepositEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NodeContract, &c.Details.Node.IsDepositingEnabled, "getDepositEnabled")
}

// Check if creating vacant minipools is currently enabled
func (c *DaoProtocolSettings) GetVacantMinipoolsEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NodeContract, &c.Details.Node.AreVacantMinipoolsEnabled, "getVacantMinipoolsEnabled")
}

// Get the minimum RPL stake per minipool as a fraction of assigned user ETH
func (c *DaoProtocolSettings) GetMinimumPerMinipoolStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NodeContract, &c.Details.Node.MinimumPerMinipoolStake.RawValue, "getMinimumPerMinipoolStake")
}

// Get the maximum RPL stake per minipool as a fraction of assigned user ETH
func (c *DaoProtocolSettings) GetMaximumPerMinipoolStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NodeContract, &c.Details.Node.MaximumPerMinipoolStake.RawValue, "getMaximumPerMinipoolStake")
}

// === RocketDAOProtocolSettingsRewards ===

// Get the total rewards amount for all rewards recipients as a fraction
func (c *DaoProtocolSettings) GetRewardsPercentageTotal(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.RewardsContract, &c.Details.Rewards.PercentageTotal.RawValue, "getRewardsClaimersPercTotal")
}

// Get the rewards interval time
func (c *DaoProtocolSettings) GetRewardsIntervalTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.RewardsContract, &c.Details.Rewards.IntervalTime.RawValue, "getRewardsClaimIntervalTime")
}

// === Universal ===

// Get all basic details
func (c *DaoProtocolSettings) GetAllDetails(mc *multicall.MultiCaller) {
	// Auction
	c.GetCreateAuctionLotEnabled(mc)
	c.GetBidOnAuctionLotEnabled(mc)
	c.GetAuctionLotMinimumEthValue(mc)
	c.GetAuctionLotMaximumEthValue(mc)
	c.GetAuctionLotDuration(mc)
	c.GetAuctionLotStartingPriceRatio(mc)
	c.GetAuctionLotReservePriceRatio(mc)

	// Deposit
	c.GetPoolDepositEnabled(mc)
	c.GetAssignPoolDepositsEnabled(mc)
	c.GetMinimumPoolDeposit(mc)
	c.GetMaximumDepositPoolSize(mc)
	c.GetMaximumPoolDepositAssignments(mc)

	// Inflation
	c.GetInflationIntervalRate(mc)
	c.GetInflationIntervalStartTime(mc)

	// Minipool
	c.GetMinipoolLaunchBalance(mc)
	c.GetFullDepositNodeAmount(mc)
	c.GetHalfDepositNodeAmount(mc)
	c.GetEmptyDepositNodeAmount(mc)
	c.GetFullDepositUserAmount(mc)
	c.GetHalfDepositUserAmount(mc)
	c.GetEmptyDepositUserAmount(mc)
	c.GetSubmitWithdrawableEnabled(mc)
	c.GetMinipoolLaunchTimeout(mc)
	c.GetBondReductionEnabled(mc)

	// Network
	c.GetOracleDaoConsensusThreshold(mc)
	c.GetSubmitBalancesEnabled(mc)
	c.GetSubmitBalancesFrequency(mc)
	c.GetSubmitPricesEnabled(mc)
	c.GetSubmitPricesFrequency(mc)
	c.GetMinimumNodeFee(mc)
	c.GetTargetNodeFee(mc)
	c.GetMaximumNodeFee(mc)
	c.GetNodeFeeDemandRange(mc)
	c.GetTargetRethCollateralRate(mc)

	// Node
	c.GetNodeRegistrationEnabled(mc)
	c.GetNodeDepositEnabled(mc)
	c.GetVacantMinipoolsEnabled(mc)
	c.GetMinimumPerMinipoolStake(mc)
	c.GetMaximumPerMinipoolStake(mc)

	// Rewards
	c.GetRewardsPercentageTotal(mc)
	c.GetRewardsIntervalTime(mc)
}

// ====================
// === Transactions ===
// ====================

// === RocketDAOProtocolSettingsAuction ===

// Set the create lot enabled flag
func (c *DaoProtocolSettings) BootstrapCreateLotEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.create.enabled", value, opts)
}

// Set the create lot enabled flag
func (c *DaoProtocolSettings) BootstrapBidOnLotEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.bidding.enabled", value, opts)
}

// Set the minimum ETH value for lots
func (c *DaoProtocolSettings) BootstrapLotMinimumEthValue(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.value.minimum", value, opts)
}

// Set the maximum ETH value for lots
func (c *DaoProtocolSettings) BootstrapLotMaximumEthValue(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.value.maximum", value, opts)
}

// Set the duration value for lots, in blocks
func (c *DaoProtocolSettings) BootstrapLotDuration(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.duration", big.NewInt(int64(value)), opts)
}

// Set the starting price ratio for lots
func (c *DaoProtocolSettings) BootstrapLotStartingPriceRatio(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.price.start", eth.EthToWei(value), opts)
}

// Set the reserve price ratio for lots
func (c *DaoProtocolSettings) BootstrapLotReservePriceRatio(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.price.reserve", eth.EthToWei(value), opts)
}

// === RocketDAOProtocolSettingsDeposit ===

// Set the deposit enabled flag
func (c *DaoProtocolSettings) BootstrapDepositEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.enabled", value, opts)
}

// Set the deposit assignments enabled flag
func (c *DaoProtocolSettings) BootstrapAssignDepositsEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.assign.enabled", value, opts)
}

// Set the minimum deposit amount
func (c *DaoProtocolSettings) BootstrapMinimumDeposit(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.minimum", value, opts)
}

// Set the maximum deposit pool size
func (c *DaoProtocolSettings) BootstrapMaximumDepositPoolSize(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.pool.maximum", value, opts)
}

// Set the max assignments per deposit
func (c *DaoProtocolSettings) BootstrapMaximumDepositAssignments(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.assign.maximum", big.NewInt(int64(value)), opts)
}

// === RocketDAOProtocolSettingsInflation ===

// Set the RPL inflation rate per interval
func (c *DaoProtocolSettings) BootstrapInflationIntervalRate(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsInflation, "rpl.inflation.interval.rate", eth.EthToWei(value), opts)
}

// Set the RPL inflation start time
func (c *DaoProtocolSettings) BootstrapInflationIntervalStartTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsInflation, "rpl.inflation.interval.start", big.NewInt(int64(value)), opts)
}

// === RocketDAOProtocolSettingsMinipool ===

// Set the flag for enabling minipool withdrawable event submissions
func (c *DaoProtocolSettings) BootstrapSubmitWithdrawableEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.submit.withdrawable.enabled", value, opts)
}

func (c *DaoProtocolSettings) BootstrapMinipoolLaunchTimeout(value time.Duration, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.launch.timeout", big.NewInt(int64(value.Seconds())), opts)
}

func (c *DaoProtocolSettings) BootstrapBondReductionEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.bond.reduction.enabled", value, opts)
}

// === RocketDAOProtocolSettingsNetwork ===

func (c *DaoProtocolSettings) BootstrapNodeConsensusThreshold(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.consensus.threshold", eth.EthToWei(value), opts)
}

func (c *DaoProtocolSettings) BootstrapSubmitBalancesEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.balances.enabled", value, opts)
}

func (c *DaoProtocolSettings) BootstrapSubmitBalancesFrequency(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.balances.frequency", big.NewInt(int64(value)), opts)
}

func (c *DaoProtocolSettings) BootstrapSubmitPricesEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.prices.enabled", value, opts)
}

func (c *DaoProtocolSettings) BootstrapSubmitPricesFrequency(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.prices.frequency", big.NewInt(int64(value)), opts)
}

func (c *DaoProtocolSettings) BootstrapMinimumNodeFee(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.minimum", eth.EthToWei(value), opts)
}

func (c *DaoProtocolSettings) BootstrapTargetNodeFee(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.target", eth.EthToWei(value), opts)
}

func (c *DaoProtocolSettings) BootstrapMaximumNodeFee(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.maximum", eth.EthToWei(value), opts)
}

func (c *DaoProtocolSettings) BootstrapNodeFeeDemandRange(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.demand.range", value, opts)
}

func (c *DaoProtocolSettings) BootstrapTargetRethCollateralRate(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.reth.collateral.target", eth.EthToWei(value), opts)
}

// === RocketDAOProtocolSettingsNode ===

func (c *DaoProtocolSettings) BootstrapNodeRegistrationEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.registration.enabled", value, opts)
}

func (c *DaoProtocolSettings) BootstrapNodeDepositEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.deposit.enabled", value, opts)
}

func (c *DaoProtocolSettings) BootstrapVacantMinipoolsEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapBool(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.vacant.minipools.enabled", value, opts)
}

func (c *DaoProtocolSettings) BootstrapMinimumPerMinipoolStake(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.per.minipool.stake.minimum", eth.EthToWei(value), opts)
}

func (c *DaoProtocolSettings) BootstrapMaximumPerMinipoolStake(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.per.minipool.stake.maximum", eth.EthToWei(value), opts)
}

// === RocketDAOProtocolSettingsRewards ===

func (c *DaoProtocolSettings) BootstrapRewardsIntervalTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return c.daoProtocol.BootstrapUint(rocketpool.ContractName_RocketDAOProtocolSettingsRewards, "rpl.rewards.claim.period.time", big.NewInt(int64(value)), opts)
}
