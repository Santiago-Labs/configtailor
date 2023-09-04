package configtailor

import (
	"testing"
)

func TestTrimPathPrefix(t *testing.T) {
	tests := []struct {
		name   string
		p      string
		prefix string
		want   string
	}{
		{
			name:   "base case",
			p:      "/home/user/go/src/configtailor",
			prefix: "/home/user/go/src",
			want:   "configtailor",
		},
		{
			name:   "get rid of just home with trailing slash",
			p:      "/home/user/go/src/configtailor",
			prefix: "/home/",
			want:   "user/go/src/configtailor",
		},
		{
			name:   "without trailing slash",
			p:      "/home/user/go/src/configtailor",
			prefix: "/home",
			want:   "user/go/src/configtailor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trimPathPrefix(tt.prefix, tt.p)
			if got != tt.want {
				t.Errorf("Mappings got (%v) wanted (%v)", got, tt.want)
			}
		})
	}
}
