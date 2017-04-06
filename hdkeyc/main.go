package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	chaincfg "github.com/btcsuite/btcd/chaincfg"
	btcutil "github.com/btcsuite/btcutil"
	b58 "github.com/btcsuite/btcutil/base58"
	keychain "github.com/btcsuite/btcutil/hdkeychain"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	cli "github.com/urfave/cli"
	"golang.org/x/crypto/ripemd160"
)

func main() {
	app := cli.NewApp()
	app.Usage = "A command line utility for manipulating HD wallet keys"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		privKeyCmd,
		pubKeyCmd,
		genKeyCmd,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var genKeyCmd = cli.Command{
	Name:            "gen",
	Usage:           "generate an HD wallet key",
	SkipFlagParsing: true,
	Action: func(c *cli.Context) error {
		f, err := exec.LookPath("hdkeygen")
		if err != nil {
			return fmt.Errorf("could not find 'hdkeygen' binary: %s", err)
		}

		cmd := exec.Command(f, c.Args()...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		return cmd.Run()
	},
}

var privKeyCmd = cli.Command{
	Name:  "priv",
	Usage: "utilities for working with hd private keys",
	Subcommands: []cli.Command{
		getMasterPubCmd,
		getChildPrivKeyCmd,
	},
}

var getMasterPubCmd = cli.Command{
	Name:  "getmasterpub",
	Usage: "derive the master public key from the given master private key",
	Action: func(c *cli.Context) error {
		if !c.Args().Present() {
			return fmt.Errorf("must pass in private key file")
		}

		data, err := ioutil.ReadFile(c.Args().First())
		if err != nil {
			return err
		}

		key, err := keychain.NewKeyFromString(string(data))
		if err != nil {
			return err
		}

		if !key.IsPrivate() {
			return fmt.Errorf("given key was not a private key")
		}

		pubk, err := key.Neuter()
		if err != nil {
			return err
		}

		fmt.Print(pubk.String())
		return nil
	},
}

var getChildPrivKeyCmd = cli.Command{
	Name:  "child",
	Usage: "derive a child private key",
	Description: `Derive a child private key from the given heirarchically deterministic
   private key and print it out.

   By default, it outputs in wallet import format for use by bitcoin and zcash.
   Optionally, you may pass the --format flag with a parameter of 'eth' to 
   signal that it should output a raw ecdsa key for use by ethereum.`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "output format for key (wif or eth)",
			Value: "wif",
		},
	},
	Action: func(c *cli.Context) error {
		format := c.String("format")

		if len(c.Args()) != 2 {
			return fmt.Errorf("must pass in private key and index")
		}

		i, err := strconv.Atoi(c.Args()[1])
		if err != nil {
			return err
		}

		data, err := ioutil.ReadFile(c.Args().First())
		if err != nil {
			return err
		}

		key, err := keychain.NewKeyFromString(string(data))
		if err != nil {
			return err
		}

		if !key.IsPrivate() {
			return fmt.Errorf("given key was not a private key")
		}

		childpriv, err := key.Child(uint32(i))
		if err != nil {
			return err
		}

		privk, err := childpriv.ECPrivKey()
		if err != nil {
			return err
		}

		switch format {
		case "wif":
			wif, err := btcutil.NewWIF(privk, &chaincfg.MainNetParams, false)
			if err != nil {
				return err
			}

			fmt.Println(wif.String())
		case "eth":
			fmt.Printf("%x\n", privk.Serialize())
		}
		return nil
	},
}

var pubKeyCmd = cli.Command{
	Name:  "pub",
	Usage: "tools for working with hd public keys",
	Subcommands: []cli.Command{
		getChildPubKeyCmd,
	},
}

var getChildPubKeyCmd = cli.Command{
	Name:  "child",
	Usage: "derive a child public key",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "output format for key (btc, zec, or eth)",
			Value: "btc",
		},
	},
	Action: func(c *cli.Context) error {
		enc := c.String("format")

		if len(c.Args()) != 2 {
			return fmt.Errorf("must pass in public key and index")
		}

		i, err := strconv.Atoi(c.Args()[1])
		if err != nil {
			return err
		}

		data, err := ioutil.ReadFile(c.Args().First())
		if err != nil {
			return err
		}

		key, err := keychain.NewKeyFromString(string(data))
		if err != nil {
			return err
		}

		if key.IsPrivate() {
			return fmt.Errorf("given key was a private key, not public")
		}

		childpub, err := key.Child(uint32(i))
		if err != nil {
			return err
		}

		addr, err := childpub.Address(&chaincfg.MainNetParams)
		if err != nil {
			return err
		}

		switch enc {
		case "btc":
			fmt.Println(b58.CheckEncode(addr.Hash160()[:], 0))
		case "zec":
			ecpk, err := childpub.ECPubKey()
			if err != nil {
				return err
			}
			uncomp := ecpk.SerializeUncompressed()
			shad := sha256.Sum256(uncomp)
			h := ripemd160.New()
			h.Write(shad[:])
			ripemd := h.Sum(nil)
			d := append([]byte{0x1c, 0xb8}, ripemd...)
			fmt.Println(b58.CheckEncode(d, 0))
		case "eth":
			ecpubkey, err := childpub.ECPubKey()
			if err != nil {
				return err
			}

			addr := ethcrypto.PubkeyToAddress(*ecpubkey.ToECDSA())
			fmt.Println(addr.Hex())
		default:
			return fmt.Errorf("unrecognized output format: %s", enc)
		}
		return nil
	},
}
