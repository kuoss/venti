package model

import (
	"testing"
	"time"

	commonmodel "github.com/prometheus/common/model"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestMarshalRule(t *testing.T) {
	tests := []struct {
		input Rule
		want  string
	}{
		{
			Rule{
				Alert: "alert1",
				Expr:  "vector(1)",
				For:   commonmodel.Duration(3 * time.Hour),
			},
			"alert: alert1\nexpr: vector(1)\nfor: 3h\n",
		},
		{
			Rule{
				Alert: "alert1",
				Expr:  "vector(1)",
				For:   commonmodel.Duration(7 * 24 * time.Hour),
			},
			"alert: alert1\nexpr: vector(1)\nfor: 1w\n",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := yaml.Marshal(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.want, string(got))
		})
	}
}

func TestUnmarshalRule(t *testing.T) {
	tests := []struct {
		input string
		want  Rule
	}{
		{
			"alert: alert1\nexpr: vector(1)\nfor: 3h\n",
			Rule{
				Alert: "alert1",
				Expr:  "vector(1)",
				For:   commonmodel.Duration(3 * time.Hour),
			},
		},
		{
			"alert: alert1\nexpr: vector(1)\nfor: 1w\n",
			Rule{
				Alert: "alert1",
				Expr:  "vector(1)",
				For:   commonmodel.Duration(7 * 24 * time.Hour),
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			var got Rule
			err := yaml.Unmarshal([]byte(tt.input), &got)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
