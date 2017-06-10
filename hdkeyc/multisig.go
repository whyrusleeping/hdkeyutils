package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	keychain "github.com/btcsuite/btcutil/hdkeychain"
	cli "github.com/urfave/cli"
	addrs "github.com/whyrusleeping/hdkeyutils/addrs"
	mktx "github.com/whyrusleeping/mktx"
)

var msigCmd = cli.Command{
	Name:  "msig",
	Usage: "manipulate HD multisig wallets",
	Subcommands: []cli.Command{
		msigRedeemScriptCmd,
		msigAddrCmd,
		msigMkSpendTxCmd,
		msigSignTxCmd,
		msigFinishTxSigCmd,
	},
}

var msigRedeemScriptCmd = cli.Command{
	Name:  "redeem-script",
	Usage: "create a multisig redeem script for the given keys",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "m",
			Usage: "'m' parameter for m of n",
		},
		cli.IntFlag{
			Name:  "n",
			Usage: "'n' parameter for m of n",
		},
		cli.IntFlag{
			Name:  "index",
			Usage: "hd wallet index of pubkeys to use",
			Value: -1,
		},
	},
	Action: func(c *cli.Context) error {
		msigscr, err := cmdMsigScript(c)
		if err != nil {
			return err
		}

		fmt.Printf("%x", msigscr)
		return nil
	},
}

func cmdMsigScript(c *cli.Context) ([]byte, error) {
	index := c.Int("index")
	if index < 0 {
		return nil, fmt.Errorf("must specify an index")
	}

	m := c.Int("m")
	if m <= 0 {
		return nil, fmt.Errorf("must specify a value for m")
	}
	n := c.Int("n")
	if n <= 0 {
		return nil, fmt.Errorf("must specify a value for n")
	}

	var keys [][]byte
	for _, keystr := range c.Args() {
		k, err := keychain.NewKeyFromString(keystr)
		if err != nil {
			return nil, err
		}
		child, err := k.Child(uint32(index))
		if err != nil {
			return nil, err
		}

		pubk, err := child.ECPubKey()
		if err != nil {
			return nil, err
		}

		keys = append(keys, pubk.SerializeCompressed())
	}

	return mktx.MakeMultisig(m, n, keys), nil
}

var msigAddrCmd = cli.Command{
	Name:  "addr",
	Usage: "create an HDM wallet address",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "m",
			Usage: "'m' parameter for m of n",
		},
		cli.IntFlag{
			Name:  "n",
			Usage: "'n' parameter for m of n",
		},
		cli.IntFlag{
			Name:  "index",
			Usage: "hd wallet index of pubkeys to use",
			Value: -1,
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "output format for key (btc or zec)",
			Value: "btc",
		},
		cli.BoolFlag{
			Name:  "testnet",
			Usage: "print testnet addrs",
		},
	},
	Action: func(c *cli.Context) error {
		msigscr, err := cmdMsigScript(c)
		if err != nil {
			return err
		}
		if c.Bool("testnet") {
			addrs.BitcoinP2SHPrefix = addrs.BitcoinTestnetP2SHPrefix
			addrs.ZcashP2SHPrefix = addrs.ZcashTestnetP2SHPrefix
		}

		var prefix []byte
		switch c.String("format") {
		case "btc":
			prefix = addrs.BitcoinP2SHPrefix
		case "zec":
			prefix = addrs.ZcashP2SHPrefix
		default:
			return fmt.Errorf("unsupported format: %s", c.String("format"))
		}

		hash := addrs.HashSha256Ripe160(msigscr)
		fmt.Println(addrs.Base58Check(hash, prefix))
		return nil
	},
}

var msigMkSpendTxCmd = cli.Command{
	Name:      "mktx",
	Usage:     "create a bare multisig spend transaction",
	ArgsUsage: "<redeem-script> <txhash> <target addr> <value>",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "prevoutindex",
			Usage: "specify the index of the previous output",
		},
	},
	Action: func(c *cli.Context) error {
		redeemscr := c.Args()[0]
		txhash := c.Args()[1]
		target := c.Args()[2]
		svalue := c.Args()[3]
		value, err := strconv.ParseInt(svalue, 10, 64)
		if err != nil {
			return err
		}

		prevtx, err := chainhash.NewHashFromStr(txhash)
		if err != nil {
			return err
		}

		redeem, err := hex.DecodeString(redeemscr)
		if err != nil {
			return err
		}

		addr, err := btcutil.DecodeAddress(target, nil)
		if err != nil {
			return err
		}

		pkscript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			return err
		}

		outp := wire.NewOutPoint(prevtx, uint32(c.Int("prevoutindex")))
		txin := wire.NewTxIn(outp, redeem)
		tx := wire.NewMsgTx(1)
		tx.AddTxIn(txin)

		txo := wire.NewTxOut(value, pkscript)
		tx.AddTxOut(txo)

		buf := new(bytes.Buffer)
		tx.Serialize(buf)
		fmt.Println(hex.EncodeToString(buf.Bytes()))

		return nil
	},
}

var msigSignTxCmd = cli.Command{
	Name:  "sign",
	Usage: "sign a multisig transaction",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "index",
			Usage: "index of HD key to sign with",
			Value: -1,
		},
	},
	Action: func(c *cli.Context) error {
		keyfi := c.Args()[0]
		txdata := c.Args()[1]

		index := c.Int("index")
		if index < 0 {
			return fmt.Errorf("must specify an index")
		}

		priv, err := loadPrivKey(keyfi)
		if err != nil {
			return err
		}
		sigkey, err := priv.Child(uint32(index))
		if err != nil {
			return err
		}
		ecpriv, err := sigkey.ECPrivKey()
		if err != nil {
			return err
		}

		rawtx, err := hex.DecodeString(txdata)
		if err != nil {
			return err
		}

		rawtx = append(rawtx, 1, 0, 0, 0)
		dhb := chainhash.DoubleHashB(rawtx)

		sig, err := ecpriv.Sign(dhb)
		if err != nil {
			return err
		}

		fmt.Printf("%x\n", sig.Serialize())
		return nil
	},
}

var msigFinishTxSigCmd = cli.Command{
	Name:  "finish",
	Usage: "complete a multisig spend tx with signatures",
	Action: func(c *cli.Context) error {
		txdata := c.Args()[0]

		var sigs [][]byte
		for _, sig := range c.Args()[1:] {
			rawsig, err := hex.DecodeString(sig)
			if err != nil {
				return err
			}

			sigs = append(sigs, rawsig)
		}

		rawtx, err := hex.DecodeString(txdata)
		if err != nil {
			return err
		}

		tx := wire.NewMsgTx(1)
		if err := tx.Deserialize(bytes.NewReader(rawtx)); err != nil {
			return err
		}

		sb := txscript.NewScriptBuilder()
		sb.AddOp(txscript.OP_0)
		for _, sig := range sigs {
			sb.AddData(append(sig, 1))
		}

		redeem := tx.TxIn[0].SignatureScript
		sb.AddData(redeem)

		outscript, err := sb.Script()
		if err != nil {
			return err
		}

		tx.TxIn[0].SignatureScript = outscript

		buf := new(bytes.Buffer)
		tx.Serialize(buf)
		fmt.Println(hex.EncodeToString(buf.Bytes()))

		return nil
	},
}
