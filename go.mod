module github.com/stateset/stateset-blockchain

go 1.16

require (
	github.com/armon/go-metrics v0.4.1
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/cosmos-sdk v0.47.4
	github.com/cosmos/go-bip39 v1.0.0
	github.com/go-kit/kit v0.12.0
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.3
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/json-iterator/go v1.1.12
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.16.0
	github.com/regen-network/cosmos-proto v0.3.1
	github.com/spf13/cast v1.5.1
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.16.0
	github.com/stateset/stateset v0.0.0-20230828213553-a4a0dca73590
	github.com/stretchr/testify v1.8.4
	github.com/tendermint/go-amino v0.16.0
	github.com/tendermint/spm v0.1.4
	github.com/tendermint/tendermint v0.37.0-rc2
	github.com/tendermint/tm-db v0.6.7
	golang.org/x/crypto v0.11.0
	google.golang.org/grpc v1.56.2
	gopkg.in/yaml.v2 v2.4.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
