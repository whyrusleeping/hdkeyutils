package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	chaincfg "github.com/btcsuite/btcd/chaincfg"
	btcutil "github.com/btcsuite/btcutil"
	keychain "github.com/btcsuite/btcutil/hdkeychain"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	cli "github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		privKeyCmd,
		pubKeyCmd,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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
			fmt.Println(addr.EncodeAddress())
		case "zec":
			fmt.Println("t" + addr.EncodeAddress())
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
