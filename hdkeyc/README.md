# hdkeyc
Cli for working with HD wallet keys.


## Installation
```
go get github.com/whyrusleeping/hdkeyutils/hdkeyc
```

## Usage
```
NAME:
   hdkeyc - A command line utility for manipulating HD wallet keys

USAGE:
   hdkeyc [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
     priv     utilities for working with hd private keys
     pub      tools for working with hd public keys
     gen      generate an HD wallet key
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Importing Keys

### Bitcoin

To import a private key:
```
$ hdkeyc priv child mymasterpriv.key 101 > key101.wif
$ bitcoin-cli importprivkey `cat key101.wif` false
```

To get a public key:
```
hdkeyc pub child mymasterpub.key 101
```

### Zcash

To import a private key:
```
$ hdkeyc priv child mymaster.key 101 > key101.wif
$ zcash-cli importprivkey `cat key101.wif` false
```

To get a public key:
```
hdkeyc pub child mymasterpub.key 101 --format=zec
```

### Ethereum

To import a private key:
```
$ hdkeyc priv child mymaster.key 101 --format=eth > key101.eth
$ geth import key101.eth 
```

To get a public key:
```
hdkeyc pub child mymasterpub.key 101 --format=eth
```

## Multisig

`hdkeyc` can be used as an HDM wallet interface.

### Creating an address
To create an HDM wallet, you will need to collect master public keys from all
parties who will be signers in the multisig.  To generate an address with a
given set of keys, run:

```
hdkeyc msig addr -m=4 -n=5 --index=19 $MPUB1 $MPUB2 $MPUB3 $MPUB4 $MPUB5
```

Where m and n are the 'm of n' parameters for the multisig setup, and index is
the HD key index to use for this address. 

You can now send transactions to that generated address.

### Spending from multisig
Spending from a multisig address is a bit more complicated than a normal
transaction.  First, you will need your redeem script. You can retrieve this by
calling `hdkeyc msig redeem-script` with the exact same parameters (in the
exact same order) as the `addr` command you ran to generate the address you are
spending from. 

```
hdkeyc msig redeem-script -m=4 -n=5 --index=19 $MPUB1 $MPUB2 $MPUB3 $MPUB4
$MPUB5
```

Next, you will need the hash of the transaction you want to spend, the address
you want to send the funds to, and the exact value (minus tx fees) that
transaction is for.

Once you've gathered all that, create the base transaction with:

``` 
hdkeyc msig mktx $REDEEM $TXINHASH $TARGET $VALUE 
```

This will output the base transaction data in hex. This is what we will be
gathering signatures over.

Next, you will need to gather 'm' signatures from keyholders. A keyholder can
sign the transaction by running:

```
hdkeyc msig sign myprivate.key $TXDATA
```

This will output the raw signature in hex. Once you have gathered 'm' of these,
you can create the final transaction. You can do this with:

```
hdkeyc msig finish $TXDATA $SIG1 $SIG2 ... $SIGM
```

The output of that command will the the hex encoded raw transaction that you
can then send off into the network.
