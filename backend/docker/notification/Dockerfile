FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /app/app ./cmd/notification

RUN touch /etc/log.log && chown 65534:65534 /etc/log.log

FROM scratch

WORKDIR /app

COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app/app .
COPY --from=builder --chown=nobody:nobody /etc/log.log /etc/log.log
COPY --from=builder /app/config.docker.yaml ./config.yaml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nobody

CMD ["/app/app"]
