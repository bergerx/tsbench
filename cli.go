package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var stdin io.ReadCloser = os.Stdin
var stdout io.Writer = os.Stdout
var stderr io.Writer = os.Stderr
var exit func(int) = os.Exit

type config struct {
	workerCount      int
	queryParamsFile  string
	connectionString string
	debug            bool
}

// parses the comman line flags and writes to the given output Writer on parse errors
func NewConfigFromFlags(args []string, output io.Writer) (*config, error) {
	c := &config{}
	commandLine := flag.NewFlagSet(args[0], flag.ContinueOnError)
	commandLine.SetOutput(output)
	commandLine.IntVar(&c.workerCount, "workers", 4, "Number of workers/threads to use")
	commandLine.StringVar(&c.queryParamsFile, "query-params-path", "", "Query params CSV file path (required)")
	commandLine.StringVar(&c.connectionString, "connection-string", "", "Postgres connection string (required)")
	commandLine.BoolVar(&c.debug, "debug", false, "Enable debug logging.")
	err := commandLine.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	if c.queryParamsFile == "" || c.connectionString == "" {
		message := "missing required flags"
		fmt.Fprintf(stderr, "%s\n", message)
		commandLine.Usage()
		return nil, errors.New(message)
	}
	return c, nil
}

func inputReader(inputFile string) (io.ReadCloser, error) {
	if inputFile == "-" {
		return stdin, nil
	}
	queryFile, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	return queryFile, nil
}

func init() {
	log.SetOutput(ioutil.Discard)
	log.SetFlags(log.Lshortfile)
}

func main() {
	startTime := time.Now()
	c, err := NewConfigFromFlags(os.Args, stderr)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		exit(1)
	}
	if c.debug {
		log.SetOutput(stderr)
	}
	inputReader, err := inputReader(c.queryParamsFile)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		exit(1)
	}
	queryLoader := NewCSVQueryLoader(inputReader, c.workerCount)
	log.Println("queryLoader created")
	queryExecutor, err := NewQueryExecutor(c.connectionString)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		exit(1)
	}
	log.Println("executor created")
	workerPool := NewWorkerPool(c.workerCount, queryLoader.OutputChannel(), *queryExecutor)
	log.Println("workerPool created")
	results := NewResults(workerPool.OutputChannel())
	log.Println("results created")
	summary, err := results.Summary()
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		exit(1)
	}
	wallClockDuration := time.Since(startTime)
	fmt.Fprintf(stdout, "Copleted %d queries with %d workers in %s:\n", workerPool.QueryCount(), c.workerCount, wallClockDuration)
	fmt.Fprintf(stdout, "  %s\n", summary)
}
