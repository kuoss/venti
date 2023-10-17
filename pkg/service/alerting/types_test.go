package alerting

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlertStateString(t *testing.T) {
	testCases := []struct {
		alertState AlertState
		want       string
		wantPanic  string
	}{
		{StateInactive, "inactive", ""},
		{StatePending, "pending", ""},
		{StateFiring, "firing", ""},
		{-1, "", "unknown alert state: -1"},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			if tc.wantPanic == "" {
				got := tc.alertState.String()
				require.Equal(t, tc.want, got)
			} else {
				require.PanicsWithError(t, tc.wantPanic, func() {
					got := tc.alertState.String()
					require.Equal(t, "", got)
				})
			}
		})
	}
}

func TestAlertingRuleState(t *testing.T) {
	testCases := []struct {
		alerts map[uint64]*Alert
		want   AlertState
	}{
		{
			map[uint64]*Alert{},
			StateInactive,
		},
		{
			map[uint64]*Alert{0: {State: StateInactive}},
			StateInactive,
		},
		{
			map[uint64]*Alert{0: {State: StatePending}},
			StatePending,
		},
		{
			map[uint64]*Alert{0: {State: StateFiring}},
			StateFiring,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			alertingRule := AlertingRule{Active: tc.alerts}
			got := alertingRule.State()
			require.Equal(t, tc.want, got)
		})
	}
}
