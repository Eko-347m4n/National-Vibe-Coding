package wallet

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	version            = byte(0x00)
	addressChecksumLen = 4
)

// Wallet menyimpan pasangan kunci publik dan privat
type Wallet struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  []byte
}

// NewWallet membuat dan mengembalikan Wallet baru
func NewWallet() *Wallet {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	return &Wallet{priv, pub}
}

// Address mengembalikan alamat wallet
func (w Wallet) Address() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)

	return []byte(address)
}

// HashPubKey melakukan hashing pada public key dengan ripemd160(sha256(pubkey))
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// checksum menghasilkan checksum untuk sebuah public key
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// DecodeAddress mendekode alamat base58 menjadi public key hash
func DecodeAddress(address string) []byte {
	pubKeyHash, err := base58.Decode(address)
	if err != nil {
		log.Panic(err)
	}
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	if bytes.Compare(actualChecksum, targetChecksum) != 0 {
		log.Panic("Error: invalid address checksum")
	}

	return pubKeyHash
}
