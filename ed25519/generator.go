package ed25519

import (
	goed25519 "crypto/ed25519"
	"fmt"
	"io"

	"github.com/dakota-xyz/recovery/registry"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/mr-tron/base58"
)

func init() {
	registry.RegisterGenerator("ELLIPTIC_CURVE_ED25519", GenerateKey)
}

// 3HaU3UsPY66wsNVRPFwvAkzxUftcMZSYPmE1B1gTPYjAgSRZh8uBTjesLTY4VR4KFcrLeRw6dVQuQ7AYXhKcFRiC
func GenerateKey(seed io.Reader) (string, string, error) {
	publicKey, privateKey, err := goed25519.GenerateKey(seed)
	if err != nil {
		return "", "", fmt.Errorf("%w", err)
	}
	pubkey := common.PublicKeyFromBytes(publicKey)
	return pubkey.ToBase58(), base58.Encode(privateKey), nil
}
