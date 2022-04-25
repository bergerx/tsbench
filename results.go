package main

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

type results struct {
	inputChan <-chan time.Duration
	durations []time.Duration
	completed bool
}

func NewResults(inputChan <-chan time.Duration) *results {
	r := &results{
		inputChan: inputChan,
	}
	go r.load()
	return r
}

func (r *results) load() {
	for result := range r.inputChan {
		r.durations = append(r.durations, result)
	}
	sort.Sort(r)
	r.completed = true
}

// Blocks until the inputChan is closed.
func (r *results) waitForResults() {
	for {
		if r.completed {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// sort.Interface implementation
func (r *results) Len() int {
	return len(r.durations)
}

// sort.Interface implementation
func (r *results) Less(i int, j int) bool {
	return r.durations[i] < r.durations[j]
}

// sort.Interface implementation
func (r *results) Swap(i int, j int) {
	d := r.durations
	d[i], d[j] = d[j], d[i]
}

func (r *results) Min() time.Duration {
	r.waitForResults()
	return r.durations[0]
}

func (r *results) Max() time.Duration {
	r.waitForResults()
	return r.durations[r.Len()-1]
}

func (r *results) Median() time.Duration {
	r.waitForResults()
	mid := (float64(r.Len()+1) / float64(2)) - 1
	midFloor := math.Floor(mid)
	if midFloor == mid {
		return r.durations[int(mid)]
	}
	prevIndex := int(midFloor)
	prevResult := r.durations[prevIndex]
	nextResult := r.durations[prevIndex+1]
	averageOfPrevAndNextItemValues := (prevResult.Nanoseconds() + nextResult.Nanoseconds()) / 2
	return time.Duration(averageOfPrevAndNextItemValues)
}

// This fails it total of duration is more than ~290 years
func (r *results) Average() time.Duration {
	r.waitForResults()
	var totalDuration time.Duration
	for _, d := range r.durations {
		totalDuration += d
	}
	averageInNanoseconds := totalDuration.Nanoseconds() / int64(r.Len())
	return time.Duration(averageInNanoseconds)
}

func (r *results) Summary() (string, error) {
	r.waitForResults()
	if r.Len() == 0 {
		return "", errors.New("no results found to summarize")
	}
	return fmt.Sprintf("min: %s, median: %s, average: %s, max: %s", r.Min(), r.Median(), r.Average(), r.Max()), nil
}
