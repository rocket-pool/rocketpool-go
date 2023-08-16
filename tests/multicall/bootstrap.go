package utils

import (
	"fmt"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/tests"
)

const (
	defaultMinMemberCount uint64 = 3
	defaultProposalCooldownTime uint64 = 0
)

func BootstrapOracleDao(mgr *tests.TestManager) error {
	rp := mgr.RocketPool

	// Get the number of nodes required for a functioning Oracle DAO
	oDao, err := trustednode.NewDaoNodeTrusted(rp)
	if err != nil {
		return fmt.Errorf("error getting oDAO binding: %w", err)
	}
	err = rp.Query(func(mc *batch.MultiCaller) error {
		oDao.GetMinimumMemberCount(mc)
		return nil
	}, nil)
	if err != nil {
		return fmt.Errorf("error getting minimum Oracle DAO node count requirement: %w", err)
	}
	minRequired := oDao.Details.MinimumMemberCount

	oDao.
}
