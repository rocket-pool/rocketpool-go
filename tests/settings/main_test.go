package settings_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/settings"
	"github.com/rocket-pool/rocketpool-go/tests"
)

var (
	mgr  *tests.TestManager
	rp   *rocketpool.RocketPool
	pdao *settings.ProtocolDaoSettings
	odao *settings.OracleDaoSettings
	opts *bind.TransactOpts
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
	pdao, err = settings.NewProtocolDaoSettings(rp)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating pdao settings binding: %w", err))
	}
	odao, err = settings.NewOracleDaoSettings(rp)
	if err != nil {
		log.Fatal(fmt.Errorf("error creating odao settings binding: %w", err))
	}

	// Create the default values for the pDAO / oDAO settings as a reference
	err = createDefaults(mgr)
	if err != nil {
		log.Fatal("error creating defaults: %w", err)
	}

	// Use the owner account for bootstrapping things
	opts = mgr.OwnerAccount.Transactor

	// Get all of the current settings
	err = rp.Query(func(mc *batch.MultiCaller) error {
		odao.GetAllDetails(mc)
		pdao.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("error querying all initial details: %w", err))
	}

	// Verify details
	EnsureSameDetails(log.Fatalf, &odaoDefaults, &odao.Details)
	EnsureSameDetails(log.Fatalf, &pdaoDefaults, &pdao.Details)

	// Run tests
	os.Exit(m.Run())
}
