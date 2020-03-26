# Simple usage with a mounted data directory:
# > docker build -t stateset .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.statesetd:/root/.statesetd -v ~/.statesetcli:/root/.statesetcli stateset statesetd init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.statesetd:/root/.statesetd -v ~/.statesetcli:/root/.statesetcli stateset statesetd start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python

# Set working directory for the build
WORKDIR /go/src/github.com/stateset/stateset-blockchain

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make tools && \
    make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/statesetd /usr/bin/statesetd
COPY --from=build-env /go/bin/statesetcli /usr/bin/statesetcli

# Run statesetd by default, omit entrypoint to ease using container with statesetcli
CMD ["statesetd"]