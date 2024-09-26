package cmd

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	_ "github.com/dakota-xyz/recovery/ed25519"
	"github.com/dakota-xyz/recovery/registry"
	_ "github.com/dakota-xyz/recovery/secp256k1"

	"github.com/hashicorp/vault/shamir"
	"golang.org/x/crypto/hkdf"
)

type KeyDerivation struct {
	// Key derivation parameters
	AddressSubID string `json:"address_sub_id"`
	NetworkID    string `json:"network_id"`
	Curve        string `json:"curve"`

	// Useful metadata for the customer
	Wallet             string   `json:"wallet"`
	Address            string   `json:"address"`
	AccountName        string   `json:"account_name"`
	CompatibleNetworks []string `json:"compatible_networks"`
}
type KeyMap struct {
	OrganizationId string          `json:"organization_id"`
	Keys           []KeyDerivation `json:"keys"`
}

const keyDerivationSeparator = ":"

// SerializeKeyDerivationParameters takes the key derivation parameters and
// serializes them as expected
func SerializeKeyDerivationParameters(organizationId, addressSubId, networkId string) ([]byte, error) {
	hash := sha256.New()
	fullStr := strings.Join([]string{organizationId, addressSubId, networkId}, keyDerivationSeparator)
	if _, err := hash.Write([]byte(fullStr)); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return hash.Sum(nil), nil
}

// RecoverSeed takes a base seed and the key derivation parameters, and returns
// the new seed as a io.Reader
func RecoverSeed(seed []byte, organizationID, addressSubID, networkID string) (io.Reader, error) {
	serializedDerivationParams, err := SerializeKeyDerivationParameters(
		organizationID,
		addressSubID,
		networkID,
	)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	privateKeySeed := hkdf.New(sha512.New, seed, nil, serializedDerivationParams)
	return privateKeySeed, nil
}

// Recover reconstructs a users keys and prints them as a CSV to the io.Writer
func Recover(targetWriter io.Writer, shard1, shard2, keyMapReader io.Reader) error {
	shard1Bytes, err := io.ReadAll(shard1)
	if err != nil {
		return fmt.Errorf("failed to read shard1: %w", err)
	}
	shard2Bytes, err := io.ReadAll(shard2)
	if err != nil {
		return fmt.Errorf("failed to read shard2: %w", err)
	}

	// Now that we have the both the client shard and ours, we can reconstruct the client key
	dek, err := shamir.Combine([][]byte{shard1Bytes, shard2Bytes})
	if err != nil {
		return fmt.Errorf("failed to reconstruct the client data encryption key: %w", err)
	}

	var keymap KeyMap
	if err := json.NewDecoder(keyMapReader).Decode(&keymap); err != nil {
		return fmt.Errorf("failed to unmarshal keymap: %w", err)
	}
	w := csv.NewWriter(targetWriter)
	w.Write([]string{"Account", "Address", "PrivateKey", "Wallet", "Compatible Networks..."})
	for _, kdp := range keymap.Keys {
		generator, exists := registry.Generators[kdp.Curve]
		if !exists {
			return fmt.Errorf("unsupported curve %s", kdp.Curve)
		}
		privateKeySeed, err := RecoverSeed(dek, keymap.OrganizationId, kdp.AddressSubID, kdp.NetworkID)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		address, privateKey, err := generator(privateKeySeed)
		if err != nil {
			return fmt.Errorf("failed to derive key: %w", err)
		}

		if kdp.Wallet == "ONCHAIN" {
			address = kdp.Address
		}

		w.Write([]string{kdp.AccountName, address, privateKey, kdp.Wallet, strings.Join(kdp.CompatibleNetworks, ",")})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}
	return nil
}
