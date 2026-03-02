# Veritas Chain (In Progress) 🎓⛓️

**A Blockchain-Based Academic Credential Verification System**

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

## Overview

Veritas Chain is a proof-of-authority blockchain system designed specifically for secure academic credential verification among universities. The system ensures data integrity through cryptographic signatures, maintains privacy through hashed certificate IDs, and provides verifiable authenticity through university-signed blocks using ECDSA cryptography.

The project features a comprehensive CLI interface for node management, identity generation, and interactive blockchain operations.

*Inspired by(stole) my boss's vision for blockchain-based credential verification - thank you!*

## Key Features

- 🔐 **Privacy-First Design**: Certificate data is cryptographically hashed before blockchain storage
- 🏛️ **Proof-of-Authority**: Only authorized universities can create and sign blocks
- 🔑 **ECDSA Cryptography**: Blocks signed with university private keys using P-256 elliptic curves
- 🤝 **University Identity System**: Authorized universities (Harvard, MIT, Stanford, Yale) with persistent identities
- ⚡ **Fast Verification**: Instant block and signature verification
- 🔍 **Tamper-Proof Records**: Immutable blocks with cryptographic signatures
- 💾 **Persistent Storage**: BadgerDB for efficient blockchain data persistence
- 🛡️ **Signature Validation**: Real-time verification of block signatures
- 🖥️ **CLI Interface**: Comprehensive command-line interface for node management and operations
- 🔧 **Interactive Mode**: Real-time blockchain interaction through terminal commands
- 📊 **Merkle Trees**: Efficient certificate verification with Merkle proof generation
- 🔄 **Chain Validation**: Comprehensive blockchain integrity validation

## Quick Start

### Prerequisites

- Go 1.25+ installed
- Git for version control
- Basic understanding of blockchain concepts

### Installation

```bash
# Clone the repository
git clone https://github.com/amanechibana/veritas-chain.git
cd veritas-chain

# Install dependencies
go mod tidy

# Build the CLI tool
go build -o veritas .
```

### Basic Usage

#### 1. Generate a Signer Key

```bash
# Generate a new P-256 private key for your node
./veritas identity keygen

# Output example:
# Generated signer key:
#   SIGNER_PRIVATE_KEY_HEX=1234567890abcdef...
#   Address=1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU
```

#### 2. Set Environment Variables

```bash
# Set your private key (use the output from keygen)
export SIGNER_PRIVATE_KEY_HEX=1234567890abcdef...

# Optional: Set authorized signers file
# cp authorized_signers.json.example authorized_signers.json
```

#### 3. Start Interactive Node

```bash
# Start the node in interactive mode
./veritas node interactive

# This will:
# 1. Load your signer from environment
# 2. Initialize or load existing blockchain
# 3. Start interactive terminal for blockchain operations
```

#### 4. Interactive Commands

```bash
veritas> help                    # Show available commands
veritas> add CERT-001,CERT-002  # Add certificates to new block
veritas> list                    # List all blocks
veritas> validate                # Validate blockchain integrity
veritas> stats                   # Show blockchain statistics
veritas> exit                    # Exit interactive mode
```

## CLI Commands

### Identity Management

```bash
# Generate a new P-256 private key
./veritas identity keygen

# Output:
# Generated signer key:
#   SIGNER_PRIVATE_KEY_HEX=1234567890abcdef...
#   Address=1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU
```

### Node Management

```bash
# Start node in interactive mode
./veritas node interactive

# Global flags available for all commands:
./veritas --verbose --config /path/to/config.yaml node interactive
```

### Interactive Commands

Once in interactive mode (`./veritas node interactive`), you can use:

```bash
# Add certificates to a new block
veritas> add CERT-001,CERT-002,CERT-003

# List all blocks in the chain
veritas> list

# Validate the entire blockchain
veritas> validate

# Show blockchain statistics
veritas> stats

# Show help
veritas> help

# Exit interactive mode
veritas> exit
```

### Environment Variables

```bash
# Required: Your private key (from identity keygen)
export SIGNER_PRIVATE_KEY_HEX=1234567890abcdef...

# Optional: Configuration file path
export VERITAS_CONFIG=/path/to/config.yaml

# Optional: Database path override
export VERITAS_DB_PATH=/custom/db/path
```

### Example Output

#### 1. Generate Signer Key
```bash
$ ./veritas identity keygen
Generated signer key:
  SIGNER_PRIVATE_KEY_HEX=1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
  Address=1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU
```

#### 2. Start Interactive Node
```bash
$ export SIGNER_PRIVATE_KEY_HEX=1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
$ ./veritas node interactive

Configuration:
  Address: 1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU
  DB Path: ./tmp/blocks_1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU
  Resolved Name: harvard
Created new blockchain with genesis block

=== Veritas Chain Interactive Mode ===
Type 'help' for available commands
Signer Address: 1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU
=====================================
```

#### 3. Interactive Commands
```bash
veritas> add CERT-001,CERT-002,CERT-003
Block added successfully!
   Height: 1
   Hash: da0d5b2cf6823599657599319ce3d8cba595fed6bbeab37bfd4799178d4f20eb
   Address: 1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU

veritas> list
Blockchain:
Block 0: Height=0, Hash=6d72c343fd0e485d977d60702cb61b41a70be3820d219c1a39c5da01418fb71c, Address=1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU
Block 1: Height=1, Hash=da0d5b2cf6823599657599319ce3d8cba595fed6bbeab37bfd4799178d4f20eb, Address=1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU

veritas> validate
  Chain validation successful

veritas> stats
Blockchain Statistics:
  Total Blocks: 2
  Total Certificates: 3

veritas> exit
Goodbye!
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Veritas Chain Node                          │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐                             │
│  │    CLI      │  │ Interactive │                             │
│  │ Interface   │  │    Mode     │                             │
│  └──────┬──────┘  └──────┬──────┘                             │
│         │                │                                    │
│         └────────────────┘                                    │
│                          │                                    │
│  ┌───────────────────────▼───────────────────────┐             │
│  │            Blockchain Core                    │             │
│  │  ┌─────────────┐  ┌─────────────┐            │             │
│  │  │   Block     │  │   Merkle    │            │             │
│  │  │ Management  │  │   Trees     │            │             │
│  │  └─────────────┘  └─────────────┘            │             │
│  │  ┌─────────────┐  ┌─────────────┐            │             │
│  │  │   Proof     │  │ Validation  │            │             │
│  │  │ of Authority│  │   Engine    │            │             │
│  │  └─────────────┘  └─────────────┘            │             │
│  └───────────────────────┬───────────────────────┘             │
│                          │                                     │
│  ┌───────────────────────▼───────────────────────┐             │
│  │            Identity System                    │             │
│  │  ┌─────────────┐  ┌─────────────┐            │             │
│  │  │   Signer    │  │   Registry  │            │             │
│  │  │ Management  │  │   (JSON)    │            │             │
│  │  └─────────────┘  └─────────────┘            │             │
│  └───────────────────────┬───────────────────────┘             │
│                          │                                     │
│  ┌───────────────────────▼───────────────────────┐             │
│  │            Storage Layer                      │             │
│  │  ┌─────────────┐  ┌─────────────┐            │             │
│  │  │   BadgerDB  │  │   File      │            │             │
│  │  │  (Blocks)   │  │  System     │            │             │
│  │  └─────────────┘  └─────────────┘            │             │
│  └───────────────────────────────────────────────┘             │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    Future Multi-Node Network                   │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │   Harvard   │    │     MIT     │    │  Stanford   │         │
│  │   Node      │◄──►│    Node     │◄──►│    Node     │         │
│  │             │    │             │    │             │         │
│  │ ┌─────────┐ │    │ ┌─────────┐ │    │ ┌─────────┐ │         │
│  │ │ gRPC    │ │    │ │ gRPC    │ │    │ │ gRPC    │ │         │
│  │ │ Server  │ │    │ │ Server  │ │    │ │ Server  │ │         │
│  │ └─────────┘ │    │ └─────────┘ │    │ └─────────┘ │         │
│  │ ┌─────────┐ │    │ ┌─────────┐ │    │ ┌─────────┐ │         │
│  │ │ gRPC    │ │    │ │ gRPC    │ │    │ │ gRPC    │ │         │
│  │ │ Client  │ │    │ │ Client  │ │    │ │ Client  │ │         │
│  │ └─────────┘ │    │ └─────────┘ │    │ └─────────┘ │         │
│  └─────────────┘    └─────────────┘    └─────────────┘         │
└─────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
veritas-chain/
├── blockchain/          # Core blockchain implementation
│   ├── block.go        # Block structure and operations
│   ├── blockchain.go   # Blockchain management and validation
│   ├── merkle.go       # Merkle tree implementation
│   └── proof.go        # Proof-of-authority consensus
├── cmd/                # CLI command implementations
│   ├── root.go         # Root command and global flags
│   ├── node.go         # Node management commands
│   └── identity.go     # Identity and key management
├── identity/           # University identity system
│   ├── identity.go     # Identity structure and cryptography
│   ├── registry.go     # Identity registry management
│   ├── signer.go       # Signer interface and implementations
│   └── utils.go        # Serialization and utility functions
├── test/               # Test utilities and examples
│   └── test_client.go  # Test client utilities
├── tmp/                # Runtime data storage
│   ├── blocks_*/       # BadgerDB blockchain data (per signer)
├── authorized_signers.json # Authorized university mappings
├── main.go             # Main application entry point
├── go.mod              # Go module dependencies
└── README.md           # This file
```

### Blockchain Features
- **Proof-of-Authority Consensus**: Only authorized universities can create blocks
- **ECDSA P-256 Signatures**: 64-byte signatures for block authentication
- **BadgerDB Storage**: High-performance key-value database for blockchain data
- **Chain Validation**: Comprehensive validation of block integrity and signatures
- **Identity Persistence**: JSON-based university identity storage

### Security Features
- **Cryptographic Hashing**: SHA-256 for certificate ID hashing
- **Digital Signatures**: ECDSA for block authentication
- **Private Key Security**: In-memory private key management
- **Signature Verification**: Real-time validation of block signatures
- **Tamper Detection**: Any modification breaks signature validation

### Current Capabilities
- Create and manage university identities
- Generate genesis blocks with university signatures
- Add blocks with certificate data
- Verify block signatures in real-time
- Persist blockchain data across sessions
- Validate entire blockchain integrity

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package tests
go test ./blockchain
go test ./identity
go test ./server
```

### Building

```bash
# Build the CLI application
go build -o veritas .

# Run the built binary
./veritas --help

# Build with specific Go version
GO_VERSION=1.25 go build -o veritas .
```

### Generating Protocol Buffers

```bash
# Install protoc and Go plugins (if not already installed)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate Go code from proto files
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/blockchain.proto
```

### Dependencies

- **BadgerDB**: High-performance key-value database for blockchain storage
- **Base58**: Address encoding/decoding for university identities
- **Cobra**: CLI framework for command-line interface
- **gRPC**: High-performance RPC framework for API services
- **Protocol Buffers**: Language-neutral data serialization
- **Go Standard Library**: Cryptography and utilities

## Current State vs Final Vision

### Current Implementation
**CLI-Enabled Single-Node Proof-of-Authority Blockchain**
- ✅ **CLI Interface**: Full command-line interface with interactive mode
- ✅ **Identity Management**: P-256 key generation and signer management
- ✅ **Merkle Trees**: Efficient certificate verification with proof generation
- ✅ **Persistent Storage**: BadgerDB with per-signer database isolation
- ✅ **Chain Validation**: Comprehensive blockchain integrity validation
- ⚠️ **Single Node**: Runs on a single machine (no networking yet)
- ⚠️ **Local Identities**: University identities managed locally

### Final Goal: Consortium Blockchain
**Multi-Node Distributed Network**
- 🔄 **P2P Networking**: Each university runs their own node with network communication
- 🔄 **Distributed Consensus**: Harvard, MIT, Stanford, Yale(For example) each host independent nodes
- 🔄 **Block Propagation**: gRPC-based block sharing across university network
- 🔄 **True Distribution**: Each university maintains their own copy of the blockchain
- 🔄 **Network Discovery**: Automatic peer discovery and connection management

## TODO / Future Enhancements

### Phase 1: Core Infrastructure ✅ COMPLETED
- ✅ **Merkle Trees**: Implement Merkle trees for efficient certificate verification
- ✅ **CLI Interface**: Command-line interface for blockchain operations
- ✅ **Identity Management**: P-256 key generation and signer management
- 🔄 **Network Layer**: Add P2P networking for multi-university consensus
- 🔄 **Node Discovery**: Implement network discovery and peer management
- 🔄 **Block Propagation**: Share blocks across university nodes

### Phase 2: Consortium Features
- 🔄 **gRPC API**: Programmatic access through gRPC services for multi-node networking
- 🔄 **Distributed Consensus**: Multi-node agreement on blockchain state
- 🔄 **Validator Rotation**: Universities take turns creating blocks
- 🔄 **Network Resilience**: Handle node failures and network partitions
- 🔄 **Cross-University Validation**: Verify blocks from other universities
- 🔄 **Peer-to-Peer Communication**: gRPC-based node-to-node communication

### Phase 3: Production Features
- 🔄 **Web Interface**: Web-based dashboard for university management or public certificate verification access
- 🔄 **Certificate Standards**: Support for standard academic credential formats (W3C Verifiable Credentials, etc.)
- 🔄 **Monitoring & Metrics**: Blockchain health monitoring and performance metrics
- 🔄 **Backup & Recovery**: Automated backup and disaster recovery procedures

## Security Considerations

- **Private Key Management**: Private keys are managed through environment variables (`SIGNER_PRIVATE_KEY_HEX`) - consider using secure key management systems (HSM, key vaults) for production
- **Database Security**: BadgerDB files are stored locally - ensure proper file permissions and consider encryption at rest
- **Network Security**: Future network layer will require TLS encryption for gRPC communication
- **Access Control**: Role-based access control for university operations (planned for multi-node setup)
- **Audit Logging**: Comprehensive logging of all blockchain operations (enhanced logging planned)
- **Signature Validation**: All blocks are cryptographically signed and validated in real-time
- **Chain Integrity**: Comprehensive blockchain validation ensures tamper detection

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired/Based on @TensorProgramming's guide to blockchain with golang on youtube

---

**Veritas Chain** - *Truth in Education, Secured by Blockchain*
