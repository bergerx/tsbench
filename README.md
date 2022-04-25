# tsbench

## Usage

It's suggested to use docker-compose setup for more consistent results.

### Run it all with a single `docker-compose`

Startup the setup, this command will exit after printing a summary line:

```shell
docker-compose up --abort-on-container-exit
```

Connect to postgres, to check the populated database:

```shell
docker-compose exec -- db psql -U postgres -d homework
```

Cleanup:

```shell
docker-compose down -v
docker-compose rm -f
docker rmi tsbench_tsbench
```

### Running tsbench without docker-compose

This allows you to run/experiment with the tsbench with ease,
we can still use docker-compose to get the timescale with the test data
up and running using the docker-compose like this:


```shell
# start timescale using docker-compose, this won't start tsbench
docker-compose up db

# if you want to reach out to postgres
docker exec -u postgres -it tsbench_db_1 psql -d homework
```

Running docker to run tsbench against the timescale provisioned by the docker-compose:

```shell
docker build --tag tsbench .
docker run --rm -it \
  -v $PWD/TimescaleDB_coding_assignment-RD_eng_setup/query_params.csv:/query_params.csv \
  --network tsbench \
  tsbench \
    -query-params-path="/query_params.csv" \
    -connection-string="host=db user=postgres database=homework password=sup3r-s3cur3-p4ssw0rd"
```

Running local go to run tsbench against the timescale provisioned by the docker-compose:

```shell
go run . \
  -query-params-path="./TimescaleDB_coding_assignment-RD_eng_setup/query_params.csv" \
  -connection-string="host=localhost user=postgres database=homework password=sup3r-s3cur3-p4ssw0rd"
```

## Run tests

```shell
go test -v --cover
```

## Assumptions, known issues, possible improvements

* Errors are printed to stderr immediately, no stats are collected.
* CSV header records are not handled.
* Only the pgx postgres driver is supported, client libraries may have impact on the results.
  More client libraries may be supported, but client library versions also goes fast.
* Tested only with go 1.18, other go versions may fail. Go versions may have impact on teh results.
* The amount of hosts will impact the usage of memory (`hostnameWorkerIndexMap map[string]int`),
  this can be replaced with some consistent hostname->id resolution function to get rid of the map in the memory. 
* Tests are run only once, results may be flaky.
* There is no effort spent for warming up database, first execution will likely perform worse than the corresponding
  executions.
* Preparing stable, repeatable, controlled hardware/environment to eliminate the impact of the underlying resources
  is out of scope. No baseline numbers are covered.
* Only human-readable output is supported.
* There is no test coverage for `query_executor.go`.
* Enabling `-debug` will likely increase the measured time due to logs going to the console.
