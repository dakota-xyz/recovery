package registry_test

import (
	"testing"

	"github.com/dakota-xyz/recovery/ed25519"
	"github.com/dakota-xyz/recovery/registry"
	"github.com/dakota-xyz/recovery/secp256k1"
	"github.com/stretchr/testify/require"
)

func TestRegisterGenerator(t *testing.T) {
	// Test case 1: Register a generator
	id := "ELLIPTIC_CURVE_SECP256K1"
	generator := secp256k1.GenerateKey
	registry.RegisterGenerator(id, generator)

	// Verify that the generator is registered
	_, exists := registry.Generators[id]
	require.True(t, exists, "Failed to register generator")

	// Test case 2: Register a duplicate generator
	generator = secp256k1.GenerateKey
	registry.RegisterGenerator(id, generator)

	// Verify that the generator is still registered
	_, exists = registry.Generators[id]
	require.True(t, exists, "Failed to register duplicate generator")

	// Test case 3: Register a new generator
	id = "ELLIPTIC_CURVE_ED25519"
	generator = ed25519.GenerateKey
	registry.RegisterGenerator(id, generator)

	// Verify that the new generator is registered
	_, exists = registry.Generators[id]
	require.True(t, exists, "Failed to register new generator")

	// Test case 4: Fetch non-existing generator
	_, exists = registry.Generators["NON_EXISTING"]

	// Verify that the generator is not registered
	require.False(t, exists, "Failed to fetch non-existing generator")
}
