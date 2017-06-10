package addrs

import (
	"crypto/sha256"
	"fmt"

	btcec "github.com/btcsuite/btcd/btcec"
	b58 "github.com/btcsuite/btcutil/base58"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160"
)

func KeyHashSha256Ripe160(ecpk *btcec.PublicKey) []byte {
	return HashSha256Ripe160(ecpk.SerializeCompressed())
}

func HashSha256Ripe160(data []byte) []byte {
	shad := sha256.Sum256(data)
	h := ripemd160.New()
	h.Write(shad[:])
	return h.Sum(nil)
}

func Base58Check(val, prefix []byte) string {
	val = append(prefix, val...)
	first := sha256.Sum256(val)
	chk := sha256.Sum256(first[:])
	return b58.Encode(append(val, chk[:4]...))
}

var BitcoinPrefix = []byte{0}
var BitcoinP2SHPrefix = []byte{5}
var BitcoinTestnetPrefix = []byte{0x6f}
var BitcoinTestnetP2SHPrefix = []byte{0xc4}

func EncodeBitcoinPubkey(k *btcec.PublicKey) string {
	val := KeyHashSha256Ripe160(k)
	return Base58Check(val, BitcoinPrefix)
}

var ZcashPrefix = []byte{0x1c, 0xb8}
var ZcashP2SHPrefix = []byte{0x1c, 0xbd}
var ZcashTestnetPrefix = []byte{0x1d, 0x25}
var ZcashTestnetP2SHPrefix = []byte{0x1c, 0xba}

func EncodeZcashPubkey(k *btcec.PublicKey) string {
	val := KeyHashSha256Ripe160(k)
	return Base58Check(val, ZcashPrefix)
}

func EncodeEthereumPubkey(k *btcec.PublicKey) string {
	addr := ethcrypto.PubkeyToAddress(*k.ToECDSA())
	return addr.Hex()
}

func DecodeZcashAddr(addr string) ([]byte, error) {
	data := b58.Decode(addr)
	if len(data) == 0 {
		return nil, fmt.Errorf("invalid base58")
	}

	if len(data) < 6 {
		return nil, fmt.Errorf("invalid zcash address")
	}
	return data[2 : len(data)-4], nil
}
