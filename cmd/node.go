package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/amanechibana/veritas-chain/blockchain"
	"github.com/amanechibana/veritas-chain/identity"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Node management commands",
	Long:  `Commands for managing Veritas Chain nodes, including starting, stopping, and configuring nodes.`,
}

// nodeInteractiveCmd represents the node interactive command
var nodeInteractiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start node in interactive mode",
	Long: `Start a Veritas Chain node in interactive mode.
This allows you to interact with the blockchain through a command-line interface.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load .env if present
		_ = godotenv.Load()

		fmt.Printf("Configuration:\n")

		// Load signer from env (required)
		signer, err := identity.LoadSignerFromEnv()
		if err != nil {
			fmt.Printf("Failed to load signer from env: %v\n", err)
			return
		}
		if signer == nil {
			fmt.Println("SIGNER_PRIVATE_KEY_HEX is required. Use 'veritas identity keygen' to generate one.")
			return
		}

		addr := string(signer.Address())
		fmt.Printf("  Address: %s\n", addr)

		// Compute per-signer DB path
		dbPath := filepath.Join("./tmp", "blocks_"+addr)
		fmt.Printf("  DB Path: %s\n", dbPath)

		// Optionally load authorized signers mapping and resolve name
		if _, err := os.Stat("authorized_signers.json"); err == nil {
			m, err := identity.LoadAuthorizedSigners("authorized_signers.json")
			if err == nil {
				if name, e2 := m.ResolveNameByAddress(addr); e2 == nil {
					fmt.Printf("  Resolved Name: %s\n", name)
				}
			}
		}

		// Initialize or continue blockchain
		var chain *blockchain.Blockchain
		if blockchain.DBExists(dbPath) {
			chain = blockchain.ContinueBlockchain(dbPath)
			fmt.Println("Loaded existing blockchain")
		} else {
			chain = blockchain.InitBlockchain(dbPath, signer)
			fmt.Println("Created new blockchain with genesis block")
		}
		defer chain.Close()

		// Start interactive mode
		startInteractiveMode(chain, signer)
	},
}

// startInteractiveMode starts the interactive terminal
func startInteractiveMode(chain *blockchain.Blockchain, signer identity.Signer) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== Veritas Chain Interactive Mode ===")
	fmt.Println("Type 'help' for available commands")
	fmt.Printf("Signer Address: %s\n", string(signer.Address()))
	fmt.Println("=====================================")

	for {
		fmt.Print("veritas> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		command := parts[0]

		switch command {
		case "help":
			showHelp()
		case "add":
			if len(parts) < 2 {
				fmt.Println("Usage: add <certificate1,certificate2,...>")
				continue
			}
			certificates := strings.Split(parts[1], ",")
			addBlock(chain, signer, certificates)
		case "list":
			listBlocks(chain)
		case "validate":
			validateChain(chain)
		case "stats":
			showStats(chain)
		case "exit", "quit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}
	}
}

func showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  add <cert1,cert2,...>  - Add a new block with certificates")
	fmt.Println("  list                   - List all blocks")
	fmt.Println("  validate               - Validate the blockchain")
	fmt.Println("  stats                  - Show blockchain statistics")
	fmt.Println("  help                   - Show this help message")
	fmt.Println("  exit/quit              - Exit interactive mode")
}

func addBlock(chain *blockchain.Blockchain, signer identity.Signer, certificates []string) {
	block, err := chain.AddBlock(certificates, signer)

	if err != nil {
		fmt.Printf("Failed to add block: %v\n", err)
		return
	}
	fmt.Printf("Block added successfully!\n")
	fmt.Printf("   Height: %d\n", block.Height)
	fmt.Printf("   Hash: %x\n", block.Hash)
	fmt.Printf("   Address: %s\n", string(block.UniversityAddress))
}

func listBlocks(chain *blockchain.Blockchain) {
	iter := chain.Iterator()
	blockCount := 0

	fmt.Println("Blockchain:")
	for blockCount < 10 { // Limit to 10 blocks
		block := iter.Next()
		fmt.Printf("Block %d: Height=%d, Hash=%x, Address=%s\n",
			blockCount, block.Height, block.Hash, string(block.UniversityAddress))

		blockCount++
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func validateChain(chain *blockchain.Blockchain) {
	err := chain.ValidateChain()
	if err != nil {
		fmt.Printf("  Chain validation failed: %v\n", err)
	} else {
		fmt.Println("  Chain validation successful")
	}
}

func showStats(chain *blockchain.Blockchain) {
	stats := chain.GetStats()
	fmt.Printf("Blockchain Statistics:\n")
	fmt.Printf("  Total Blocks: %d\n", stats.BlockCount)
	fmt.Printf("  Total Certificates: %d\n", stats.CertificateCount)
}

func init() {
	rootCmd.AddCommand(nodeCmd)

	// Add node subcommands
	nodeCmd.AddCommand(nodeInteractiveCmd)
}
