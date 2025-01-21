package validation

import (
	"fmt"
	"testing"
)

func TestTownExists(t *testing.T) {
	if !initiated {
		err := Initiate(TIBIADATA_API_TESTING)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Define test data with expected results
	testData := map[string]struct {
		towns   []string
		want    bool
		wantErr bool
	}{
		"empty": {
			towns:   []string{""},
			want:    false,
			wantErr: false,
		},
		"unknown": {
			towns:   []string{"anything"},
			want:    false,
			wantErr: false,
		},
		"carlin": {
			towns:   []string{"carlin", "CARLIN", "CaRlIn"},
			want:    true,
			wantErr: false,
		},
		"ab'dendriel": {
			towns:   []string{"ab'dendriel", "AB'DENDRIEL", "Ab'DeNdRiEl"},
			want:    true,
			wantErr: false,
		},
		"port hope": {
			towns:   []string{"port hope", "PORT HOPE", "PoRt HoPe"},
			want:    true,
			wantErr: false,
		},
		"port hope with +": {
			towns:   []string{"port+hope", "PORT+HOPE", "PoRt+HoPe"},
			want:    true,
			wantErr: false,
		},
	}

	// Iterate over test data and run subtests
	for name, data := range testData {
		for _, town := range data.towns {
			t.Run(fmt.Sprintf("%s (%s)", name, town), func(t *testing.T) {
				got, err := TownExists(town)
				if (err != nil) != data.wantErr {
					t.Errorf("TownExists() error = %v, wantErr %v", err, data.wantErr)
					return
				}
				if got != data.want {
					t.Errorf("TownExists() = %v, want %v", got, data.want)
				}
			})
		}
	}
}
