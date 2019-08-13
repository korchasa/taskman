FROM korchasa/go-build:latest as build
WORKDIR /app

ADD ./*.go .
ADD ./go.mod .
RUN \
    go test -cover -v -count=1 && \
    golangci-lint run ./... && \
    go install -i github.com/korchasa/taskman

ADD ./example ./example
RUN \
    cd ./example && \
    golangci-lint run ./... && \
    go run ./main.go Hello -who=me -times=5