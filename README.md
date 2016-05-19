I did this because I found no reasonable way to sign ethereum transaction offline - on an air-gapped device, or at least on booted Tails. I just want to avoid typing cold storage password on my workstation.

I wanted to do this on Tails, as that's quite guaranteed to be clean and forgetful. First I tried to do it in Python as in http://vitalik.ca/files/python\_cold\_wallet\_instructions.txt, but it turned out to be too much pain to run the pyethereum on Tails. Then there is the Icebox but that's javascript: https://github.com/ConsenSys/icebox :/.

These tools are using go-ethereum (geth) codebase.

As geth account manager works with keystore directory, the dir is what you need to pass to these tools.

## Tools

### txsign

`txsign` will create and sign a transaction that you can later publish on https://etherscan.io/pushTx, or via geth or pyethereum if you want.

### newaccount
`newaccount` will ask for password, and then generate new account with private key from randomness from crypto/rand (for details, try to find function storeNewKey in go-ethereum code).

It needs a directory (possibly empty) where the new account is stored.


