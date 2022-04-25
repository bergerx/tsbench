package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type workerPool struct {
	workers                []*queryWorker
	hostnameWorkerIndexMap map[string]int
	inputChan              <-chan QueryParams
	outputChan             chan time.Duration
	wg                     *sync.WaitGroup
	executor               executor
	queryCount             int
}

// Start a new worker pool in background and return pointer to the pool.
func NewWorkerPool(workerCount int, inputChan <-chan QueryParams, e executor) *workerPool {
	wp := &workerPool{
		hostnameWorkerIndexMap: make(map[string]int),
		inputChan:              inputChan,
		outputChan:             make(chan time.Duration, workerCount),
		wg:                     &sync.WaitGroup{},
		executor:               e,
	}
	for i := 0; i < workerCount; i++ {
		worker := newWorker(i, wp.outputChan, e)
		wp.workers = append(wp.workers, worker)
		wp.wg.Add(1)
		log.Printf("starting %s", worker.Name())
		go worker.Run(wp.wg)
	}
	// initialize global random generator
	rand.Seed(time.Now().Unix())
	go wp.run()
	log.Println("wp created and ready to use")
	return wp
}

func (wp *workerPool) run() {
	for queryParams := range wp.inputChan {
		wp.routeQueryParamsToWorker(queryParams)
		wp.queryCount++
	}
	for _, worker := range wp.workers {
		worker.Stop()
	}
	wp.wg.Wait()
	log.Println("closing workerPool output channel ...")
	close(wp.outputChan)
}

// This channel will be closed when input channel is closed and all workers are completed.
func (wp *workerPool) OutputChannel() <-chan time.Duration {
	return wp.outputChan
}

func (wp *workerPool) QueryCount() int {
	return wp.queryCount
}

func (wp *workerPool) routeQueryParamsToWorker(queryParams QueryParams) {
	worker := wp.hostnameToWorker(queryParams.Hostname)
	log.Printf("routing queryParams for %s to %s", queryParams.Hostname, worker.Name())
	worker.Send(queryParams)
}

func (wp *workerPool) hostnameToWorker(hostname string) *queryWorker {
	workerIndex, ok := wp.hostnameWorkerIndexMap[hostname]
	if !ok {
		workerIndex = wp.pickWorkerIndexForHost(hostname)
	}
	return wp.workers[workerIndex]
}

// We pick a random worker for the host when needed.
func (wp *workerPool) pickWorkerIndexForHost(hostname string) int {
	index := rand.Intn(len(wp.workers))
	wp.hostnameWorkerIndexMap[hostname] = index
	log.Printf("picked worker %d for hostanme %s for the rest", index, hostname)
	return index
}
