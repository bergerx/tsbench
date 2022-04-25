package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeExecutor struct{}

func (e fakeExecutor) Execute()                                     {}
func (e fakeExecutor) Stop()                                        {}
func (e fakeExecutor) Measure(q QueryParams) (time.Duration, error) { return q.End.Sub(q.Start), nil }

var _ executor = fakeExecutor{}

var now = time.Now()

func TestWorkerPool(t *testing.T) {
	tests := []struct {
		name        string
		workerCount int
		queryParams []QueryParams
		want        []time.Duration
	}{
		{
			name:        "empty query params",
			workerCount: 4,
		},
		{
			name:        "one query-params with many workers",
			workerCount: 4,
			queryParams: []QueryParams{{Start: now, End: now.Add(time.Minute)}},
			want:        []time.Duration{time.Minute},
		},
		{
			name:        "4 query params with 2 workers",
			workerCount: 2,
			queryParams: []QueryParams{
				{Start: now, End: now.Add(time.Minute)},
				{Start: now, End: now.Add(2 * time.Minute)},
				{Start: now, End: now.Add(3 * time.Minute)},
				{Start: now, End: now.Add(4 * time.Minute)},
			},
			want: []time.Duration{time.Minute, 2 * time.Minute, 3 * time.Minute, 4 * time.Minute},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputChan := queryParamsSliceToChannel(tt.queryParams)
			wp := NewWorkerPool(tt.workerCount, inputChan, fakeExecutor{})
			durations := durationsChannelToSlice(wp.OutputChannel())
			assert.Equal(t, tt.want, durations)
			assert.Equal(t, len(tt.want), wp.QueryCount())
		})
	}
}
