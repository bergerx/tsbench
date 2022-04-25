package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type queryResult struct {
	minute   time.Time
	min, max float64
	host     string
}

type executor interface {
	Execute()
	Stop()
	Measure(q QueryParams) (time.Duration, error)
}

type queryExecutor struct {
	dbpool *pgxpool.Pool
}

var _ executor = queryExecutor{}

func NewQueryExecutor(connectionString string) (*queryExecutor, error) {
	log.Println("Connecting to postgres...")
	dbpool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	log.Println("Connected to postgres.")
	return &queryExecutor{dbpool: dbpool}, nil
}

func (qe queryExecutor) Execute() {
	var greeting string
	err := qe.dbpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(stderr, "QueryRow failed: %v\n", err)
		exit(1)
	}
}

func (qe queryExecutor) Stop() {
	log.Println("Closing pgx pool ...")
	qe.dbpool.Close()
}

func (qe queryExecutor) Measure(q QueryParams) (time.Duration, error) {
	start := time.Now()
	log.Printf("Running a query for %s ...", q.Hostname)
	rows, err := qe.dbpool.Query(context.Background(), `
		SELECT
			time_bucket('1 minute', ts) AS minute,
			MIN(usage) AS min,
			MAX(usage) AS max,
			host
		FROM cpu_usage
		WHERE
			host = $1 and
			ts >= $2 and
			ts <= $3
		GROUP BY minute, host;`, q.Hostname, q.Start, q.End)
	if err != nil {
		return 0, fmt.Errorf("query failed: %v", err)
	}
	for rows.Next() {
		r := queryResult{}
		err := rows.Scan(&r.minute, &r.min, &r.max, &r.host)
		if err != nil {
			return 0, fmt.Errorf("query iteration failed: %v", err)
		}
		log.Printf("query result: %v", r)
	}
	end := time.Now()
	return end.Sub(start), nil
}
