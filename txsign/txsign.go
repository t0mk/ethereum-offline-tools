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

	// gas limit in txs 18.05.2016:  0.00000002 Ether == 200000000000 wei
	gas_limit_wei_float := big.NewFloat(math.Pow10(18) * 0.00000002)
	if (!gas_limit_wei_float.IsInt()) || gas_limit_wei_float.IsInf() {
		log.Fatal("gas limit number fail")
	}
	gas_limit_wei, accuracy := gas_limit_wei_float.Int(nil)
	if accuracy != big.Exact {
		log.Fatal("gas limit number conversion fail")
	}

	// gas prices in Mist 18.5.2016
	// 0.000658154
	// 0.00072398
	// 0.000796378
	// 0.000875984
	// 0.000963646
	// 0.00106
	// 0.001166
	// 0.0012826
	// 0.00141086
	// 0.001551946
	// 0.00170713
	gas_price_wei_float := big.NewFloat(math.Pow10(18) * 0.000963646)
	if (!gas_price_wei_float.IsInt()) || gas_price_wei_float.IsInf() {
		log.Fatal("gas price number fail")
	}
	gas_price_wei, accuracy := gas_price_wei_float.Int(nil)
	if accuracy != big.Exact {
		log.Fatal("gas price number conversion fail")
	}

	amount_wei_float := big.NewFloat(math.Pow10(18) * amount_ether)
	if (!amount_wei_float.IsInt()) || amount_wei_float.IsInf() {
		log.Fatal("amount number fail")
	}
	amount_wei, accuracy := amount_wei_float.Int(nil)
	if accuracy != big.Exact {
		log.Fatal("amount number conversion fail")
	}

	// http://ethereum.stackexchange.com/questions/3386/create-and-sign-offline-raw-transactions
	//
	// transaction := types.NewTransaction(nonce, recipient, value, gasLimit, gasPrice, input)
	// signature, _ := crypto.Sign(transaction.SigHash().Bytes(), key)
	// signed, _ := tx.WithSignature(signature)

	tx := types.NewTransaction(nonce, to_addr, amount_wei, gas_limit_wei,
		gas_price_wei, nil)

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

}
