# hdkeyutils

Tools for working with HD keys.

This repo consists of two tools. One, `hdkeygen` for less-insecurely generating
private HD master keys. And another, `hdkeyc` for deriving child public and
private keys from the master private keys.


## Installation

### `hdkeygen`
Basic installation is simply:
```
go get -u github.com/whyrusleeping/hdkeyutils/hdkeygen
```

More detailed instructions can be found in the [hdkeygen README](./hdkeygen/README.md).

### `hdkeyc`
```
go get -u github.com/whyrusleeping/hdkeyutils/hdkeyc
```

### Install both at once
```
go get -u github.com/whyrusleeping/hdkeyutils/...
```

## Usage 

Usage instructions can be found in the READMEs for each tool.
