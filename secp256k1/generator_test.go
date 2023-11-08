package secp256k1_test

import (
	"bytes"
	"testing"

	"github.com/dakota-xyz/recovery/cmd"
	"github.com/dakota-xyz/recovery/secp256k1"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGenerateKey(t *testing.T) {
	// Given a seed
	organizationID := uuid.NewString()
	addressSubID := uuid.NewString()
	networkID := "ethereum-mainnet"
	seed, err := cmd.SerializeKeyDerivationParameters(organizationID,
		addressSubID, networkID)
	require.NoError(t, err)

	// When we generate the key
	privateKey, address, err := secp256k1.GenerateKey(bytes.NewReader(seed))

	// Then it succeeds
	require.NoError(t, err)
	require.NotEmpty(t, address)
	require.NotEmpty(t, privateKey)
}
