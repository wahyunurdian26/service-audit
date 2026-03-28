FROM golang:alpine AS builder

WORKDIR /app

# Copy the common shared module
COPY common ./common

# Copy the protobuf client module if applicable (transaction, gateway, etc might need it)
# We will copy the whole repo into builder since it's cleaner for local module resolution
COPY . .

WORKDIR /app/audit-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/audit-service/main .
CMD ["./main"]
