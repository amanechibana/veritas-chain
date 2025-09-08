# Veritas Chain 🎓⛓️

**A Blockchain-Based Academic Credential Verification System**

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)]()

## Overview

Veritas Chain is currently a proof-of-authority blockchain system designed specifically for secure academic credential verification among universities. The system ensures data integrity through cryptographic signatures, maintains privacy through hashed certificate IDs, and provides verifiable authenticity through university-signed blocks using ECDSA cryptography.

I stole this idea as a project from my boss so thank you :p. I hope this will end up working somewhat...

## Key Features

- 🔐 **Privacy-First Design**: Certificate data is cryptographically hashed before blockchain storage
- 🏛️ **Proof-of-Authority**: Only authorized universities can create and sign blocks
- 🔑 **ECDSA Cryptography**: Blocks signed with university private keys using P-256 elliptic curves
- 🤝 **University Identity System**: Authorized universities (Harvard, MIT, Stanford, Yale) with persistent identities
- ⚡ **Fast Verification**: Instant block and signature verification
- 🔍 **Tamper-Proof Records**: Immutable blocks with cryptographic signatures
- 💾 **Persistent Storage**: BadgerDB for efficient blockchain data persistence
- 🛡️ **Signature Validation**: Real-time verification of block signatures

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

```

### Basic Usage

```bash
# Run the blockchain demo
go run .

# The program will:
# 1. Initialize university identity system
# 2. Create a Harvard university identity
# 3. Create genesis block signed by Harvard
# 4. Add blocks with certificate data
# 5. Verify all block signatures
# 6. Display comprehensive blockchain information
```

### Example Output

```
Identities saved successfully
University identity created for harvard: Address=1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU
--------------------------------
University Identity Details:
  Name: harvard
  Address: 1H8vrviwK5Ep83sDkP8m8XsYpprVNiB8dU
  Public Key X: 6dc7aa29d12dd27777d882bee69c9efad24c9d665af4f734a3e4bf28c09ff132
  Public Key Y: 4d63f3b2bc9824c4973ccea3fa95744618583e99c01e4147bced8b7f61bc2271
--------------------------------
Genesis block and blockchain created: LastHash=6d72c343fd0e485d977d60702cb61b41a70be3820d219c1a39c5da01418fb71c
--------------------------------
Block 1 created: Height=1, Hash=da0d5b2cf6823599657599319ce3d8cba595fed6bbeab37bfd4799178d4f20eb
Block 1 signature: ed34134f19f8dcd8a32037e83efce1a3e5b35451a46178b8238d80af28773bebc5a6e5e94ed9fae01d3f25aaedbbc5ea0ff041f0bc581725cac3b59078b7d0f2
Block 1 signature length: 64 bytes
✅ Block 1 signature verification: VALID
--------------------------------
Chain validation passed!
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Harvard       │    │      MIT        │    │   Stanford      │
│   University    │    │   University    │    │   University    │
│   Identity      │    │   Identity      │    │   Identity      │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────────┐
                    │  Veritas Chain  │
                    │  (BadgerDB)     │
                    │                 │
                    │ ┌─────────────┐ │
                    │ │ Block N     │ │
                    │ │ ├─Timestamp │ │
                    │ │ ├─Hash      │ │
                    │ │ ├─PrevHash  │ │
                    │ │ ├─Height    │ │
                    │ │ ├─CertHashes│ │
                    │ │ └─Signature │ │
                    │ └─────────────┘ │
                    └─────────────────┘
```

## Project Structure

```
veritas-chain/
├── blockchain/          # Core blockchain implementation
│   ├── block.go        # Block structure and operations
│   ├── blockchain.go   # Blockchain management and validation
│   └── proof.go        # Proof-of-authority consensus
├── identity/           # University identity system
│   ├── identity.go     # Identity structure and cryptography
│   ├── identities.go   # Identity management and persistence
│   └── utils.go        # Serialization and utility functions
├── tmp/                # Runtime data storage
│   ├── blocks/         # BadgerDB blockchain data
│   └── identities.data # University identity persistence
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
```

### Building

```bash
# Build the application
go build -o veritas main.go

# Run the built binary
./veritas
```

### Dependencies

- **BadgerDB**: High-performance key-value database
- **Base58**: Address encoding/decoding
- **Go Standard Library**: Cryptography and utilities

## Current State vs Final Vision

### Current Implementation
**Single-Node Proof-of-Authority Blockchain**
- Runs on a single machine
- All university identities managed locally
- No network communication between universities
- Centralized execution with cryptographic authorization

### Final Goal: Consortium Blockchain
**Multi-Node Distributed Network**
- Each university runs their own node
- Harvard, MIT, Stanford, Yale each host independent nodes
- P2P networking for block propagation and consensus
- True distributed consensus across university network
- Each university maintains their own copy of the blockchain

## TODO / Future Enhancements

### Phase 1: Core Infrastructure
- **Merkle Trees**: Implement Merkle trees for efficient certificate verification
- **Network Layer**: Add P2P networking for multi-university consensus
- **Node Discovery**: Implement network discovery and peer management
- **Block Propagation**: Share blocks across university nodes

### Phase 2: Consortium Features
- **Distributed Consensus**: Multi-node agreement on blockchain state
- **Validator Rotation**: Universities take turns creating blocks
- **Network Resilience**: Handle node failures and network partitions
- **Cross-University Validation**: Verify blocks from other universities

### Phase 3: Production Features
- **API Interface**: REST/gRPC APIs for external integration
- **CLI Tools**: Command-line interface for blockchain operations
- **Web Interface**: Web-based dashboard for university management or public certificate verification access 
- **Certificate Standards**: Support for standard academic credential formats

## Security Considerations

- **Private Key Management**: Currently stores private keys in JSON (should be moved to secure storage)
- **Network Security**: Future network layer will require TLS encryption
- **Access Control**: Role-based access control for university operations
- **Audit Logging**: Comprehensive logging of all blockchain operations

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by @TensorProgramming's guide to blockchain with golang on youtube

---

**Veritas Chain** - *Truth in Education, Secured by Blockchain*