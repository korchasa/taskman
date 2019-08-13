FROM korchasa/go-build:latest as build
WORKDIR /app
ADD . .

#ENV GOFLAGS="-mod=vendor"
RUN \
    cd example && \
    mv ./main.go.txt ./main.go && \
    echo "  ## Install deps" && \
    go get && \
    echo "  ## Test" && \
    go test -cover -v -count=1 && \
    echo "  ## Lint" && \
    golangci-lint run ./... && \
    echo "  ## Build" && \
    go build && \
    ./example && \
    echo "  ## Done"