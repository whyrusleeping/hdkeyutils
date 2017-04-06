package addrs

import (
	"crypto/sha256"

	btcec "github.com/btcsuite/btcd/btcec"
	b58 "github.com/btcsuite/btcutil/base58"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ripemd160"
)

func KeyHashSha256Ripe160(ecpk *btcec.PublicKey) []byte {
	uncomp := ecpk.SerializeUncompressed()
	shad := sha256.Sum256(uncomp)
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

func EncodeBitcoinPubkey(k *btcec.PublicKey) string {
	val := KeyHashSha256Ripe160(k)
	return Base58Check(val, []byte{0})
}

func EncodeZcashPubkey(k *btcec.PublicKey) string {
	val := KeyHashSha256Ripe160(k)
	return Base58Check(val, []byte{0x1c, 0xb8})
}

func EncodeEthereumPubkey(k *btcec.PublicKey) string {
	addr := ethcrypto.PubkeyToAddress(*k.ToECDSA())
	return addr.Hex()
}
