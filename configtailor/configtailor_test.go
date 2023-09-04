package configtailor_test

import (
	"reflect"
	"configtailor/configtailor"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMappings(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    map[string][]string
	}{
		{
			name:    "cell mapping",
			content: "cell:us1,us2,us3,us4",
			want: map[string][]string{
				"cell": {"us1", "us2", "us3", "us4"},
			},
		},
		{
			name:    "cell + region mapping",
			content: "cell:us1,us2,us3,us4:region:us-west-2,eu-west-1",
			want: map[string][]string{
				"cell":   {"us1", "us2", "us3", "us4"},
				"region": {"us-west-2", "eu-west-1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := configtailor.Mappings(tt.content)
			require.NoError(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mappings got (%v) wanted (%v)", got, tt.want)
			}
		})
	}
}
