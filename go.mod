module github.com/hukkinj1/cosmosvanity

go 1.17

require (
	github.com/cosmos/cosmos-sdk v0.45.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.35.1
)

require (
	github.com/btcsuite/btcd v0.22.0-beta // indirect
	github.com/cosmos/btcutil v1.0.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sasha-s/go-deadlock v0.2.1-0.20190427202633-1595213edefa // indirect
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

// see https://github.com/cosmos/cosmos-sdk/issues/8469
replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
