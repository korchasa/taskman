FROM korchasa/go-build:latest as build
WORKDIR /app
ADD . .

#ENV GOFLAGS="-mod=vendor"
RUN \
    echo "  ## Install deps" && \
    go install && \
    echo "  ## Test" && \
    go test -cover -v -count=1 && \
    echo "  ## Lint" && \
    golangci-lint run ./... && \
    echo "  ## Build" && \
    gitget all && \
    go build -o app -ldflags "-s -w -X 'main.version=$(gitget version)' -X 'main.commit=$(gitget commit)' -X 'main.date=$(gitget date)'" . && \
    echo "  ## Done"