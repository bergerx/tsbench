package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type queryWorker struct {
	index      int
	counter    int
	inputChan  chan QueryParams
	outputChan chan<- time.Duration
	executor   executor
}

func newWorker(index int, outputChan chan<- time.Duration, e executor) *queryWorker {
	return &queryWorker{
		index:      index,
		inputChan:  make(chan QueryParams),
		outputChan: outputChan,
		executor:   e,
	}
}

func (w *queryWorker) Name() string {
	return fmt.Sprintf("worker-%d", w.index)
}

func (w *queryWorker) Run(wg *sync.WaitGroup) {
	for queryParams := range w.inputChan {
		w.counter++
		duration, err := w.runQuery(queryParams)
		if err != nil {
			fmt.Fprintf(stderr, "error running query: %v\n", err)
			continue
		}
		w.outputChan <- duration
	}
	log.Printf("%s is done, processed %d queries ...", w.Name(), w.counter)
	wg.Done()
}

func (w *queryWorker) Stop() {
	log.Printf("closing input channel for %s after processing %d queries ...", w.Name(), w.counter)
	close(w.inputChan)
}

func (w *queryWorker) Send(queryParams QueryParams) {
	w.inputChan <- queryParams
}

func (w *queryWorker) runQuery(q QueryParams) (time.Duration, error) {
	return w.executor.Measure(q)
}
