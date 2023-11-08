package cmd_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/dakota-xyz/recovery/cmd"
	"github.com/dakota-xyz/recovery/registry"
	"github.com/hashicorp/vault/shamir"
	"github.com/stretchr/testify/require"
)

const seed = "43e04e7cc735c3a4f611ead5475ad910fd2e97398d4284d99b73e69267e4c6c6ab8518ce0e93423e4de3aa4768c13eab7df39c9e13f93e531431b8f28ab3359d4c82700589f5ddd1e721a98fdf777aab8304b3f9b112aeb7c77bed76d7640404c788b2f844d1c0974fdf1697b24f435fdd89bd5bd91893a3b66ecfdb5434c5bd"
const backupJSON = `{
	"organization_id": "e65ccfaa-01b7-4a00-aec5-0fcc25d7eba7",
	"encrypted_backup": "-----BEGIN PGP MESSAGE-----\nVersion: GopenPGP 2.7.4\nComment: https://gopenpgp.org\n\nwcFMA/iE9y3x4R3yARAAh9CmKXIK4dBzwioe1bJo/igJ2hRDiZ19qTqQ8w/oBluD\nSs7I+sHMOpznsMsm4/9qMv1n5xmAEtIbhP84jZYfodN8PisFWlW9z/HwG0V2BEdX\n+0CRveq28uyUvsyxQLmb3HLkdIIj/CfN/MhBHgzs1z/Sxt6zvcLVrMESkd6joTc8\nZqlofDndWk8TdixgrGOk7HmgYgq6EsZCzY6RQOqmBiD1c62LE2F5AP873avPGcde\nYAXpW6rMDnF3pFPBXak7SJxaHKz5ITxMkjh3ObLo4GDIFsUQE8smAgeNuvGSjkXo\ngyaSnwKsfH1txIg/UokkzzUbbtD1VsEJoaS7L11qC2FCQXPaw2CKw4F7eagoqEUN\nVhX/WWi09jttbJybE3Ywt/p6XSZChCjgqLafozVHuLvApy3XxACCSlO+W9Tow50N\nnGAQZxtwgMC8lsKWE7xAFzdXfh7dYMdgyp7ERTzX5AaqQ2brDVfhwBbU6xnsPofd\nROUjyV+oAY1m5g1e15Z+ri60RTQfq/CQdYdUBI8KtceREsSVdDwVh6kc/4uXRY8D\nhWkbu9FSGHEVxHUBcgPykUw/IH7fnXvq3G81V7nRjDoRbdP028nnqNFsMgsU2vth\nb71zi0dHjccnti0sS/XdHWzagTj9h+q+hdQrCqB5AcxWYiRGxffUZTIf45S0k27S\nsgFwfN9fO/0Crun1TATHrcDp+iUdGhqGd/wRfhZxvl0rsoYC22wL94iq94Y52Yyt\n1EwZ3aSV2JsN05mbTuE39cX3ze4pE7cZhTB2t25oKL8rrwU4V9uYrDIr0tsCiEgn\nFnMOOMsVCYX5KST21LAGc41AZ5pqfxX9ggZgSDwtrSlsBI45GsgEVmP0ttVq8era\nvn9z8DqHNQ9kQB7dfUW8lxJzH694jQAl8TXCHEC0HUpYkDg=\n=y3F/\n-----END PGP MESSAGE-----",
	"keys": [
		{
		"address_sub_id": "21d3969c-8a56-46d9-be38-b53c21294e54",
		"network_id": "solana-devnet",
		"curve": "ELLIPTIC_CURVE_ED25519"
		},
		{
		"address_sub_id": "a09d020a-1bc8-47b6-a208-0fd47cd05e66",
		"network_id": "ethereum-hardhat",
		"curve": "ELLIPTIC_CURVE_SECP256K1"
		}
	]
}`

func TestFixture(t *testing.T) {
	inputs := []struct {
		OrganizationID     string
		AddressSubID       string
		Network            string
		Curve              string
		ExpectedAddress    string
		ExpectedPrivateKey string
	}{
		{
			OrganizationID:     "e65ccfaa-01b7-4a00-aec5-0fcc25d7eba7",
			AddressSubID:       "a09d020a-1bc8-47b6-a208-0fd47cd05e66",
			Network:            "ethereum-hardhat",
			Curve:              "ELLIPTIC_CURVE_SECP256K1",
			ExpectedAddress:    "0x0D7ad5799E3DB77c8258b9700E4f94Fcb092C64B",
			ExpectedPrivateKey: "0x222d55b028c7896058d28af1d44c55d45264c470f2a93e7b013076e68b7bfa25",
		},
		{
			OrganizationID:     "e65ccfaa-01b7-4a00-aec5-0fcc25d7eba7",
			AddressSubID:       "21d3969c-8a56-46d9-be38-b53c21294e54",
			Network:            "solana-devnet",
			Curve:              "ELLIPTIC_CURVE_ED25519",
			ExpectedAddress:    "4tZpnxbJbkCDFFCpbmb4y7wsH366kxeb57R8owi67qi8",
			ExpectedPrivateKey: "2tFuN9PCkTYsDV6rq8RauJZEmyBs7x8rLoSAFYD5JcQMCzVzStq45VeUVDDghGqXaYm8muC8YECzgoqTkyPph8gp",
		},
	}
	t.Parallel()
	for _, input := range inputs {
		testName := fmt.Sprintf("%s %s %s %s", input.OrganizationID, input.AddressSubID, input.Network, input.Curve)
		t.Run(testName, func(t *testing.T) {
			// Given a generator
			generator, exists := registry.Generators[input.Curve]
			require.True(t, exists, "no generator exists for curve %s", input.Curve)

			// When we generate the key
			seedBytes, err := hex.DecodeString(seed)
			require.NoError(t, err)
			privateKeySeed, err := cmd.RecoverSeed(seedBytes, input.OrganizationID, input.AddressSubID, input.Network)
			require.NoError(t, err)
			address, privateKey, err := generator(privateKeySeed)
			require.NoError(t, err)

			// Then it matches expectation
			require.Equal(t, input.ExpectedPrivateKey, privateKey)
			require.Equal(t, input.ExpectedAddress, address)
		})
	}
}

func TestRecover(t *testing.T) {
	// Given a sharded seed
	seedBytes, err := hex.DecodeString(seed)
	require.NoError(t, err)
	parts, err := shamir.Split(seedBytes, 3, 2)
	require.NoError(t, err)
	part1, part2 := parts[0], parts[1]

	// When we recover
	var targetWriter strings.Builder
	err = cmd.Recover(&targetWriter, bytes.NewReader(part1), bytes.NewReader(part2), strings.NewReader(backupJSON))
	require.NoError(t, err)

	// Then we find the expected addresses and keys
	csvContents := targetWriter.String()
	require.Contains(t, csvContents, "0x0D7ad5799E3DB77c8258b9700E4f94Fcb092C64B")
	require.Contains(t, csvContents, "0x222d55b028c7896058d28af1d44c55d45264c470f2a93e7b013076e68b7bfa25")
	require.Contains(t, csvContents, "4tZpnxbJbkCDFFCpbmb4y7wsH366kxeb57R8owi67qi8")
	require.Contains(t, csvContents, "2tFuN9PCkTYsDV6rq8RauJZEmyBs7x8rLoSAFYD5JcQMCzVzStq45VeUVDDghGqXaYm8muC8YECzgoqTkyPph8gp")
}
