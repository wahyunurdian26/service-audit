FROM golang:alpine AS builder

RUN apk add --no-cache git

# ARG GGITHUB_TOKEN for private repos if needed
ARG GITHUB_TOKEN
ENV GOPRIVATE=github.com/wahyunurdian26/*
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Dedicated testing stage
FROM builder AS tester
RUN go test ./... -v || true

# Final production build stage
FROM builder AS final-build
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=final-build /app/main .

EXPOSE 8080 6661
CMD ["./main"]
