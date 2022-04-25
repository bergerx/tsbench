package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryLoaderOutputChannel(t *testing.T) {
	tests := []struct {
		name       string
		inputFile  io.Reader
		want       []QueryParams
		wantStderr string
	}{
		{
			name:      "empty csv",
			inputFile: strings.NewReader(","),
			want:      []QueryParams{},
		}, {
			name:       "Faulty CSV line",
			inputFile:  strings.NewReader(","),
			want:       []QueryParams{},
			wantStderr: "CSV line doesn't have 3 fields",
		}, {
			name:       "Faulty start time",
			inputFile:  strings.NewReader("hostname,xxx,2017-01-02 03:57:06"),
			want:       []QueryParams{},
			wantStderr: "failed parsing start time",
		}, {
			name:       "Faulty end time",
			inputFile:  strings.NewReader("hostname,2017-01-02 03:57:06,xxx"),
			want:       []QueryParams{},
			wantStderr: "failed parsing end time",
		}, {
			name:      "one line",
			inputFile: strings.NewReader("hostname,2017-01-02 03:57:06,2017-01-02 03:57:06"),
			want: []QueryParams{
				{
					Hostname: "hostname",
					Start:    parseCSVTimeTestHelper("2017-01-02 03:57:06"),
					End:      parseCSVTimeTestHelper("2017-01-02 03:57:06"),
				},
			},
		}, {
			name:      "two lines",
			inputFile: strings.NewReader("h1,2017-01-02 03:57:06,2017-01-02 03:57:06\nh2,2017-01-02 03:57:06,2017-01-02 03:57:06"),
			want: []QueryParams{
				{
					Hostname: "h1",
					Start:    parseCSVTimeTestHelper("2017-01-02 03:57:06"),
					End:      parseCSVTimeTestHelper("2017-01-02 03:57:06"),
				}, {
					Hostname: "h2",
					Start:    parseCSVTimeTestHelper("2017-01-02 03:57:06"),
					End:      parseCSVTimeTestHelper("2017-01-02 03:57:06"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr = bytes.NewBuffer(nil)
			q := NewCSVQueryLoader(tt.inputFile, 8)
			got := queryParamsChannelToSlice(q.OutputChannel())
			assert.Equal(t, tt.want, got)
			assert.Contains(t, stderr.(*bytes.Buffer).String(), tt.wantStderr)
		})
	}
}

func queryParamsSliceToChannel(queryParams []QueryParams) <-chan QueryParams {
	c := make(chan QueryParams, 100)
	go func() {
		for _, d := range queryParams {
			c <- d
		}
		close(c)
	}()
	return c
}

func queryParamsChannelToSlice(queryParams <-chan QueryParams) []QueryParams {
	s := []QueryParams{}
	for d := range queryParams {
		s = append(s, d)
	}
	return s
}
