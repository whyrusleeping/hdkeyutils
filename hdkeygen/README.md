# hdkeygen

A simple (hopefully auditable) tool to generate hd wallets
in a slightly less insecure manner.

This tool is kept separate from the other utility `hdkeyc` to make
reading through it in its entirety easier to do.

## Installation

### Easy mode
```
go get .
```

### Trust 'why' mode
For this, make sure you have an ipfs daemon running.
```
$ # either install gx this way, or via prebuilt binaries on dist.ipfs.io
$ go get github.com/whyrusleeping/gx
$ go get github.com/whyrusleeping/gx-go
$ # Now, ensure gx uses ipfs to fetch data
$ export IPFS_API="http://localhost:5001"
$ gx install
$ gx-go rewrite
$ go install .
```

### Hard mode
- Download and install a version of go that you trust (given your threat model).
- Fetch a trusted copy of `github.com/btcsuite/btcd` and `github.com/btcsuite/btcutil` and install it into your `$GOPATH`
- run `go install`

## Usage
```
Usage of ./hdkeygen:
  -output string
        name of keyfile to output (default "output.key")
  -randlen int
        number of bytes of randomness to read from randomness source (default 8192)
  -randsrc string
        filename of alternative randomness source
```

Start by running:
```
$ hdkeygen
```

The program will prompt you to type some randomness, once done, close stdin with ctrl+d.

The program will then proceed to generate your key, and write it to the output file.

# License
MIT, whyrusleeping
