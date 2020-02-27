module github.com/e-money/em-ledger

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200224190540-bb7e803cb929
	github.com/emirpasic/gods v1.12.0
	github.com/gorilla/mux v1.7.4
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.4.3
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.1
	github.com/tendermint/tm-db v0.4.0
	github.com/tidwall/gjson v1.3.2
	github.com/tidwall/sjson v1.0.4
)

// replace github.com/cosmos/cosmos-sdk => ./tmpvendor/cosmos-sdk

// replace github.com/tendermint/tendermint => ./tmpvendor/tendermint
