package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/amanechibana/veritas-chain/identity"
	"github.com/spf13/cobra"
)

// identityCmd represents the identity command
var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Identity and key management commands",
	Long:  `Commands for generating and inspecting signer keys for Veritas Chain nodes.`,
}

// identityKeygenCmd generates a new P-256 private key and prints the signer env and address
var identityKeygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a new signer key",
	Long:  `Generate a new P-256 private key and print SIGNER_PRIVATE_KEY_HEX and derived address.`,
	Run: func(cmd *cobra.Command, args []string) {
		priv, addr := generateKeyAndAddress()
		hexD := hex.EncodeToString(priv.D.Bytes())
		fmt.Println("Generated signer key:")
		fmt.Printf("  SIGNER_PRIVATE_KEY_HEX=%s\n", hexD)
		fmt.Printf("  Address=%s\n", addr)
	},
}

func generateKeyAndAddress() (*ecdsa.PrivateKey, string) {
	curve := elliptic.P256()
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}
	pubBytes := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	id := &identity.Identity{PrivateKey: *priv, PublicKey: pubBytes}
	return priv, string(id.Address())
}

func init() {
	rootCmd.AddCommand(identityCmd)

	// Add identity subcommands
	identityCmd.AddCommand(identityKeygenCmd)
}
