package dao_test

import (
	"log"
	"os"
	"testing"

	"github.com/rocket-pool/rocketpool-go/dao"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/settings"
	"github.com/rocket-pool/rocketpool-go/tests"
)

var (
	mgr  *tests.TestManager
	rp   *rocketpool.RocketPool
	pdao *settings.ProtocolDaoSettings
	odao *settings.OracleDaoSettings
	dp   *dao.DaoProposal
	dnt  *trustednode.DaoNodeTrusted
	dnta *trustednode.DaoNodeTrustedActions
	dntp *trustednode.DaoNodeTrustedProposals

	odao1 *tests.Account
	odao2 *tests.Account
	odao3 *tests.Account
)

func TestMain(m *testing.M) {
	// Make the test manager
	var err error
	mgr, err = tests.NewTestManager()
	if err != nil {
		log.Fatalf("error getting test manager: %s", err.Error())
	}
	rp = mgr.RocketPool

	// Make the pDAO / oDAO bindings
	pdao, err = settings.NewProtocolDaoSettings(rp)
	if err != nil {
		fail("error creating pdao settings binding: %s", err.Error())
	}
	odao, err = settings.NewOracleDaoSettings(rp)
	if err != nil {
		fail("error creating odao settings binding: %s", err.Error())
	}
	dp, err = dao.NewDaoProposal(rp)
	if err != nil {
		fail("error creating DP: %s", err.Error())
	}
	dnt, err = trustednode.NewDaoNodeTrusted(rp)
	if err != nil {
		fail("error creating DNT: %s", err.Error())
	}
	dnta, err = trustednode.NewDaoNodeTrustedActions(rp)
	if err != nil {
		fail("error creating DNTA: %s", err.Error())
	}
	dntp, err = trustednode.NewDaoNodeTrustedProposals(rp)
	if err != nil {
		fail("error creating DNTP: %s", err.Error())
	}

	// Initialize the network
	err = mgr.InitializeDeployment()
	if err != nil {
		fail("error initializing deployment: %s", err.Error())
	}
	odao1 = mgr.NonOwnerAccounts[0]
	odao2 = mgr.NonOwnerAccounts[1]
	odao3 = mgr.NonOwnerAccounts[2]

	// Run tests
	code := m.Run()

	// Revert to the baseline after testing is done
	cleanup()

	// Done
	os.Exit(code)
}

func fail(format string, args ...any) {
	log.Printf(format, args...)
	cleanup()
	os.Exit(1)
}

func cleanup() {
	err := mgr.RevertToBaseline()
	if err != nil {
		log.Fatalf("error reverting to baseline snapshot: %s\nPlease restart Hardhat as the state will now be corrupted for other tests", err.Error())
	}
}
