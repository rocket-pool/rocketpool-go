package state

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/v2/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
)

// Gets a Protocol DAO proposal's details using the efficient multicall contract
func GetProtocolDaoProposalDetails(rp *rocketpool.RocketPool, contracts *NetworkContracts, proposalID uint64) (*protocol.ProtocolDaoProposal, error) {
	opts := &bind.CallOpts{
		BlockNumber: contracts.ElBlockNumber,
	}

	// Make the proposal
	prop, err := protocol.NewProtocolDaoProposal(rp, proposalID)
	if err != nil {
		return nil, err
	}

	// Get all of the parameters
	err = rp.Query(func(mc *batch.MultiCaller) error {
		eth.QueryAllFields(prop, mc)
		return nil
	}, opts)
	return prop, err
}

// Gets all Protocol DAO proposal details using the efficient multicall contract
func GetAllProtocolDaoProposalDetails(rp *rocketpool.RocketPool, contracts *NetworkContracts) ([]*protocol.ProtocolDaoProposal, error) {
	opts := &bind.CallOpts{
		BlockNumber: contracts.ElBlockNumber,
	}

	mgr, err := protocol.NewProtocolDaoManager(rp)
	if err != nil {
		return nil, err
	}
	err = rp.Query(nil, opts, mgr.ProposalCount)
	if err != nil {
		return nil, err
	}
	propCount := mgr.ProposalCount.Formatted()

	return mgr.GetProposals(propCount, true, opts)
}
