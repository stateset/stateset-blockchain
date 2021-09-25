#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/statesetd/${BINARY:-statesetd}
ID=${ID:-0}
LOG=${LOG:-statesetd.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'statesetd' E.g.: -e BINARY=statesetd_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export STATESETDHOME="/statesetd/node${ID}/statesetd"

if [ -d "$(dirname "${STATESETDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${STATESETDHOME}" "$@" | tee "${STATESETDHOME}/${LOG}"
else
  "${BINARY}" --home "${STATESETDHOME}" "$@"
fi