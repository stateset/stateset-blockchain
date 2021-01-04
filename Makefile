PACKAGES=$(shell GO111MODULE=on go list -mod=readonly ./...)

MODULES = account agreement bank contact factoring invoice purchaseorder loan market slashing

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=statesetd \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

define \n


endef

benchmark:
	@go test -mod=readonly -bench=. $(PACKAGES)

buidl: build

build: build_cli build_daemon

download:
	go mod download

build_cli:
	@go build $(BUILD_FLAGS) -o bin/statesetcli cmd/statesetcli/*.go

build_daemon:
	@go build $(BUILD_FLAGS) -o bin/statesetd cmd/statesetd/*.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/statesetd cmd/statesetd/*.go
	GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/statesetcli cmd/statesetcli/*.go

doc:
	@echo "--> Wait a few seconds and visit http://localhost:6060/pkg/github.com/stateset/stateset-blockchain/"
	godoc -http=:6060

export:
	@bin/statesetd export

create-wallet:
	bin/statesetcli keys add validator --home ~/.octopus

init:
	rm -rf ~/.statesetd
	bin/statesetd init statenode $(shell bin/statesetcli keys show validator -a --home ~/.octopus)
	bin/statesetd add-genesis-account $(shell bin/statesetcli keys show validator -a --home ~/.octopus) 10000000000ustates
	bin/statesetd gentx --name=validator --amount 10000000000ustates --home-client ~/.octopus
	bin/statesetd collect-gentxs

install:
	@go install $(BUILD_FLAGS) ./cmd/statesetd
	@go install $(BUILD_FLAGS) ./cmd/statesetcli
	@echo "Installed statesetd and statesetcli ..."
	@statesetd version --long

reset:
	bin/statesetd unsafe-reset-all

restart: build_daemon reset start

start:
	bin/statesetd start --inv-check-period 10 --log_level "main:info,state:info,*:error,app:info,account:info,statebank:info,agreement:info,invoice:info,loan:info,market:info,stateslashing:info,statefactoring:info"

check:
	@echo "--> Running golangci"
	@golangci-lint run --tests=false --skip-files=\\btest_common.go

dep_graph: ; $(foreach dir, $(MODULES), godepgraph -s -novendor github.com/stateset/stateset-blockchain/x/$(dir) | dot -Tpng -o x/$(dir)/dep.png${\n})

install_tools_macos:
	brew install dep && brew upgrade dep
	brew install golangci/tap/golangci-lint
	brew upgrade golangci/tap/golangci-lint

go_test:
	@go test $(PACKAGES)

test: go_test

test_cover:
	@go test $(PACKAGES) -v -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic
	@go tool cover -html=coverage.txt

version:
	@bin/statesetd version --long

########################################
### Local validator nodes using docker and docker-compose

build-docker-statesetdnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: localnet-stop
	@if ! [ -f build/node0/statesetd/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/statesetd:Z stateset/statesetdnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

########################################

.PHONY: benchmark buidl build build_cli build_daemon check dep_graph test test_cover update_deps \
build-docker-statesetdnode localnet-start localnet-stop