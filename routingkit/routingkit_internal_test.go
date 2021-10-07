package routingkit

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseAsMeters(t *testing.T) {
	tests := []struct {
		val         string
		expectedVal float64
		expectErr   bool
	}{
		// imperial format
		{
			val:         `13'11"`,
			expectedVal: 4.241800,
		},
		{
			val:         `13'`,
			expectedVal: 3.9624,
		},
		{
			val:       `'`,
			expectErr: true,
		},
		{
			val:       `13'bdo`,
			expectErr: true,
		},
		// decimal format
		{
			val:         `13`,
			expectedVal: 13.0,
		},
		{
			val:         `13.1`,
			expectedVal: 13.1,
		},
		{
			val:         `13 m`,
			expectedVal: 13.0,
		},
		{
			val:         `13.1 m`,
			expectedVal: 13.1,
		},
		{
			val:       `13.1 moo`,
			expectErr: true,
		},
		{
			val:       `13,1`,
			expectErr: true,
		},
		{
			val:       ``,
			expectErr: true,
		},
		{
			val:       `13.1m`,
			expectErr: true,
		},
	}
	for i, test := range tests {
		val, err := parseAsMeters(test.val)
		if err != nil && test.expectErr {
			continue
		}
		if err != nil && !test.expectErr {
			t.Errorf("[%d] did not expect a parsing error, got %v", i, err)
			continue
		} else if err == nil && test.expectErr {
			t.Errorf("[%d] expected a parsing error but got none", i)
			continue
		}
		if diff := cmp.Diff(test.expectedVal, val, floatComparer); diff != "" {
			t.Errorf("[%d]: (-want, +got):\n%s", i, diff)
		}
	}
}

func TestParseAsTonnes(t *testing.T) {
	tests := []struct {
		val         string
		expectedVal float64
		expectErr   bool
	}{
		// tonnes
		{
			val:         `3`,
			expectedVal: 3.0,
		},
		{
			val:         `3.0`,
			expectedVal: 3.0,
		},
		{
			val:         `3 t`,
			expectedVal: 3.0,
		},
		{
			val:         `3.1 t`,
			expectedVal: 3.1,
		},
		{
			val:         `3 st`,
			expectedVal: 2.72155,
		},
		{
			val:         `10000 lbs`,
			expectedVal: 4.535924,
		},
		{
			val:         `10000 kg`,
			expectedVal: 10,
		},
		{
			val:         `12 lt`,
			expectedVal: 12.1926,
		},
		{
			val:         `30 cwt`,
			expectedVal: 1.524,
		},
		{
			val:       `30cwt`,
			expectErr: true,
		},
		{
			val:       `30 dogs`,
			expectErr: true,
		},
	}
	for i, test := range tests {
		val, err := parseAsTonnes(test.val)
		if err != nil && test.expectErr {
			continue
		}
		if err != nil && !test.expectErr {
			t.Errorf("[%d] did not expect a parsing error, got %v", i, err)
			continue
		} else if err == nil && test.expectErr {
			t.Errorf("[%d] expected a parsing error but got none", i)
			continue
		}
		if diff := cmp.Diff(test.expectedVal, val, floatComparer); diff != "" {
			t.Errorf("[%d]: (-want, +got):\n%s", i, diff)
		}
	}
}

var floatComparer = cmp.Comparer(func(x, y float64) bool {
	diff := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	if math.IsNaN(diff / mean) {
		return true
	}
	return (diff / mean) < 0.00001
})
