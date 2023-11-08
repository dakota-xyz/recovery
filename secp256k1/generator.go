package secp256k1

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/dakota-xyz/recovery/registry"
	decred "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/sha3"
)

func init() {
	registry.RegisterGenerator("ELLIPTIC_CURVE_SECP256K1", GenerateKey)
}

// GenerateKey satisfies the registry type
func GenerateKey(seed io.Reader) (string, string, error) {
	derivedPrivateKey, err := decred.GeneratePrivateKeyFromRand(seed)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}
	publicKeyBytes := elliptic.Marshal(decred.S256(), derivedPrivateKey.PubKey().X(), derivedPrivateKey.PubKey().Y())
	address, err := PublicKeyToAddress(publicKeyBytes)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}
	return string(address), "0x" + derivedPrivateKey.Key.String(), nil
}

// PublicKeyToAddress converts the public key to a valid ethereum address format
func PublicKeyToAddress(publicKey []byte) ([]byte, error) {
	hasher := sha3.NewLegacyKeccak256()
	if _, err := hasher.Write(publicKey[1:]); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	rawAddress := hasher.Sum(nil)[12:]
	// 40 bytes from the address plus 2 bytes for the 0x prefix
	var buf [42]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], rawAddress)

	// compute checksum
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2] //nolint:gomnd // Part of the algorithm
		if i%2 == 0 {
			hashByte = hashByte >> 4 //nolint:gomnd // Part of the algorithm
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return buf[:], nil
}
