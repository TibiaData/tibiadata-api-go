package validation

import "testing"

func TestTownExists(t *testing.T) {
	if !initiated {
		err := Initiate(TIBIADATA_API_TESTING)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name    string
		town    string
		want    bool
		wantErr bool
	}{
		{
			name:    "empty",
			town:    "",
			want:    false,
			wantErr: false,
		}, {
			name:    "unknown",
			town:    "anything",
			want:    false,
			wantErr: false,
		}, {
			name:    "carlin lower case",
			town:    "carlin",
			want:    true,
			wantErr: false,
		}, {
			name:    "carlin upper case",
			town:    "CARLIN",
			want:    true,
			wantErr: false,
		}, {
			name:    "carlin mixed case",
			town:    "CaRlIn",
			want:    true,
			wantErr: false,
		}, {
			name:    "ab'dendriel lower case",
			town:    "ab'dendriel",
			want:    true,
			wantErr: false,
		}, {
			name:    "ab'dendriel upper case",
			town:    "AB'DENDRIEL",
			want:    true,
			wantErr: false,
		}, {
			name:    "ab'dendriel mixed case",
			town:    "Ab'DeNdRiEl",
			want:    true,
			wantErr: false,
		}, {
			name:    "port hope lower case",
			town:    "port hope",
			want:    true,
			wantErr: false,
		}, {
			name:    "port hope upper case",
			town:    "PORT HOPE",
			want:    true,
			wantErr: false,
		}, {
			name:    "port hope mixed case",
			town:    "PoRt HoPe",
			want:    true,
			wantErr: false,
		}, {
			name:    "port hope lower case with '+'",
			town:    "port+hope",
			want:    true,
			wantErr: false,
		}, {
			name:    "port hope upper case with '+'",
			town:    "PORT+HOPE",
			want:    true,
			wantErr: false,
		}, {
			name:    "port hope mixed case with '+'",
			town:    "PoRt+HoPe",
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TownExists(tt.town)
			if (err != nil) != tt.wantErr {
				t.Errorf("TownExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TownExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
