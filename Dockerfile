FROM korchasa/go-build:latest as build
WORKDIR /app

ENV GOFLAGS=""
ADD ./*.go ./
ADD ./go.mod .
RUN set -ex && \
    golangci-lint run ./... && \
    go test ./... && \
    go install -i github.com/korchasa/taskman

WORKDIR ./example
ADD ./example .
RUN set -ex && \
    go run ./main.go | sed -r "s/[[:cntrl:]]\[[0-9]{1,3}m//g" | sed -r "s/go-build[0-9]*\/b001\/exe\///g" > ./actual1 && \
    diff -s -u ./expected1 ./actual1 && \
    go run ./main.go Hello -who=me -times=5 -show | sed -r "s/[[:cntrl:]]\[[0-9]{1,3}m//g" > ./actual2 && \
    diff -s -u ./expected2 ./actual2 && \
    go run ./main.go Exec "echo hello" | sed -r "s/[[:cntrl:]]\[[0-9]{1,3}m//g" > ./actual3_1 && \
    diff -s -u ./expected3 ./actual3_1
