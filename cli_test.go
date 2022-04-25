package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func parseCSVTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}

// this is a helper method to parse times, meant to be used in the table tests
func parseCSVTimeTestHelper(s string) time.Time {
	t, _ := parseCSVTime(s)
	return t
}

func TestParseCliToConfig(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantErr    bool
		wantStderr string
		wantConfig *config
	}{
		{
			name:       "without any parameter",
			args:       []string{"./tsbench"},
			wantErr:    true,
			wantStderr: "missing required flag",
			wantConfig: nil,
		},
		{
			name:       "without query-params-file",
			args:       []string{"./tsbench", "-connection-string=connection=string goes=here"},
			wantErr:    true,
			wantStderr: "missing required flag",
			wantConfig: nil,
		},
		{
			name:       "without connection-string",
			args:       []string{"./tsbench", "-query-params-path=file"},
			wantErr:    true,
			wantStderr: "missing required flag",
			wantConfig: nil,
		},
		{
			name: "with all parameters",
			args: []string{"./tsbench",
				"-query-params-path=file",
				"-workers=10",
				"-connection-string=connection=string goes=here",
				"-debug"},
			wantErr:    false,
			wantStderr: "",
			wantConfig: &config{10, "file", "connection=string goes=here", true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr = bytes.NewBuffer(nil)
			c, err := NewConfigFromFlags(tt.args, stderr)
			assert.EqualValues(t, tt.wantErr, err != nil)
			assert.Contains(t, stderr.(*bytes.Buffer).String(), tt.wantStderr)
			assert.Equal(t, tt.wantConfig, c)
		})
	}
}
