package rocketpool

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/rocket-pool/rocketpool-go/v2/rocketpool"
	"github.com/rocket-pool/rocketpool-go/v2/tests"
)

var (
	mgr *tests.TestManager
	rp  *rocketpool.RocketPool
)

func TestMain(m *testing.M) {
	var err error
	mgr, err = tests.NewTestManager()
	if err != nil {
		log.Fatal(fmt.Sprintf("error getting test manager: %s", err.Error()))
	}
	rp = mgr.RocketPool

	// Run tests
	os.Exit(m.Run())

}
