I did this because I found no reasonable way to sign ethereum transaction offline - on an air-gapped device, or at least on booted Tails. I just want to avoid typing cold storage password on my workstation.

I wanted to do this on Tails, as that's quite guaranteed to be clean and forgetful. First I tried to do it in Python as in http://vitalik.ca/files/python_cold_wallet_instructions.txt, but it turned out to be too much pain to run the pyethereum on Tails. Then, there is the Icebox https://github.com/ConsenSys/icebox but that's javascript, meaning the browser will get to see the password too :/.

Then I realized I can just put it together in Go. That way, I can theoretically easily cross-compile and run it on different archs (RP, some handheld, ...). Also, I haven't found 64-bit x86 build of Tails, only 32-bit (i386), so I need to cross-compile anyway if I want to run it on Tails.

The tools are usign the go-ethereum (geth) codebase from https://github.com/ethereum/go-ethereum

As geth account manager works with keystore directory, the dir is what you need to pass to these tools.

# Tools

The tools can be easily cross-compiled with xgo (`go get github.com/karalabe/xgo`):

```
xgo --go=latest --targets=linux/arm-7 -v github.com/t0mk/ethereum-offline-tools/txsign
xgo --go=latest --targets=linux/arm-7 -v github.com/t0mk/ethereum-offline-tools/newaccount
```

IIUC xgo is building the code in container from pulled Docker image. This sucks for paranoid people I guess. But I think running it on air-gapped device should be OK.

## txsign

`txsign` will create and sign a transaction that you can later publish on https://etherscan.io/pushTx, or via geth or pyethereum if you want.

### gas price and gas limit

I am not sure of how these are interpreted. I just took values from recent transactions on the blockchain. See the source.

### Nonce

You need to set nonce for the transaction manually. Nonce is, IIUC, the count of transaction from an address. First transaction from new address has nonce 0, second has nonce 1, etc.

It's reasonable to find out nonce for your sending address from eth node with updated blockchain. You can see some info about all your accounts with geth with javascript source `counts.js` in the root of this repo.

If you have Ethereum-Wallet running:

```
~/Ethereum-Wallet-linux64-0-7-3/resources/node/geth/geth --exec 'loadScript("counts.js")'  --maxpeers 0 --nodiscover --networkid 3301 attach
```

If you have no geth running but you are confident that your local copy of blockchain is recent enough:

```
~/Ethereum-Wallet-linux64-0-7-3/resources/node/geth/geth --exec 'loadScript("counts.js")'  --maxpeers 0 --nodiscover --networkid 3301 console
```

Your nonce is `count+1` of the account you want to send from.

Naturally, if you have never sent anything from this address before, the nonce will be zero.

## newaccount
`newaccount` will ask for password, and then generate new account with private key from randomness from crypto/rand (for details, try to find function storeNewKey in go-ethereum code).

It needs a directory (possibly empty) where the new account is stored.


