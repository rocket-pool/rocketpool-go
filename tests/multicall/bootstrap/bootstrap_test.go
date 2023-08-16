package bootstrap

import (
	"fmt"
	"testing"
	"time"

	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/settings"
)

const (
	// DNT
	dnt_MinMemberCount uint64 = 3
	dnt_MemberCount    uint64 = 0

	// oDAO
	odao_ChallengeCooldown time.Duration = 168 * time.Hour // 7 days

	// pDAO
)

func TestBoostrapFunctions(t *testing.T) {
	// Revert to the baseline at the end of the test
	t.Cleanup(func() {
		err := mgr.RevertToBaseline()
		if err != nil {
			t.Fatal(fmt.Errorf("error reverting to baseline snapshot: %w", err))
		}
	})

	// Initializers
	dnt, err := trustednode.NewDaoNodeTrusted(rp)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating DNT binding: %w", err))
	}
	odao, err := settings.NewOracleDaoSettings(rp)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating odao settings binding: %w", err))
	}
	pdao, err := settings.NewProtocolDaoSettings(rp)
	if err != nil {
		t.Fatal(fmt.Errorf("error creating pdao settings binding: %w", err))
	}

	// Get all of the current settings
	err = rp.Query(func(mc *batch.MultiCaller) error {
		dnt.GetAllDetails(mc)
		odao.GetAllDetails(mc)
		pdao.GetAllDetails(mc)
		return nil
	}, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error querying all initial details: %w", err))
	}

	// Verify DNT details
	if dnt.Details.MinimumMemberCount.Formatted() != dnt_MinMemberCount {
		t.Errorf("Expected DNT MinimumMemberCount = %d but it was %d", dnt_MinMemberCount, dnt.Details.MinimumMemberCount.Formatted())
	}
	if dnt.Details.MemberCount.Formatted() != dnt_MemberCount {
		t.Errorf("Expected DNT MemberCount = %d but it was %d", dnt_MemberCount, dnt.Details.MemberCount.Formatted())
	}

	// Verify oDAO settings
	if odao.Details.Members.ChallengeCooldown.Formatted() != odao_ChallengeCooldown {
		t.Errorf("Expcted oDAO challenge cooldown = %s but it was %s", odao_ChallengeCooldown, odao.Details.Members.ChallengeCooldown.Formatted())
	}

	/*
		// Bootstrap Oracle DAO members
		dnt.BootstrapMember()

		// Bootstrap a contract upgrade
		dnt.BootstrapUpgrade()
	*/
}
