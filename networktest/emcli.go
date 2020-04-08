// This software is Copyright (c) 2019 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package networktest

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

const (
	EMCLI = "./build/emcli"

	// gjson paths
	QGetMintableEUR = "value.mintable.#(denom==\"eeur\").amount"
	QGetBalanceEUR  = "value.Account.value.coins.#(denom==\"eeur\").amount"
)

type Emcli struct {
	node     string
	chainid  string
	keystore *KeyStore
}

func (cli Emcli) QueryIssuers() ([]byte, error) {
	return execCmdAndCollectResponse(cli.addQueryFlags("q", "issuers"))
}

func (cli Emcli) QueryInflation() ([]byte, error) {
	return execCmdAndCollectResponse(cli.addQueryFlags("q", "inflation"))
}

func (cli Emcli) AuthorityCreateIssuer(authority, issuer Key, denoms ...string) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "authority", "create-issuer", authority.name, issuer.GetAddress(), strings.Join(denoms, ","))
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) AuthorityDestroyIssuer(authority, issuer Key) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "authority", "destroy-issuer", authority.name, issuer.GetAddress())
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) AuthoritySetMinGasPrices(authority Key, minGasPrices string, params ...string) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "authority", "set-gas-prices", authority.name, minGasPrices)
	args = append(args, params...)
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) QueryMinGasPrices() ([]byte, error) {
	args := cli.addQueryFlags("query", "authority", "gas-prices")
	return execCmdAndCollectResponse(args)
}

func (cli Emcli) QueryTransaction(txhash string) ([]byte, error) {
	args := cli.addQueryFlags("query", "tx", txhash)
	return execCmdAndCollectResponse(args)
}

func (cli Emcli) QueryRewards(delegator string) (gjson.Result, error) {
	args := cli.addQueryFlags("query", "distribution", "rewards", delegator)

	bz, err := execCmdAndCollectResponse(args)
	if err != nil {
		return gjson.Result{}, err
	}

	return gjson.ParseBytes(bz), nil
}

// NOTE Hardcoded to eeur for now.
func (cli Emcli) QueryAccount(account string) (balance, mintable int, err error) {
	args := cli.addQueryFlags("query", "account", account)
	bz, err := execCmdAndCollectResponse(args)
	if err != nil {
		return 0, 0, err
	}

	queryresponse := gjson.ParseBytes(bz)

	v := queryresponse.Get(QGetBalanceEUR)
	balance, _ = strconv.Atoi(v.Str)

	v = queryresponse.Get(QGetMintableEUR)
	if v.Exists() {
		mintable, _ = strconv.Atoi(v.Str)
	}

	return
}

func (cli Emcli) QueryAccountJson(account string) ([]byte, error) {
	args := cli.addQueryFlags("query", "account", account)
	return execCmdAndCollectResponse(args)
}

func (cli Emcli) QueryMarketInstruments() ([]byte, error) {
	args := cli.addQueryFlags("query", "market", "instruments")
	return execCmdAndCollectResponse(args)
}

func (cli Emcli) QueryMarketInstrument(source, destination string) ([]byte, error) {
	args := cli.addQueryFlags("query", "market", "instrument", source, destination)
	return execCmdAndCollectResponse(args)
}

func (cli Emcli) QueryValidators() (gjson.Result, error) {
	args := cli.addQueryFlags("query", "staking", "validators")
	bz, err := execCmdAndCollectResponse(args)
	if err != nil {
		return gjson.Result{}, err
	}

	return gjson.ParseBytes(bz), nil
}

func (cli Emcli) QueryDelegations(account string) ([]byte, error) {
	args := cli.addQueryFlags("query", "staking", "delegations", account)
	return execCmdAndCollectResponse(args)
}

func (cli Emcli) IssuerIncreaseMintableAmount(issuer, liquidityprovider Key, amount string) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "issuer", "increase-mintable", issuer.name, liquidityprovider.GetAddress(), amount)
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) IssuerRevokeMinting(issuer, liquidityprovider Key) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "issuer", "revoke-mint", issuer.name, liquidityprovider.GetAddress())
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) IssuerDecreaseMintableAmount(issuer, liquidityprovider Key, amount string) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "issuer", "decrease-mintable", issuer.name, liquidityprovider.GetAddress(), amount)
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) IssuerSetInflation(issuer Key, denom string, inflation string) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "issuer", "set-inflation", issuer.name, denom, inflation)
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) LiquidityProviderMint(key Key, amount string) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "liquidityprovider", "mint", key.name, amount)
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) LiquidityProviderBurn(key Key, amount string) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "liquidityprovider", "burn", key.name, amount)
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) MarketAddOrder(key Key, source, destination, cid string) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "market", "add", source, destination, cid, "--from", key.name)
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) MarketCancelOrder(key Key, cid string) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "market", "cancel", cid, "--from", key.name)
	return execCmdWithInput(args, KeyPwd)
}

func (cli Emcli) UnjailValidator(key Key) (string, bool, error) {
	args := cli.addTransactionFlags("tx", "slashing", "unjail", "--from", key.name)
	return execCmdWithInput(args, KeyPwd)
}

func extractTxHash(bz []byte) (txhash string, success bool, err error) {
	json := gjson.ParseBytes(bz)

	txhashjson := json.Get("txhash")
	successjson := gjson.ParseBytes(bz).Get("logs.0.success")

	if !txhashjson.Exists() || !successjson.Exists() {
		return "", false, fmt.Errorf("could not find status fields in response %v", string(bz))
	}

	return txhashjson.Str, successjson.Bool(), nil
}

func execCmdWithInput(arguments []string, input string) (string, bool, error) {
	cmd := exec.Command(EMCLI, arguments...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", false, err
	}

	_, err = io.WriteString(stdin, input+"\n")
	if err != nil {
		return "", false, err
	}

	//fmt.Println(" *** Running command: ", EMCLI, strings.Join(arguments, " "))
	bz, err := cmd.CombinedOutput()
	//fmt.Println(" *** Output", string(bz))
	if err != nil {
		return "", false, err
	}

	return extractTxHash(bz)
}

func execCmdAndCollectResponse(arguments []string) ([]byte, error) {
	//fmt.Println(" *** Running command: ", EMCLI, strings.Join(arguments, " "))
	bz, err := exec.Command(EMCLI, arguments...).CombinedOutput()
	//fmt.Println(" *** Output: ", string(bz))
	return bz, err
}

func (cli Emcli) addQueryFlags(arguments ...string) []string {
	return cli.addNetworkFlags(arguments)
}

func (cli Emcli) addTransactionFlags(arguments ...string) []string {
	arguments = append(arguments,
		"--home", cli.keystore.path,
		"--yes",
	)

	return cli.addNetworkFlags(arguments)
}

func (cli Emcli) addNetworkFlags(arguments []string) []string {
	return append(arguments,
		"--node", cli.node,
		"--chain-id", cli.chainid,
		"--output", "json",
	)
}
