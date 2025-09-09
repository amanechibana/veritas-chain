package cmd

import (
	"fmt"

	"github.com/amanechibana/veritas-chain/blockchain"
	"github.com/spf13/cobra"
)

// blockchainCmd represents the blockchain command
var blockchainCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "Blockchain management commands",
	Long:  `Commands for managing the blockchain, including viewing blocks, validating the chain, and managing certificates.`,
}

// blockchainListCmd represents the blockchain list command
var blockchainListCmd = &cobra.Command{
	Use:   "list",
	Short: "List blockchain data",
	Long:  `List various blockchain data (blocks, certificates, etc.)`,
}

// blockchainListBlocksCmd represents the blockchain list blocks command
var blockchainListBlocksCmd = &cobra.Command{
	Use:   "blocks",
	Short: "List all blocks",
	Long:  `List all blocks in the blockchain with optional filtering and formatting.`,
	Run: func(cmd *cobra.Command, args []string) {
		showAll, _ := cmd.Flags().GetBool("all")
		limit, _ := cmd.Flags().GetInt("limit")

		if limit == 0 {
			limit = 100
		}

		if showAll {
			fmt.Println("Showing all blocks (this may take a while)")
		} else {
			fmt.Printf("Showing first %d blocks\n", limit)
		}

		chain := blockchain.ContinueBlockchain()
		defer chain.Close()

		iter := chain.Iterator()
		blockCount := 0

		for blockCount < limit {
			block := iter.Next()
			fmt.Printf("Block %d: Height=%d, Hash=%x, Signature=%x\n",
				blockCount, block.Height, block.Hash, block.Signature)

			blockCount++

			// Stop when we reach the genesis block (PrevHash is empty)
			if len(block.PrevHash) == 0 {
				break
			}
		}

	},
}

// blockchainValidateCmd represents the blockchain validate command
var blockchainValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the blockchain",
	Long:  `Validate the integrity of the entire blockchain by checking all blocks and their signatures.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !blockchain.DBExists() {
			fmt.Println("No blockchain found. Please initialize a blockchain first.")
			return
		}

		chain := blockchain.ContinueBlockchain()
		defer chain.Close()

		fmt.Println("Validating blockchain...")

		err := chain.ValidateChain()
		if err != nil {
			fmt.Println("✗ Blockchain validation failed:", err)
			return
		}

		fmt.Println("✓ All blocks validated successfully")
		fmt.Println("✓ All signatures verified")
		fmt.Println("✓ Chain integrity confirmed")
	},
}

// blockchainInfoCmd represents the blockchain info command
var blockchainInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show blockchain information",
	Long:  `Display general information about the blockchain, including statistics and network status.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !blockchain.DBExists() {
			fmt.Println("No blockchain found. Please initialize a blockchain first.")
			return
		}

		chain := blockchain.ContinueBlockchain()
		defer chain.Close()

		blockchainStats := chain.GetStats()

		fmt.Println("Blockchain Information:")
		fmt.Printf("  Total Blocks: %d\n", blockchainStats.BlockCount)
		fmt.Printf("  Total Certificates: %d\n", blockchainStats.CertificateCount)
		fmt.Println("  Network Status: idk maybe healthy?")
	},
}

func init() {
	rootCmd.AddCommand(blockchainCmd)

	// Add blockchain subcommands
	blockchainCmd.AddCommand(blockchainListCmd)
	blockchainCmd.AddCommand(blockchainValidateCmd)
	blockchainCmd.AddCommand(blockchainInfoCmd)

	// Add list subcommands
	blockchainListCmd.AddCommand(blockchainListBlocksCmd)

	// Persistent flags for list command
	blockchainListCmd.PersistentFlags().BoolP("all", "a", false, "Show all data, including hidden")

	// Block listing flags
	blockchainListBlocksCmd.Flags().IntP("limit", "l", 10, "Maximum number of blocks to show")
}
