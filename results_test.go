package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResultsSummary(t *testing.T) {
	tests := []struct {
		name      string
		durations []time.Duration
		want      string
		wantErr   bool
	}{
		{
			name:      "no results",
			durations: []time.Duration{},
			want:      "",
			wantErr:   true,
		}, {
			name: "single result",
			durations: []time.Duration{
				time.Minute,
			},
			want: "min: 1m0s, median: 1m0s, average: 1m0s, max: 1m0s",
		}, {
			name: "two results",
			durations: []time.Duration{
				time.Minute,
				2 * time.Minute,
			},
			want: "min: 1m0s, median: 1m30s, average: 1m30s, max: 2m0s",
		}, {
			name: "five results",
			durations: []time.Duration{
				time.Hour,
				time.Minute,
				2 * time.Minute,
				time.Hour,
				time.Minute,
			},
			want: "min: 1m0s, median: 2m0s, average: 24m48s, max: 1h0m0s",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputChan := durationsSliceToChan(tt.durations)
			r := NewResults(inputChan)
			got, err := r.Summary()
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func durationsSliceToChan(durations []time.Duration) <-chan time.Duration {
	c := make(chan time.Duration, 100)
	go func() {
		for _, d := range durations {
			c <- d
		}
		close(c)
	}()
	return c
}

func durationsChannelToSlice(durationChan <-chan time.Duration) []time.Duration {
	var durations []time.Duration
	for d := range durationChan {
		durations = append(durations, d)
	}
	return durations
}
