FROM golang:1.24.4-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .

RUN mkdir -p /app \
    && go build -o /app/submgr ./cmd/api

RUN go build -o /app/healthcheck ./cmd/healthcheck

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/submgr /app/submgr
COPY --from=builder /app/healthcheck /app/healthcheck

WORKDIR /app
ENTRYPOINT ["/app/submgr"]

