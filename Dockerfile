FROM cosmwasm/go-ext-builder:0001-alpine AS rust-builder

ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python

WORKDIR /go/src/github.com/stateset/stateset-blockchain

COPY go.* /go/src/github.com/stateset/stateset-blockchain/

RUN apk add --no-cache git \
    && go mod download github.com/CosmWasm/go-cosmwasm \
    && export GO_WASM_DIR=$(go list -f "{{ .Dir }}" -m github.com/CosmWasm/go-cosmwasm) \
    && cd ${GO_WASM_DIR} \
    && cargo build --release --features backtraces --example muslc \
    && mv ${GO_WASM_DIR}/target/release/examples/libmuslc.a /lib/libgo_cosmwasm_muslc.a


FROM cosmwasm/go-ext-builder:0001-alpine AS go-builder

WORKDIR /go/src/github.com/stateset/stateset-blockchain

RUN apk add --no-cache git libusb-dev linux-headers

COPY . .
COPY --from=rust-builder /lib/libgo_cosmwasm_muslc.a /lib/libgo_cosmwasm_muslc.a

RUN BUILD_TAGS=muslc make update-swagger-docs build

FROM alpine:3

WORKDIR /root

COPY --from=go-builder /go/src/github.com/stateset/stateset-blockchain/build/statesetd /usr/local/bin/statesetd
COPY --from=go-builder /go/src/github.com/stateset/stateset-blockchain/build/statesetcli /usr/local/bin/statesetcli

CMD [ "statesetd", "--help" ]