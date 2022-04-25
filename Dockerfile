# syntax=docker/dockerfile:1

FROM golang:1.18-buster AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /app/tsbench /tsbench
USER nonroot:nonroot
ENTRYPOINT ["/tsbench"]
