package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "veritas",
		Short: "Veritas Chain CLI",
		Long:  "Command-line interface for Veritas Chain",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to Veritas Chain!")
			cmd.Help()
		},
	}

	var nodeStartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start a Veritas Chain node",
		Long:  "Start a Veritas Chain node",
		Run: func(cmd *cobra.Command, args []string) {
			// Get flag values
			port, _ := cmd.Flags().GetInt("port")
			university, _ := cmd.Flags().GetString("university")
			verbose, _ := cmd.Flags().GetBool("verbose")

			fmt.Printf("Starting Veritas Chain node...\n")
			fmt.Printf("Port: %d\n", port)
			fmt.Printf("University: %s\n", university)
			if verbose {
				fmt.Println("Verbose mode enabled")
			}
		},
	}

	// Add flags to the node start command
	nodeStartCmd.Flags().IntP("port", "p", 8080, "Port to run the node on")
	nodeStartCmd.Flags().StringP("university", "u", "harvard", "University name for this node")

	// Make port flag required
	nodeStartCmd.MarkFlagRequired("port")

	// Add a create identity command with different flag types
	var createIdentityCmd = &cobra.Command{
		Use:   "create-identity",
		Short: "Create a new university identity",
		Long:  "Create a new university identity with specified parameters",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			force, _ := cmd.Flags().GetBool("force")

			fmt.Printf("Creating identity for: %s\n", name)
			if force {
				fmt.Println("Force mode: overwriting existing identity")
			}
		},
	}

	// Add flags to create identity command
	createIdentityCmd.Flags().StringP("name", "n", "", "University name (required)")
	createIdentityCmd.Flags().BoolP("force", "f", false, "Force creation, overwrite existing")

	// Make name required
	createIdentityCmd.MarkFlagRequired("name")

	// Add a list command with persistent flags
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List blockchain data",
		Long:  "List various blockchain data (blocks, identities, etc.)",
	}

	// Add persistent flags that apply to all subcommands
	listCmd.PersistentFlags().StringP("format", "f", "table", "Output format (table, json, yaml)")
	listCmd.PersistentFlags().BoolP("all", "a", false, "Show all data, including hidden")

	// Add subcommands to list
	var listBlocksCmd = &cobra.Command{
		Use:   "blocks",
		Short: "List all blocks",
		Run: func(cmd *cobra.Command, args []string) {
			format, _ := cmd.Flags().GetString("format")
			showAll, _ := cmd.Flags().GetBool("all")
			limit, _ := cmd.Flags().GetInt("limit")

			fmt.Printf("Listing blocks in %s format\n", format)
			if showAll {
				fmt.Println("Showing all blocks")
			}
			fmt.Printf("Limit: %d\n", limit)
		},
	}

	listBlocksCmd.Flags().IntP("limit", "l", 10, "Maximum number of blocks to show")
	listCmd.AddCommand(listBlocksCmd)

	rootCmd.AddCommand(nodeStartCmd)
	rootCmd.AddCommand(createIdentityCmd)
	rootCmd.AddCommand(listCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
