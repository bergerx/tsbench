package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"time"
)

type QueryParams struct {
	Hostname   string
	Start, End time.Time
}

type queryLoader struct {
	inputReader io.Reader
	outputChan  chan QueryParams
}

func NewCSVQueryLoader(inputReader io.Reader, chanSize int) queryLoader {
	q := queryLoader{
		outputChan:  make(chan QueryParams, chanSize),
		inputReader: inputReader,
	}
	go q.load()
	return q
}

// This channel will be closed when all items in the query are read.
func (q *queryLoader) OutputChannel() <-chan QueryParams {
	return q.outputChan
}

// reads from the input as CSV and populates output and error chans.
// expects to find hostname, start time, and end time fields in the provided CSV
func (q *queryLoader) load() {
	cr := csv.NewReader(q.inputReader)
	for {
		log.Println("reading a csv record")
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(stderr, "failed parsing file: %s\n", err)
			break
		}
		hostname, start, end, err := q.parseLine(record)
		if err != nil {
			fmt.Fprintln(stderr, err.Error())
			continue
		}
		q.outputChan <- QueryParams{Hostname: hostname, Start: start, End: end}
	}
	log.Println("finished reading CSV, closing loader output channel...")
	close(q.outputChan)
}

func (q *queryLoader) parseLine(record []string) (hostname string, start time.Time, end time.Time, err error) {
	if len(record) != 3 {
		err = fmt.Errorf("CSV line doesn't have 3 fields: %s", record)
		return
	}
	hostname = record[0]
	start, err = time.Parse("2006-01-02 15:04:05", record[1])
	if err != nil {
		err = fmt.Errorf("failed parsing start time: %s", err)
		return
	}
	end, err = time.Parse("2006-01-02 15:04:05", record[2])
	if err != nil {
		err = fmt.Errorf("failed parsing end time: %s", err)
		return
	}
	return
}
