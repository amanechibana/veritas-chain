package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "veritas",
	Short: "Veritas Chain CLI",
	Long: `Veritas Chain is a proof-of-authority blockchain for university certificate verification.

This CLI provides commands to manage blockchain nodes, identities, and certificates.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Veritas Chain!")
		fmt.Println("Use 'veritas --help' to see available commands.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags that apply to all commands
	rootCmd.PersistentFlags().StringP("config", "c", "", "Config file (default is $HOME/.veritas.yaml)")
}
