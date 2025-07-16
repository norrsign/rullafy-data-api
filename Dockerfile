# syntax=docker/dockerfile:1.7
FROM golang:1.24.5 AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /goapi  .

# ──────────────────────────────────────────────────────────────────────────────
# Runtime layer
# ──────────────────────────────────────────────────────────────────────────────
FROM alpine:3.20
ENV GOAPI_PORT=8080
ENV GOAPI_LISTEN=0.0.0.0

COPY --from=builder /goapi /usr/local/bin/goapi

RUN printf '%s\n' \
  '#!/bin/sh' \
  'set -e' \
  'port=${GOAPI_PORT:-8080}' \
  'listen=${GOAPI_LISTEN:-0.0.0.0}' \
  'cmd="/usr/local/bin/goapi server start --port $port --listen $listen"' \
  '[ -n "$GOAPI_REALM_URL" ] && cmd="$cmd --realm-url $GOAPI_REALM_URL"' \
  'exec $cmd' \
  > /usr/local/bin/entrypoint.sh \
  && chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
EXPOSE 8080
