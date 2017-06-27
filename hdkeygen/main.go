package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	chaincfg "github.com/btcsuite/btcd/chaincfg"
	keychain "github.com/btcsuite/btcutil/hdkeychain"
)

var chainParams = &chaincfg.Params{
	// From chaincfg
	HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
	HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub
}

func fatal(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, f, args...)
	os.Exit(1)
}

func main() {
	randsrc := flag.String("randsrc", "", "filename of alternative randomness source")
	randbytes := flag.Int("randlen", 8192, "number of bytes of randomness to read from randomness source")
	output := flag.String("output", "output.key", "name of keyfile to output")
	seedhex := flag.String("seedhex", "", "optionally specify entire seed data in hex")
	flag.Parse()

	var r io.Reader = rand.Reader
	if *randsrc == "" {
		fmt.Println("  > No alt random source given, using go's 'crypto/rand.Reader'")
	} else {
		fi, err := os.Open(*randsrc)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer fi.Close()
		r = fi
	}

	fmt.Println("Starting Less Insecure Key Generator...")

	var seed []byte
	shx := *seedhex
	if len(shx) > 0 {
		sdata, err := hex.DecodeString(shx)
		if err != nil {
			fatal("failed to decode seedhex: %s", err)
		}

		seed = sdata
		fmt.Printf("Using seed data of \"%x\"\n", seed)
	} else {
		// Read data from users keyboard
		fmt.Println("Please enter some randomness, press Ctrl+d when youre done")
		userRandom := new(bytes.Buffer)
		io.Copy(userRandom, os.Stdin)

		fmt.Printf("Read %d bytes of random data from the keyboard\n", userRandom.Len())

		// Now get some randomness from either crypto/rand, or whatever the user
		// wants to use instead
		fmt.Printf("Now reading %d bytes of random data from cryptographic randomness source\n", *randbytes)
		otherRand := make([]byte, *randbytes)
		_, err := io.ReadFull(r, otherRand)
		if err != nil {
			fatal("Error while reading randomness: %s\n", err)
		}

		// Hash first the otherRand data then the users input
		h := sha512.New()
		h.Write(otherRand)
		h.Write(userRandom.Bytes())
		seed = h.Sum(nil)

		// TODO: maybe don't print this? its basically the private key
		fmt.Printf("Generated master key seed from randomness:\n%x\n", seed)
	}

	// Create the private key from the seed we generated
	fmt.Println("Now creating master private key...")
	newmasterkey, err := keychain.NewMaster(seed, chainParams)
	if err != nil {
		fatal("Error generating master key from seed: %s\n", err)
	}

	fmt.Printf("Key created, writing out to %s\n", *output)
	outfi, err := os.Create(*output)
	if err != nil {
		fatal("Failed to open output file: %s\n", err)
	}
	defer outfi.Close()

	_, err = outfi.Write([]byte(newmasterkey.String()))
	if err != nil {
		fatal("Error writing key to output file: %s\n", err)
	}

	// just to be safe and leave it in fewer random memory addresses
	newmasterkey.Zero()

	fmt.Println("Success!")
}
