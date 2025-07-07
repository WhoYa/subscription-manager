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

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/submgr /app/submgr

WORKDIR /app
ENTRYPOINT ["/app/submgr"]

