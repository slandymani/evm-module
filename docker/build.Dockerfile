### build stage
FROM golang:1.21-bookworm AS build-env

WORKDIR /go/src/github.com/slandymani/evm-module

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

RUN go install github.com/MinseokOh/toml-cli@latest
RUN go install github.com/cosmos/cosmos-sdk/cosmovisor/cmd/cosmovisor@latest

### run stage
FROM alpine:3.18

WORKDIR /root

COPY --from=build-env /go/bin/toml-cli /usr/bin/toml-cli
COPY --from=build-env /go/bin/cosmovisor /usr/bin/cosmovisor
COPY --from=build-env /go/src/github.com/slandymani/evm-module/build/evmmd /usr/bin/evmmd

RUN apk add --no-cache \
    ca-certificates jq \
    curl bash \
    vim lz4 \
    tini \
    gcompat
    
RUN addgroup -g 1000 evmm \
    && adduser -S -h /home/evmm -D evmm -u 1000 -G evmm

USER 1000
WORKDIR /home/evmm

ENV DAEMON_NAME=evmmd
ENV DAEMON_HOME=/home/evmm/.evmmd
ENV DAEMON_ALLOW_DOWNLOAD_BINARIES=true
ENV DAEMON_RESTART_AFTER_UPGRADE=true
ENV UNSAFE_SKIP_BACKUP=false

RUN mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin $DAEMON_HOME/cosmovisor/upgrades \
    && cp /usr/bin/evmmd $DAEMON_HOME/cosmovisor/genesis/bin/evmmd

EXPOSE 26656 26657 1317 9090 8545 8546

ENTRYPOINT ["/sbin/tini", "--"]

CMD ["cosmovisor"]
