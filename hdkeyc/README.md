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
