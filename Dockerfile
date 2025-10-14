FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM node:lts-trixie-slim AS node

FROM ghcr.io/astral-sh/uv:python3.13-trixie-slim

COPY --from=node /usr/local/bin/node /usr/local/bin/node
#COPY --from=node /usr/local/include/node /usr/local/include/node
COPY --from=node /usr/local/lib/node_modules /usr/local/lib/node_modules
RUN ln -s /usr/local/lib/node_modules/npm/bin/npm-cli.js /usr/local/bin/npm && \
    ln -s /usr/local/lib/node_modules/npm/bin/npx-cli.js /usr/local/bin/npx && \
    ln -s /usr/local/bin/node /usr/local/bin/nodejs

RUN apt-get update && \
    apt-get install -y --no-install-recommends curl wget libsqlcipher1 libsqlcipher-dev && \
    rm -rf /var/lib/apt/lists/*


COPY --from=builder /app/build/mcp-proxy /main
ENTRYPOINT ["/main"]
CMD ["--config", "/config/config.json"]
