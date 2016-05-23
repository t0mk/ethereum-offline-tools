package main

import (
	"fmt"
	"log"
	"math"
	"math/big"

	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	t0mkcommon "github.com/t0mk/ethereum-offline-tools/common"
)

func usage() {
	fmt.Println("This creates new transaction and signs it with account from given keystore. You can then publish the transaction e.g. on https://etherscan.io/pushTx")
	fmt.Println("Usage:")
	fmt.Println(os.Args[0], "keystore_dir from_address to_address amount_in_ether nonce")
	fmt.Println("Example:")
	fmt.Println(os.Args[0], "/home/tomk/.ethereum/keystore 0x124c3349b52adb253ff92838684b6261d793be5d 0x326fca031b2af033fcc5ad0d44981a54f917b3b2 0.75 4")
	os.Exit(2)
}

func main() {
	// keystore_dir, from_address, to_address, amount, nonce

	if len(os.Args) != 6 {
		usage()
	}
	amount_ether, err := strconv.ParseFloat(os.Args[4], 64)
	if err != nil {
		log.Fatal(err)
	}
	nonce, err := strconv.ParseUint(os.Args[5], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	stat, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if !stat.IsDir() {
		log.Fatal(os.Args[1], "is not a directory")
	}
	from_arg_addr := os.Args[2]
	to_arg_addr := os.Args[3]
	from_addr := common.HexToAddress(from_arg_addr)
	to_addr := common.HexToAddress(to_arg_addr)
	am := accounts.NewManager(os.Args[1], 262144, 1)

	acc := t0mkcommon.GetAccountForAddr(am, from_addr)
	base_prompt := "Password to unlock " + from_arg_addr
	p := t0mkcommon.GetCheckedPassword(base_prompt)

	fmt.Println("Password read and checked. Now trying to unlock account", acc.File)

	err = am.TimedUnlock(acc, p, 0)
	if err == nil {
		fmt.Println("Success unlocking account for", from_arg_addr)
	} else {
		log.Fatal(err)
	}

	// gas for tx is 21000
	gas_limit := common.String2Big("121000")

	// >>> 20 * 10**9
	// 20000000000
	gas_price := common.String2Big("20000000000")

	// Is using String2Big retarded in this context? I guess, but I refuse to
	// investigate the int types and their limits now in 2016 ...

	amount_wei_float := big.NewFloat(math.Pow10(18) * amount_ether)
	if (!amount_wei_float.IsInt()) || amount_wei_float.IsInf() {
		log.Fatal("amount number fail")
	}
	amount_wei, accuracy := amount_wei_float.Int(nil)
	if accuracy != big.Exact {
		log.Fatal("amount number conversion fail")
	}

	max_tx_cost_wei := new(big.Int).Mul(gas_price, gas_limit)
	max_tx_cost_eth := new(big.Float).Mul(new(big.Float).SetInt(max_tx_cost_wei), big.NewFloat(math.Pow10(-18)))

	// http://ethereum.stackexchange.com/questions/3386/create-and-sign-offline-raw-transactions
	//
	// transaction := types.NewTransaction(nonce, recipient, value, gasLimit, gasPrice, input)
	// signature, _ := crypto.Sign(transaction.SigHash().Bytes(), key)
	// signed, _ := tx.WithSignature(signature)

	tx := types.NewTransaction(nonce, to_addr, amount_wei, gas_limit,
		gas_price, nil)

	signature, err := am.Sign(from_addr, tx.SigHash().Bytes())
	if err != nil {
		log.Fatal(err)
	}

	signed_tx, err := tx.WithSignature(signature)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("transaction:")
	fmt.Println(signed_tx)
	fmt.Println()
	fmt.Println("Maximum transaction cost in Eth in 7 decimals:", max_tx_cost_eth.Text('f', 6))
	fmt.Println()
	fmt.Println("You can publish the Hex at https://etherscan.io/pushTx")

	fmt.Println()
	fmt.Println("After you publish, you can see your transaction at:")
	fmt.Println(fmt.Sprintf("https://etherscan.io/tx/0x%x", signed_tx.Hash()))

}
