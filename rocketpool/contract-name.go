package rocketpool

type ContractName string

const (
	ContractName_RocketAuctionManager                  ContractName = "rocketAuctionManager"
	ContractName_RocketDAOProtocol                     ContractName = "rocketDAOProtocol"
	ContractName_RocketDAOProtocolSettingsAuction      ContractName = "rocketDAOProtocolSettingsAuction"
	ContractName_RocketDAOProtocolSettingsDeposit      ContractName = "rocketDAOProtocolSettingsDeposit"
	ContractName_RocketDAOProtocolSettingsInflation    ContractName = "rocketDAOProtocolSettingsInflation"
	ContractName_RocketDAOProtocolSettingsMinipool     ContractName = "rocketDAOProtocolSettingsMinipool"
	ContractName_RocketDAOProtocolSettingsNetwork      ContractName = "rocketDAOProtocolSettingsNetwork"
	ContractName_RocketDAOProtocolSettingsNode         ContractName = "rocketDAOProtocolSettingsNode"
	ContractName_RocketDAOProtocolSettingsRewards      ContractName = "rocketDAOProtocolSettingsRewards"
	ContractName_RocketDAONodeTrustedActions           ContractName = "rocketDAONodeTrustedActions"
	ContractName_RocketDAONodeTrusted                  ContractName = "rocketDAONodeTrusted"
	ContractName_RocketDAONodeTrustedProposals         ContractName = "rocketDAONodeTrustedProposals"
	ContractName_RocketDAONodeTrustedSettingsMembers   ContractName = "rocketDAONodeTrustedSettingsMembers"
	ContractName_RocketDAONodeTrustedSettingsMinipool  ContractName = "rocketDAONodeTrustedSettingsMinipool"
	ContractName_RocketDAONodeTrustedSettingsProposals ContractName = "rocketDAONodeTrustedSettingsProposals"
	ContractName_rocketDAONodeTrustedSettingsRewards   ContractName = "rocketDAONodeTrustedSettingsRewards"
	ContractName_RocketDepositPool                     ContractName = "rocketDepositPool"
	ContractName_RocketNetworkBalances                 ContractName = "rocketNetworkBalances"
	ContractName_RocketNetworkFees                     ContractName = "rocketNetworkFees"
	ContractName_RocketNetworkPenalties                ContractName = "rocketNetworkPenalties"
	ContractName_RocketNetworkPrices                   ContractName = "rocketNetworkPrices"
	ContractName_RocketNodeDeposit                     ContractName = "rocketNodeDeposit"
	ContractName_RocketNodeDistributorFactory          ContractName = "rocketNodeDistributorFactory"
	ContractName_RocketNodeDistributorDelegate         ContractName = "rocketNodeDistributorDelegate"
	ContractName_RocketNodeManager                     ContractName = "rocketNodeManager"
	ContractName_RocketNodeStaking                     ContractName = "rocketNodeStaking"
	ContractName_RocketMerkleDistributorMainnet        ContractName = "rocketMerkleDistributorMainnet"
	ContractName_RocketRewardsPool                     ContractName = "rocketRewardsPool"
	ContractName_RocketStorage                         ContractName = "rocketStorage"
	ContractName_RocketTokenRETH                       ContractName = "rocketTokenRETH"
	ContractName_RocketTokenRPLFixedSupply             ContractName = "rocketTokenRPLFixedSupply"
	ContractName_RocketTokenRPL                        ContractName = "rocketTokenRPL"
)

// List of all singleton contract names
var ContractNames = []ContractName{
	ContractName_RocketAuctionManager,
	ContractName_RocketDAOProtocol,
	ContractName_RocketDAOProtocolSettingsAuction,
	ContractName_RocketDAOProtocolSettingsDeposit,
	ContractName_RocketDAOProtocolSettingsInflation,
	ContractName_RocketDAOProtocolSettingsMinipool,
	ContractName_RocketDAOProtocolSettingsNetwork,
	ContractName_RocketDAOProtocolSettingsNode,
	ContractName_RocketDAOProtocolSettingsRewards,
	ContractName_RocketDAONodeTrustedActions,
	ContractName_RocketDAONodeTrusted,
	ContractName_RocketDAONodeTrustedProposals,
	ContractName_RocketDAONodeTrustedSettingsMembers,
	ContractName_RocketDAONodeTrustedSettingsMinipool,
	ContractName_RocketDAONodeTrustedSettingsProposals,
	ContractName_rocketDAONodeTrustedSettingsRewards,
	ContractName_RocketDepositPool,
	ContractName_RocketNetworkBalances,
	ContractName_RocketNetworkFees,
	ContractName_RocketNetworkPenalties,
	ContractName_RocketNetworkPrices,
	ContractName_RocketNodeDeposit,
	ContractName_RocketNodeDistributorFactory,
	ContractName_RocketNodeManager,
	ContractName_RocketNodeStaking,
	ContractName_RocketMerkleDistributorMainnet,
	ContractName_RocketRewardsPool,
	ContractName_RocketStorage,
	ContractName_RocketTokenRETH,
	ContractName_RocketTokenRPLFixedSupply,
	ContractName_RocketTokenRPL,
}

// List of all instanceable contract names
var InstanceContractNames = []ContractName{
	ContractName_RocketNodeDistributorDelegate,
}
