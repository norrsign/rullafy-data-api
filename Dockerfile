FROM goapi:latest AS builder

RUN echo "Listing /app/goapi:" && ls -la /app/goapi


WORKDIR /app/norrsign/rullafy-data-api

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /goapi  .

# ──────────────────────────────────────────────────────────────────────────────
# Runtime layer
# ──────────────────────────────────────────────────────────────────────────────
FROM alpine:3.20


COPY --from=builder /goapi /usr/local/bin/goapi

RUN printf '%s\n' \
  '#!/bin/sh' \
  'set -e' \
  'cmd="/usr/local/bin/goapi server start"' \
  'exec $cmd' \
  > /usr/local/bin/entrypoint.sh \
  && chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
EXPOSE 8080
