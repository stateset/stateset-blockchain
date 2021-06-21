FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/stateset/stateset-blockchain

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make install

# Final image
FROM alpine:edge

ENV stateset /stateset

# Install ca-certificates
RUN apk add --update ca-certificates

RUN addgroup stateset && \
    adduser -S -G stateset stateset -h "$STATESET"

USER stateset

WORKDIR $STATESET

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/statesetd /usr/bin/statesetd

# Run statesetd by default, omit entrypoint to ease using container with statesetcli
CMD ["statesetd"]