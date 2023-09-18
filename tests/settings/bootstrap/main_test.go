package bootstrap_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/dao/oracle"
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
	opts    *bind.TransactOpts
)

func TestMain(m *testing.M) {
	// Make the test manager
	var err error
	mgr, err = tests.NewTestManager()
	if err != nil {
		log.Fatal(fmt.Sprintf("error getting test manager: %s", err.Error()))
	}
	rp = mgr.RocketPool

	// Make the pDAO / oDAO bindings
	pdaoMgr, err = protocol.NewProtocolDaoManager(rp)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating pdao manager binding: %w", err))
	}
	odaoMgr, err = oracle.NewOracleDaoManager(rp)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating odao manager binding: %w", err))
	}

	// Create the default values for the pDAO / oDAO settings as a reference
	err = tests.CreateDefaults(mgr)
	if err != nil {
		log.Fatal("error creating defaults: %w", err)
	}

	// Use the owner account for bootstrapping things
	opts = mgr.OwnerAccount.Transactor

	// Get all of the current settings
	err = rp.Query(func(mc *batch.MultiCaller) error {
		odaoMgr.Settings.GetAllDetails(mc)
		pdaoMgr.Settings.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("error querying all initial details: %w", err))
	}

	// Verify details
	settings_test.EnsureSameDetails(log.Fatalf, &tests.ODaoDefaults, odaoMgr.Settings.OracleDaoSettingsDetails)
	settings_test.EnsureSameDetails(log.Fatalf, &tests.PDaoDefaults, pdaoMgr.Settings.ProtocolDaoSettingsDetails)

	// Run tests
	os.Exit(m.Run())
}
