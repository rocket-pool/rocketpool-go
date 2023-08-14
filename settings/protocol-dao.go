package settings

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for Protocol DAO settings
type ProtocolDaoSettings struct {
	Details           ProtocolDaoSettingsDetails
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
type ProtocolDaoSettingsDetails struct {
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

// Creates a new ProtocolDaoSettings binding
func NewProtocolDaoSettings(rp *rocketpool.RocketPool) (*ProtocolDaoSettings, error) {
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

	return &ProtocolDaoSettings{
		Details:     ProtocolDaoSettingsDetails{},
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
func (c *ProtocolDaoSettings) GetCreateAuctionLotEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.IsCreateLotEnabled, "getCreateLotEnabled")
}

// Check if lot bidding is currently enabled
func (c *ProtocolDaoSettings) GetBidOnAuctionLotEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.IsBidOnLotEnabled, "getBidOnLotEnabled")
}

// Get the minimum lot size in ETH
func (c *ProtocolDaoSettings) GetAuctionLotMinimumEthValue(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotMinimumEthValue, "getLotMinimumEthValue")
}

// Get the maximum lot size in ETH
func (c *ProtocolDaoSettings) GetAuctionLotMaximumEthValue(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotMaximumEthValue, "getLotMaximumEthValue")
}

// Get the lot duration, in blocks
func (c *ProtocolDaoSettings) GetAuctionLotDuration(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotDuration.RawValue, "getLotDuration")
}

// Get the lot starting price relative to current ETH price, as a fraction
func (c *ProtocolDaoSettings) GetAuctionLotStartingPriceRatio(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotStartingPriceRatio.RawValue, "getStartingPriceRatio")
}

// Get the reserve price relative to current ETH price, as a fraction
func (c *ProtocolDaoSettings) GetAuctionLotReservePriceRatio(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.AuctionContract, &c.Details.Auction.LotReservePriceRatio.RawValue, "getReservePriceRatio")
}

// === RocketDAOProtocolSettingsDeposit ===

// Check if deposits are currently enabled
func (c *ProtocolDaoSettings) GetPoolDepositEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.DepositContract, &c.Details.Deposit.IsDepositingEnabled, "getDepositEnabled")
}

// Check if deposit assignments are currently enabled
func (c *ProtocolDaoSettings) GetAssignPoolDepositsEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.DepositContract, &c.Details.Deposit.AreDepositAssignmentsEnabled, "getAssignDepositsEnabled")
}

// Get the minimum deposit to the deposit pool
func (c *ProtocolDaoSettings) GetMinimumPoolDeposit(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.DepositContract, &c.Details.Deposit.MinimumDeposit, "getMinimumDeposit")
}

// Get the maximum size of the deposit pool
func (c *ProtocolDaoSettings) GetMaximumDepositPoolSize(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.DepositContract, &c.Details.Deposit.MaximumDepositPoolSize, "getMaximumDepositPoolSize")
}

// Get the maximum assignments per deposit transaction
func (c *ProtocolDaoSettings) GetMaximumPoolDepositAssignments(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.DepositContract, &c.Details.Deposit.MaximumAssignmentsPerDeposit.RawValue, "getMaximumDepositAssignments")
}

// === RocketDAOProtocolSettingsInflation ===

// Get the RPL inflation rate per interval
func (c *ProtocolDaoSettings) GetInflationIntervalRate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.InflationContract, &c.Details.Inflation.IntervalRate.RawValue, "getInflationIntervalRate")
}

// Get the RPL inflation start time
func (c *ProtocolDaoSettings) GetInflationIntervalStartTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.InflationContract, &c.Details.Inflation.StartTime.RawValue, "getInflationIntervalStartTime")
}

// === RocketDAOProtocolSettingsMinipool ===

// Get the minipool launch balance
func (c *ProtocolDaoSettings) GetMinipoolLaunchBalance(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.LaunchBalance, "getLaunchBalance")
}

// Get the amount required from the node for a full deposit
func (c *ProtocolDaoSettings) GetFullDepositNodeAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.FullDepositNodeAmount, "getFullDepositNodeAmount")
}

// Get the amount required from the node for a half deposit
func (c *ProtocolDaoSettings) GetHalfDepositNodeAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.HalfDepositNodeAmount, "getHalfDepositNodeAmount")
}

// Get the amount required from the node for an empty deposit
func (c *ProtocolDaoSettings) GetEmptyDepositNodeAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.EmptyDepositNodeAmount, "getEmptyDepositNodeAmount")
}

// Get the amount required from the pool stakers for a full deposit
func (c *ProtocolDaoSettings) GetFullDepositUserAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.FullDepositUserAmount, "getFullDepositUserAmount")
}

// Get the amount required from the pool stakers for a half deposit
func (c *ProtocolDaoSettings) GetHalfDepositUserAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.HalfDepositUserAmount, "getHalfDepositUserAmount")
}

// Get the amount required from the pool stakers for an empty deposit
func (c *ProtocolDaoSettings) GetEmptyDepositUserAmount(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.EmptyDepositUserAmount, "getEmptyDepositUserAmount")
}

// Check if minipool withdrawable event submissions are currently enabled
func (c *ProtocolDaoSettings) GetSubmitWithdrawableEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.IsSubmitWithdrawableEnabled, "getSubmitWithdrawableEnabled")
}

// Get the timeout period, in seconds, for prelaunch minipools to launch
func (c *ProtocolDaoSettings) GetMinipoolLaunchTimeout(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.LaunchTimeout.RawValue, "getLaunchTimeout")
}

// Check if minipool bond reductions are currently enabled
func (c *ProtocolDaoSettings) GetBondReductionEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipool.IsBondReductionEnabled, "getBondReductionEnabled")
}

// === RocketDAOProtocolSettingsNetwork ===

// Get the threshold of Oracle DAO nodes that must reach consensus on oracle data to commit it
func (c *ProtocolDaoSettings) GetOracleDaoConsensusThreshold(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.OracleDaoConsensusThreshold.RawValue, "getNodeConsensusThreshold")
}

// Check if network balance submission is enabled
func (c *ProtocolDaoSettings) GetSubmitBalancesEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.IsSubmitBalancesEnabled, "getSubmitBalancesEnabled")
}

// Get the frequency, in blocks, at which network balances should be submitted by the Oracle DAO
func (c *ProtocolDaoSettings) GetSubmitBalancesFrequency(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.SubmitBalancesFrequency.RawValue, "getSubmitBalancesFrequency")
}

// Check if network price submission is enabled
func (c *ProtocolDaoSettings) GetSubmitPricesEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.IsSubmitPricesEnabled, "getSubmitPricesEnabled")
}

// Get the frequency, in blocks, at which network prices should be submitted by the Oracle DAO
func (c *ProtocolDaoSettings) GetSubmitPricesFrequency(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.SubmitPricesFrequency.RawValue, "getSubmitPricesFrequency")
}

// Get the minimum node commission rate
func (c *ProtocolDaoSettings) GetMinimumNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.MinimumNodeFee.RawValue, "getMinimumNodeFee")
}

// Get the target node commission rate
func (c *ProtocolDaoSettings) GetTargetNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.TargetNodeFee.RawValue, "getTargetNodeFee")
}

// Get the maximum node commission rate
func (c *ProtocolDaoSettings) GetMaximumNodeFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.MaximumNodeFee.RawValue, "getMaximumNodeFee")
}

// Get the range of node demand values to base fee calculations on
func (c *ProtocolDaoSettings) GetNodeFeeDemandRange(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.NodeFeeDemandRange, "getNodeFeeDemandRange")
}

// Get the target collateralization rate for the rETH NetworkContract as a fraction
func (c *ProtocolDaoSettings) GetTargetRethCollateralRate(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NetworkContract, &c.Details.Network.TargetRethCollateralRate.RawValue, "getTargetRethCollateralRate")
}

// === RocketDAOProtocolSettingsNode ===

// Check if node registration is currently enabled
func (c *ProtocolDaoSettings) GetNodeRegistrationEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NodeContract, &c.Details.Node.IsRegistrationEnabled, "getRegistrationEnabled")
}

// Check if node deposits are currently enabled
func (c *ProtocolDaoSettings) GetNodeDepositEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NodeContract, &c.Details.Node.IsDepositingEnabled, "getDepositEnabled")
}

// Check if creating vacant minipools is currently enabled
func (c *ProtocolDaoSettings) GetVacantMinipoolsEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NodeContract, &c.Details.Node.AreVacantMinipoolsEnabled, "getVacantMinipoolsEnabled")
}

// Get the minimum RPL stake per minipool as a fraction of assigned user ETH
func (c *ProtocolDaoSettings) GetMinimumPerMinipoolStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NodeContract, &c.Details.Node.MinimumPerMinipoolStake.RawValue, "getMinimumPerMinipoolStake")
}

// Get the maximum RPL stake per minipool as a fraction of assigned user ETH
func (c *ProtocolDaoSettings) GetMaximumPerMinipoolStake(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.NodeContract, &c.Details.Node.MaximumPerMinipoolStake.RawValue, "getMaximumPerMinipoolStake")
}

// === RocketDAOProtocolSettingsRewards ===

// Get the total RPL rewards amount for all rewards recipients as a fraction
func (c *ProtocolDaoSettings) GetRewardsPercentageTotal(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.RewardsContract, &c.Details.Rewards.PercentageTotal.RawValue, "getRewardsClaimersPercTotal")
}

// Get the rewards interval time
func (c *ProtocolDaoSettings) GetRewardsIntervalTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.RewardsContract, &c.Details.Rewards.IntervalTime.RawValue, "getRewardsClaimIntervalTime")
}

// === Universal ===

// Get all basic details
func (c *ProtocolDaoSettings) GetAllDetails(mc *multicall.MultiCaller) {
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
func (c *ProtocolDaoSettings) BootstrapCreateAuctionLotEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.create.enabled", value, opts)
}

// Set the create lot enabled flag
func (c *ProtocolDaoSettings) BootstrapBidOnAuctionLotEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.bidding.enabled", value, opts)
}

// Set the minimum ETH value for lots
func (c *ProtocolDaoSettings) BootstrapAuctionLotMinimumEthValue(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.value.minimum", value, opts)
}

// Set the maximum ETH value for lots
func (c *ProtocolDaoSettings) BootstrapAuctionLotMaximumEthValue(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.value.maximum", value, opts)
}

// Set the duration value for lots, in blocks
func (c *ProtocolDaoSettings) BootstrapAuctionLotDuration(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.lot.duration", value, opts)
}

// Set the starting price ratio for lots
func (c *ProtocolDaoSettings) BootstrapAuctionLotStartingPriceRatio(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.price.start", value, opts)
}

// Set the reserve price ratio for lots
func (c *ProtocolDaoSettings) BootstrapAuctionLotReservePriceRatio(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsAuction, "auction.price.reserve", value, opts)
}

// === RocketDAOProtocolSettingsDeposit ===

// Set the deposit enabled flag
func (c *ProtocolDaoSettings) BootstrapPoolDepositEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.enabled", value, opts)
}

// Set the deposit assignments enabled flag
func (c *ProtocolDaoSettings) BootstrapAssignPoolDepositsEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.assign.enabled", value, opts)
}

// Set the minimum deposit amount
func (c *ProtocolDaoSettings) BootstrapMinimumPoolDeposit(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.minimum", value, opts)
}

// Set the maximum deposit pool size
func (c *ProtocolDaoSettings) BootstrapMaximumDepositPoolSize(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.pool.maximum", value, opts)
}

// Set the max assignments per deposit
func (c *ProtocolDaoSettings) BootstrapMaximumPoolDepositAssignments(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsDeposit, "deposit.assign.maximum", value, opts)
}

// === RocketDAOProtocolSettingsInflation ===

// Set the RPL inflation rate per interval
func (c *ProtocolDaoSettings) BootstrapInflationIntervalRate(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsInflation, "rpl.inflation.interval.rate", value, opts)
}

// Set the RPL inflation start time
func (c *ProtocolDaoSettings) BootstrapInflationIntervalStartTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsInflation, "rpl.inflation.interval.start", value, opts)
}

// === RocketDAOProtocolSettingsMinipool ===

// Set the flag for enabling minipool withdrawable event submissions
func (c *ProtocolDaoSettings) BootstrapSubmitWithdrawableEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.submit.withdrawable.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapMinipoolLaunchTimeout(value time.Duration, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.launch.timeout", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapBondReductionEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsMinipool, "minipool.bond.reduction.enabled", value, opts)
}

// === RocketDAOProtocolSettingsNetwork ===

func (c *ProtocolDaoSettings) BootstrapOracleDaoConsensusThreshold(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.consensus.threshold", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapSubmitBalancesEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.balances.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapSubmitBalancesFrequency(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.balances.frequency", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapSubmitPricesEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.prices.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapSubmitPricesFrequency(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.submit.prices.frequency", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapMinimumNodeFee(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.minimum", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapTargetNodeFee(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.target", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapMaximumNodeFee(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.maximum", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapNodeFeeDemandRange(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.node.fee.demand.range", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapTargetRethCollateralRate(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNetwork, "network.reth.collateral.target", value, opts)
}

// === RocketDAOProtocolSettingsNode ===

func (c *ProtocolDaoSettings) BootstrapNodeRegistrationEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.registration.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapNodeDepositEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.deposit.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapVacantMinipoolsEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.vacant.minipools.enabled", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapMinimumPerMinipoolStake(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.per.minipool.stake.minimum", value, opts)
}

func (c *ProtocolDaoSettings) BootstrapMaximumPerMinipoolStake(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsNode, "node.per.minipool.stake.maximum", value, opts)
}

// === RocketDAOProtocolSettingsRewards ===

func (c *ProtocolDaoSettings) BootstrapRewardsIntervalTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoProtocol, rocketpool.ContractName_RocketDAOProtocolSettingsRewards, "rpl.rewards.claim.period.time", value, opts)
}
