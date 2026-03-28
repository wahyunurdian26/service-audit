FROM golang:1.25.6-alpine AS builder

# Install system dependencies
RUN apk add --no-cache git ca-certificates tzdata wget

# 1. Pre-download tools (Cached independently)
ARG GRPC_HEALTH_PROBE_VERSION=v0.4.37
RUN wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# 2. Setup Private Repo
ARG GITHUB_TOKEN
ENV GOPRIVATE=github.com/wahyunurdian26/*
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

WORKDIR /app

# 3. Dependencies (Using BuildKit Cache Mount for Go modules)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 4. Copy source
COPY . .

# 5. Tester Stage
FROM builder AS tester
RUN --mount=type=cache,target=/go/pkg/mod \
    go test ./... -v

# 6. Binary Builder Stage
FROM builder AS binary-builder
RUN --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=linux go build -o main .

# 7. Final Output Stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# Copy binary and health probe from previous stages
COPY --from=binary-builder /app/main .
COPY --from=builder /bin/grpc_health_probe /usr/bin/local/grpc_health_probe

EXPOSE 8080 6661
CMD ["./main"]
