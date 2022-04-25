#!/bin/sh
psql -U postgres -d homework -c "\COPY cpu_usage FROM /docker-entrypoint-initdb.d/02_cpu_usage.csv CSV HEADER"