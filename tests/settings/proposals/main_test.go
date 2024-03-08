package proposals_test

import (
	"log"
	"os"
	"testing"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/node-manager-core/eth"
	"github.com/rocket-pool/rocketpool-go/dao/oracle"
	"github.com/rocket-pool/rocketpool-go/dao/proposals"
	"github.com/rocket-pool/rocketpool-go/dao/protocol"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/tests"
	settings_test "github.com/rocket-pool/rocketpool-go/tests/settings"
)

var (
	mgr     *tests.TestManager
	rp      *rocketpool.RocketPool
	pdaoMgr *protocol.ProtocolDaoManager
	odaoMgr *oracle.OracleDaoManager
	dpm     *proposals.DaoProposalManager

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
	pdaoMgr, err = protocol.NewProtocolDaoManager(rp)
	if err != nil {
		fail("error creating pdao manager binding: %s", err.Error())
	}
	odaoMgr, err = oracle.NewOracleDaoManager(rp)
	if err != nil {
		fail("error creating odao manager binding: %s", err.Error())
	}
	dpm, err = proposals.NewDaoProposalManager(rp)
	if err != nil {
		fail("error creating DPM: %s", err.Error())
	}

	// Create the default values for the pDAO / oDAO settings as a reference
	err = tests.CreateDefaults(mgr)
	if err != nil {
		fail("error creating defaults: %s", err.Error())
	}

	// Get all of the current settings
	err = rp.Query(func(mc *batch.MultiCaller) error {
		eth.QueryAllFields(odaoMgr.Settings, mc)
		eth.QueryAllFields(pdaoMgr.Settings, mc)
		return nil
	}, nil)
	if err != nil {
		fail("error querying all initial details: %s", err.Error())
	}

	// Verify details
	settings_test.EnsureSameDetails(log.Fatalf, &tests.ODaoDefaults, odaoMgr.Settings)
	settings_test.EnsureSameDetails(log.Fatalf, &tests.PDaoDefaults, pdaoMgr.Settings)

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
