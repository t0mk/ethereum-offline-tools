package main

import (
	"fmt"
	"log"
	"math"
	"math/big"

	"os"
	"strconv"

	t0mkcommon "../common"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func usage() {
	fmt.Println("This creates new transaction and signs it with account from given keystore. You can then publish the transaction e.g. on https://etherscan.io/pushTx")
	fmt.Println("Usage:")
	fmt.Println(os.Args[0], "keystore_dir fromAddress toAddress amount_in_ether nonce")
	fmt.Println("Example:")
	fmt.Println(os.Args[0], "/home/tomk/.ethereum/keystore 0x124c3349b52adb253ff92838684b6261d793be5d 0x326fca031b2af033fcc5ad0d44981a54f917b3b2 0.75 4")
	os.Exit(2)
}

func main() {
	// keystore_dir, fromAddress, toAddress, amount, nonce

	if len(os.Args) != 6 {
		usage()
	}
	amountEther, err := strconv.ParseFloat(os.Args[4], 64)
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
	fromArgAddr := os.Args[2]
	toArgAddr := os.Args[3]
	fromAddr := common.HexToAddress(fromArgAddr)
	toAddr := common.HexToAddress(toArgAddr)
	// am := accounts.NewManager(os.Args[1], 262144, 1)
	ksPath := os.Args[1]
	ks := keystore.NewKeyStore(ksPath, 262144, 1)
	if !ks.HasAddress(fromAddr) {
		log.Fatal(fmt.Sprintf("Address %s is not in %s", fromArgAddr, ksPath))
	}
	accName := accounts.Account{Address: fromAddr}
	acc, err := ks.Find(accName)

	basePrompt := "Password to unlock " + fromArgAddr
	p := t0mkcommon.GetCheckedPassword(basePrompt)

	fmt.Println("Password read and checked. Now trying to unlock account", acc)

	err = ks.TimedUnlock(acc, p, 0)
	if err == nil {
		fmt.Println("Success unlocking account for", fromArgAddr)
	} else {
		log.Fatal(err)
	}

	ws := ks.Wallets()
	log.Println(ws)

	gasLimit := new(big.Int)
	_, err = fmt.Sscan("121000", gasLimit)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice := new(big.Int)
	_, err = fmt.Sscan("20000000000", gasPrice)
	if err != nil {
		log.Fatal(err)
	}

	amountWeiFloat := big.NewFloat(math.Pow10(18) * amountEther)
	if (!amountWeiFloat.IsInt()) || amountWeiFloat.IsInf() {
		log.Fatal("amount number fail")
	}
	amountWei, accuracy := amountWeiFloat.Int(nil)
	if accuracy != big.Exact {
		log.Fatal("amount number conversion fail")
	}

	maxTxCostWei := new(big.Int).Mul(gasPrice, gasLimit)
	maxTxCostEth := new(big.Float).Mul(new(big.Float).SetInt(maxTxCostWei), big.NewFloat(math.Pow10(-18)))

	// http://ethereum.stackexchange.com/questions/3386/create-and-sign-offline-raw-transactions
	//
	// transaction := types.NewTransaction(nonce, recipient, value, gasLimit, gasPrice, input)
	// signature, _ := crypto.Sign(transaction.SigHash().Bytes(), key)
	// signed, _ := tx.WithSignature(signature)

	tx := types.NewTransaction(nonce, toAddr, amountWei, gasLimit, gasPrice, nil)

	signedTx, err := ks.SignTx(acc, tx, nil)
	if err != nil {
		log.Fatal(err)
	}

	/*
		signer := types.HomesteadSigner{}

		signature, err := am.Sign(fromAddr, tx.SigHash(signer).Bytes())
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println(signature)
		//fmt.Println(tx)

		signedTx, err := tx.WithSignature(signer, signature)
		if err != nil {
			log.Fatal(err)
		}
	*/

	fmt.Println("transaction:")
	fmt.Println(signedTx)
	fmt.Println()
	fmt.Println("Maximum transaction cost in Eth in 7 decimals:", maxTxCostEth.Text('f', 6))
	fmt.Println()
	fmt.Println("You can publish the Hex at https://etherscan.io/pushTx")

	fmt.Println()
	fmt.Println("After you publish, you can see your transaction at:")
	fmt.Println(fmt.Sprintf("https://etherscan.io/tx/0x%x", signedTx.Hash()))

}
