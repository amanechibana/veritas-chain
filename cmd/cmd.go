package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/amanechibana/veritas-chain/blockchain"
	"github.com/spf13/cobra"
)

// Node represents a Veritas Chain node
type Node struct {
	Port       int
	University string
	Verbose    bool
	Blockchain *blockchain.Blockchain
	Server     *http.Server
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// AddBlockRequest represents a request to add a block
type AddBlockRequest struct {
	Certificates []string `json:"certificates"`
	Address      string   `json:"address"`
}

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
		Long:  "Start a Veritas Chain node with HTTP API",
		Run: func(cmd *cobra.Command, args []string) {
			// Get flag values
			port, _ := cmd.Flags().GetInt("port")
			university, _ := cmd.Flags().GetString("university")
			verbose, _ := cmd.Flags().GetBool("verbose")

			// Create and start the node
			node := &Node{
				Port:       port,
				University: university,
				Verbose:    verbose,
			}

			if err := node.Start(); err != nil {
				log.Fatal(err)
			}
		},
	}

	// Add flags to the node start command
	nodeStartCmd.Flags().IntP("port", "p", 8080, "Port to run the node on")
	nodeStartCmd.Flags().StringP("university", "u", "harvard", "University name for this node")
	nodeStartCmd.Flags().BoolP("verbose", "v", false, "Enable verbose logging")

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

	// Add CLI commands that work with running nodes
	var addBlockCmd = &cobra.Command{
		Use:   "add-block",
		Short: "Add a block to the blockchain",
		Long:  "Add a block to the blockchain (works with running node or standalone)",
		Run: func(cmd *cobra.Command, args []string) {
			certificates, _ := cmd.Flags().GetString("certificates")
			address, _ := cmd.Flags().GetString("address")
			nodeURL, _ := cmd.Flags().GetString("node")

			if nodeURL != "" {
				// Send to running node
				addBlockToNode(nodeURL, certificates, address)
			} else {
				// Standalone operation
				addBlockStandalone(certificates, address)
			}
		},
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Get blockchain status",
		Long:  "Get blockchain status (works with running node or standalone)",
		Run: func(cmd *cobra.Command, args []string) {
			nodeURL, _ := cmd.Flags().GetString("node")

			if nodeURL != "" {
				// Get from running node
				getStatusFromNode(nodeURL)
			} else {
				// Standalone operation
				getStatusStandalone()
			}
		},
	}

	var blocksCmd = &cobra.Command{
		Use:   "blocks",
		Short: "List blocks",
		Long:  "List blocks (works with running node or standalone)",
		Run: func(cmd *cobra.Command, args []string) {
			nodeURL, _ := cmd.Flags().GetString("node")

			if nodeURL != "" {
				// Get from running node
				getBlocksFromNode(nodeURL)
			} else {
				// Standalone operation
				getBlocksStandalone()
			}
		},
	}

	// Add flags
	addBlockCmd.Flags().StringP("certificates", "c", "", "Comma-separated certificate IDs (required)")
	addBlockCmd.Flags().StringP("address", "a", "", "Signer address (required)")
	addBlockCmd.Flags().StringP("node", "n", "", "Node URL (e.g., http://localhost:8080)")
	addBlockCmd.MarkFlagRequired("certificates")
	addBlockCmd.MarkFlagRequired("address")

	statusCmd.Flags().StringP("node", "n", "", "Node URL (e.g., http://localhost:8080)")
	blocksCmd.Flags().StringP("node", "n", "", "Node URL (e.g., http://localhost:8080)")

	rootCmd.AddCommand(nodeStartCmd)
	rootCmd.AddCommand(createIdentityCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addBlockCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(blocksCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Start initializes and starts the node server
func (n *Node) Start() error {
	fmt.Printf("Starting Veritas Chain node...\n")
	fmt.Printf("Port: %d\n", n.Port)
	fmt.Printf("University: %s\n", n.University)
	if n.Verbose {
		fmt.Println("Verbose mode enabled")
	}

	// Ensure tmp directory exists
	if err := os.MkdirAll("./tmp", 0755); err != nil {
		return fmt.Errorf("failed to create tmp directory: %v", err)
	}

	// Initialize blockchain
	if err := n.initializeBlockchain(); err != nil {
		return fmt.Errorf("failed to initialize blockchain: %v", err)
	}

	// Setup HTTP server
	mux := http.NewServeMux()
	n.setupRoutes(mux)

	n.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", n.Port),
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Node server starting on port %d...\n", n.Port)
		if err := n.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down node...")
	return n.Stop()
}

// Stop gracefully shuts down the node
func (n *Node) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if n.Blockchain != nil {
		n.Blockchain.Close()
	}

	return n.Server.Shutdown(ctx)
}

// initializeBlockchain initializes or loads the blockchain
func (n *Node) initializeBlockchain() error {
	// Initialize university identity system
	if err := blockchain.InitializeUniversityIdentities(); err != nil {
		return err
	}

	// Create university identity if it doesn't exist
	_, err := blockchain.GetUniversityIdentity(n.University)
	if err != nil {
		// Create the identity
		address, err := blockchain.CreateUniversityIdentity(n.University)
		if err != nil {
			return err
		}
		fmt.Printf("Created university identity for %s: %s\n", n.University, address)
	}

	// Initialize or continue blockchain
	if blockchain.DBExists() {
		n.Blockchain = blockchain.ContinueBlockchain()
		fmt.Println("Loaded existing blockchain")
	} else {
		// Get university identity for genesis block
		universityIdentity, err := blockchain.GetUniversityIdentity(n.University)
		if err != nil {
			return err
		}
		n.Blockchain = blockchain.InitBlockchain(universityIdentity)
		fmt.Println("Created new blockchain with genesis block")
	}

	return nil
}

// setupRoutes configures the HTTP routes
func (n *Node) setupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", n.handleRoot)
	mux.HandleFunc("/status", n.handleStatus)
	mux.HandleFunc("/blocks", n.handleBlocks)
	mux.HandleFunc("/add-block", n.handleAddBlock)
	mux.HandleFunc("/verify", n.handleVerify)
	mux.HandleFunc("/identities", n.handleIdentities)
}

// HTTP handlers
func (n *Node) handleRoot(w http.ResponseWriter, r *http.Request) {
	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"node":      n.University,
			"port":      n.Port,
			"endpoints": []string{"/status", "/blocks", "/add-block", "/verify", "/identities"},
		},
	}
	n.writeJSON(w, response)
}

func (n *Node) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Count blocks
	iter := n.Blockchain.Iterator()
	blockCount := 0
	for {
		block := iter.Next()
		blockCount++
		if len(block.PrevHash) == 0 { // reached genesis
			break
		}
	}

	// Validate chain
	chainValid := true
	validationError := ""
	if err := n.Blockchain.ValidateChain(); err != nil {
		chainValid = false
		validationError = err.Error()
	}

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"node":           n.University,
			"last_hash":      fmt.Sprintf("%x", n.Blockchain.LastHash),
			"block_count":    blockCount,
			"chain_valid":    chainValid,
			"validation_err": validationError,
		},
	}
	n.writeJSON(w, response)
}

func (n *Node) handleBlocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var blocks []map[string]interface{}
	iter := n.Blockchain.Iterator()
	blockCount := 0
	for {
		block := iter.Next()
		blockData := map[string]interface{}{
			"height":           block.Height,
			"hash":             fmt.Sprintf("%x", block.Hash),
			"previous_hash":    fmt.Sprintf("%x", block.PrevHash),
			"timestamp":        block.Timestamp,
			"certificate_count": block.GetCertificateCount(),
			"signature":        fmt.Sprintf("%x", block.Signature),
		}
		blocks = append(blocks, blockData)
		blockCount++
		if len(block.PrevHash) == 0 { // reached genesis
			break
		}
	}

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"blocks": blocks,
			"count":  blockCount,
		},
	}
	n.writeJSON(w, response)
}

func (n *Node) handleAddBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AddBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := APIResponse{
			Success: false,
			Error:   "Invalid JSON request",
		}
		n.writeJSON(w, response)
		return
	}

	// Get identity by address
	universityIdentity, err := blockchain.GetIdentityByAddress(req.Address)
	if err != nil {
		response := APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to get identity for address %s: %v", req.Address, err),
		}
		n.writeJSON(w, response)
		return
	}

	// Add the block
	block := n.Blockchain.AddBlock(req.Certificates, universityIdentity)

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"height":           block.Height,
			"hash":             fmt.Sprintf("%x", block.Hash),
			"signature":        fmt.Sprintf("%x", block.Signature),
			"certificate_count": len(req.Certificates),
			"signed_by":        req.Address,
		},
	}
	n.writeJSON(w, response)
}

func (n *Node) handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := n.Blockchain.ValidateChain()
	response := APIResponse{
		Success: err == nil,
		Data: map[string]interface{}{
			"valid": err == nil,
		},
	}
	if err != nil {
		response.Error = err.Error()
	}
	n.writeJSON(w, response)
}

func (n *Node) handleIdentities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var identities []map[string]interface{}
	for name, address := range blockchain.AuthorizedUniversities {
		if address != "" {
			identity, err := blockchain.GetUniversityIdentity(name)
			if err == nil {
				identityData := map[string]interface{}{
					"name":       name,
					"address":    address,
					"public_key": map[string]interface{}{
						"x": fmt.Sprintf("%x", identity.PrivateKey.PublicKey.X.Bytes()),
						"y": fmt.Sprintf("%x", identity.PrivateKey.PublicKey.Y.Bytes()),
					},
				}
				identities = append(identities, identityData)
			}
		}
	}

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"identities": identities,
			"count":      len(identities),
		},
	}
	n.writeJSON(w, response)
}

// writeJSON writes a JSON response
func (n *Node) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// CLI functions for interacting with nodes or standalone operations

// addBlockToNode sends a block to a running node
func addBlockToNode(nodeURL, certificates, address string) {
	certList := strings.Split(certificates, ",")
	for i, cert := range certList {
		certList[i] = strings.TrimSpace(cert)
	}

	req := AddBlockRequest{
		Certificates: certList,
		Address:      address,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return
	}

	resp, err := http.Post(nodeURL+"/add-block", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if apiResp.Success {
		fmt.Printf("✅ Block added successfully!\n")
		fmt.Printf("Data: %+v\n", apiResp.Data)
	} else {
		fmt.Printf("❌ Failed to add block: %s\n", apiResp.Error)
	}
}

// addBlockStandalone adds a block without a running node
func addBlockStandalone(certificates, address string) {
	// Ensure tmp directory exists
	if err := os.MkdirAll("./tmp", 0755); err != nil {
		fmt.Printf("Failed to create tmp directory: %v\n", err)
		return
	}

	// Check if blockchain exists
	if !blockchain.DBExists() {
		fmt.Println("Blockchain not found. Please run 'veritas start' first to initialize.")
		return
	}

	// Initialize university identity system
	if err := blockchain.InitializeUniversityIdentities(); err != nil {
		fmt.Printf("Failed to initialize university identities: %v\n", err)
		return
	}

	// Get identity by address
	universityIdentity, err := blockchain.GetIdentityByAddress(address)
	if err != nil {
		fmt.Printf("Failed to get identity for address %s: %v\n", address, err)
		return
	}

	// Continue existing blockchain
	bc := blockchain.ContinueBlockchain()
	defer bc.Close()

	// Parse certificate IDs
	certIDs := strings.Split(certificates, ",")
	for i, cert := range certIDs {
		certIDs[i] = strings.TrimSpace(cert)
	}

	// Add the block
	block := bc.AddBlock(certIDs, universityIdentity)
	fmt.Printf("✅ Block created: Height=%d, Hash=%x\n", block.Height, bc.LastHash)
	fmt.Printf("Block signature: %x\n", block.Signature)
	fmt.Printf("Signed by address: %s\n", address)

	// Verify the signature
	if block.Verify(universityIdentity.PrivateKey.PublicKey) {
		fmt.Println("✅ Block signature verification: VALID")
	} else {
		fmt.Println("❌ Block signature verification: INVALID")
	}
}

// getStatusFromNode gets status from a running node
func getStatusFromNode(nodeURL string) {
	resp, err := http.Get(nodeURL + "/status")
	if err != nil {
		fmt.Printf("Error getting status: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if apiResp.Success {
		fmt.Printf("📊 Blockchain Status:\n")
		fmt.Printf("Data: %+v\n", apiResp.Data)
	} else {
		fmt.Printf("❌ Failed to get status: %s\n", apiResp.Error)
	}
}

// getStatusStandalone gets status without a running node
func getStatusStandalone() {
	// Check if blockchain exists
	if !blockchain.DBExists() {
		fmt.Println("❌ Blockchain not found. Please run 'veritas start' first to initialize.")
		return
	}

	// Continue existing blockchain
	bc := blockchain.ContinueBlockchain()
	defer bc.Close()

	// Count blocks
	iter := bc.Iterator()
	blockCount := 0
	for {
		block := iter.Next()
		blockCount++
		if len(block.PrevHash) == 0 { // reached genesis
			break
		}
	}

	// Validate chain
	chainValid := true
	validationError := ""
	if err := bc.ValidateChain(); err != nil {
		chainValid = false
		validationError = err.Error()
	}

	fmt.Printf("📊 Blockchain Status:\n")
	fmt.Printf("  Last Hash: %x\n", bc.LastHash)
	fmt.Printf("  Block Count: %d\n", blockCount)
	fmt.Printf("  Chain Valid: %t\n", chainValid)
	if !chainValid {
		fmt.Printf("  Validation Error: %s\n", validationError)
	}
}

// getBlocksFromNode gets blocks from a running node
func getBlocksFromNode(nodeURL string) {
	resp, err := http.Get(nodeURL + "/blocks")
	if err != nil {
		fmt.Printf("Error getting blocks: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	if apiResp.Success {
		fmt.Printf("📦 Blocks:\n")
		fmt.Printf("Data: %+v\n", apiResp.Data)
	} else {
		fmt.Printf("❌ Failed to get blocks: %s\n", apiResp.Error)
	}
}

// getBlocksStandalone gets blocks without a running node
func getBlocksStandalone() {
	// Check if blockchain exists
	if !blockchain.DBExists() {
		fmt.Println("❌ Blockchain not found. Please run 'veritas start' first to initialize.")
		return
	}

	// Continue existing blockchain
	bc := blockchain.ContinueBlockchain()
	defer bc.Close()

	fmt.Printf("📦 All blocks in chain:\n")
	iter := bc.Iterator()
	blockCount := 0
	for {
		block := iter.Next()
		fmt.Printf("Block %d: Height=%d, Hash=%x, Certificates=%d\n",
			blockCount, block.Height, block.Hash, block.GetCertificateCount())

		// Show signature information
		if len(block.Signature) > 0 {
			fmt.Printf("  Signature: %x (length: %d bytes)\n", block.Signature, len(block.Signature))
		} else {
			fmt.Printf("  ⚠️  No signature found\n")
		}

		blockCount++
		if len(block.PrevHash) == 0 { // reached genesis
			break
		}
	}
}
