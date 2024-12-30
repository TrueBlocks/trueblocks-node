
# Technical Specification for TrueBlocks Node

## Title Page

**Technical Specification for TrueBlocks Node**  
Version: 1.0  
Published: [Insert Date]  
Author: [Your Company/Name]  

---

## Table of Contents

1. [Introduction](./TechnicalSpec_1.md)
2. [System Architecture](./TechnicalSpec_2.md)
3. [Core Functionalities](./TechnicalSpec_3.md)
4. [Technical Design](./TechnicalSpec_4.md)
5. [Supported Chains and RPCs](./TechnicalSpec_5.md)
6. [Command-Line Interface](./TechnicalSpec_6.md)
7. [Performance and Scalability](./TechnicalSpec_7.md)
8. [Integration Points](./TechnicalSpec_8.md)
9. [Testing and Validation](./TechnicalSpec_9.md)
10. [Appendices](./TechnicalSpec_10.md)


# Section 1: Introduction

## Purpose of the Technical Specification

This document defines the technical architecture, design, and functionalities of TrueBlocks Node, enabling developers and engineers to understand its internal workings and design principles.

## Intended Audience

This specification is for:
- Developers working on TrueBlocks Node or integrating it into applications.
- System architects designing systems that use TrueBlocks Node.
- Technical professionals looking for a detailed understanding of the system.

## Scope and Objectives

The specification covers:
- High-level architecture.
- Core functionalities such as blockchain indexing, REST API, and address monitoring.
- Design principles, including scalability, error handling, and integration with IPFS.
- Supported chains, RPC requirements, and testing methodologies.

# Section 2: System Architecture

## High-Level Architecture Diagram

(Include a diagram here if needed. Replace this text with a Markdown-compatible diagram or a link to an image.)

## Key Components Overview

1. **Blockchain Indexer**: Handles blockchain data collection and indexing.
2. **REST API Server**: Exposes APIs for data access.
3. **IPFS Integrator**: Manages decentralized storage.
4. **Configuration Manager**: Parses `.env` files and other configurations.

## Interactions Between Components

- The Blockchain Indexer collects data from RPC endpoints and stores it in the local database.
- The REST API retrieves indexed data and exposes it via endpoints.
- The IPFS Integrator uploads and pins indexed data to IPFS for decentralized access.

# Section 3: Core Functionalities

## Blockchain Indexing

Indexes blockchain data for fast and efficient retrieval. Supports multiple chains and tracks transactions.

## REST API

Exposes indexed data through a REST API. Includes endpoints for:
- Retrieving transactions and blocks.
- Accessing monitored address data.

## Address Monitoring

Allows tracking of specific blockchain addresses. Captures transactions and updates in real-time.

## IPFS Integration

Pins portions of the Unchained Index to IPFS for decentralized and tamper-proof storage.

# Section 4: Technical Design

## Configuration Files and Environment Variables

TrueBlocks Node uses a `.env` file for configuration. Key variables include:
- `TB_NODE_DATADIR`: Directory for storing data.
- `TB_NODE_MAINNETRPC`: RPC endpoint for Ethereum mainnet.
- `TB_NODE_CHAINS`: List of chains to index.

## Initialization Process

1. Validate `.env` configuration.
2. Connect to RPC endpoints for the specified chains.
3. Initialize the blockchain index if necessary.

## Data Flow and Processing

- **Input**: Blockchain data retrieved via RPC.
- **Processing**: Indexing, storing, and optionally pinning data to IPFS.
- **Output**: Indexed data accessible through the REST API.

## Error Handling and Logging

Logs are written to the console with adjustable levels (`Debug`, `Info`, `Warn`, `Error`). Errors during initialization or RPC interactions are logged and reported.

# Section 5: Supported Chains and RPCs

## List of Supported Blockchains

TrueBlocks Node supports Ethereum mainnet and other EVM-compatible chains like:
- Sepolia
- Gnosis
- Optimism

## Requirements for RPC Endpoints

Each chain requires a valid RPC endpoint. For example:
- `TB_NODE_MAINNETRPC`: Mainnet RPC URL.
- `TB_NODE_SEPOLIARPC`: Sepolia RPC URL.

## Handling Multiple Chains

To enable multiple chains, set `TB_NODE_CHAINS` in the `.env` file:
```env
TB_NODE_CHAINS="mainnet,sepolia,gnosis"
```
Ensure each chain has a corresponding RPC endpoint.

# Section 6: Command-Line Interface

## Available Commands and Options

### Initialization
```bash
./trueblocks-node --init all
```
- Options: `all`, `blooms`, `none`

### Scraper
```bash
./trueblocks-node --scrape on
```
- Enables or disables the blockchain scraper.

### REST API
```bash
./trueblocks-node --api on
```
- Starts the API server.

### Sleep Duration
```bash
./trueblocks-node --sleep 60
```
- Sets the duration (in seconds) between updates.

## Detailed Behavior for Each Command

1. **`--init`**: Controls how the blockchain index is initialized.
2. **`--scrape`**: Toggles the blockchain scraper.
3. **`--api`**: Starts or stops the API server.

# Section 7: Performance and Scalability

## Performance Benchmarks

TrueBlocks Node is designed to handle high-throughput blockchain data. Typical performance benchmarks include:
- Processing speed: ~500 blocks per second (depending on RPC response time).
- REST API response time: <50ms for standard queries.

## Strategies for Handling Large-Scale Data

1. Use high-performance RPC endpoints with low latency.
2. Increase local storage capacity to handle large blockchain data.
3. Scale horizontally by running multiple instances of TrueBlocks Node for different chains.

## Resource Optimization Guidelines

- Limit the number of chains processed simultaneously to reduce system load.
- Configure `--sleep` duration to balance processing speed with system resource usage.

# Section 8: Integration Points

## Integration with External APIs

TrueBlocks Node exposes data through a REST API, making it compatible with external applications. Example use cases:
- Fetching transaction details for a given address.
- Retrieving block information for analysis.

## Interfacing with IPFS

Data indexed by TrueBlocks Node can be pinned to IPFS for decentralized storage:
```bash
./trueblocks-node --ipfs on
```

## Customizing for Specific Use Cases

Users can tailor the configuration by:
- Adjusting `.env` variables to include specific chains and RPC endpoints.
- Writing custom scripts to query the REST API and process the data.

# Section 9: Testing and Validation

## Unit Testing

Unit tests cover:
- Blockchain indexing logic.
- Configuration parsing and validation.
- REST API endpoint functionality.

Run tests with:
```bash
go test ./...
```

## Integration Testing

Integration tests ensure all components work together as expected. Tests include:
- RPC connectivity validation.
- Multi-chain indexing workflows.

## Testing Guidelines for Developers

1. Use mock RPC endpoints for testing without consuming live resources.
2. Validate `.env` configuration in test environments before deployment.
3. Automate tests with CI/CD pipelines to ensure reliability.

# Section 10: Appendices

## Glossary of Technical Terms

- **EVM**: Ethereum Virtual Machine, the runtime environment for smart contracts.
- **RPC**: Remote Procedure Call, a protocol for interacting with blockchain nodes.
- **IPFS**: InterPlanetary File System, a decentralized storage solution.

## References and Resources

- [TrueBlocks GitHub Repository](https://github.com/TrueBlocks/trueblocks-node)
- [TrueBlocks Official Website](https://trueblocks.io)
- [Ethereum Developer Documentation](https://ethereum.org/en/developers/)
- [IPFS Documentation](https://docs.ipfs.io)

# Index for Technical Specification

- **Address Monitoring**: Section 3, Core Functionalities
- **API Access**: Section 3, Core Functionalities
- **Architecture Overview**: Section 2, System Architecture
- **Blockchain Indexing**: Section 3, Core Functionalities
- **Configuration Files**: Section 4, Technical Design
- **Data Flow**: Section 4, Technical Design
- **Error Handling**: Section 4, Technical Design
- **Integration Points**: Section 8, Integration Points
- **IPFS Integration**: Section 3, Core Functionalities; Section 8, Integration Points
- **Logging**: Section 4, Technical Design
- **Performance Benchmarks**: Section 7, Performance and Scalability
- **REST API**: Section 3, Core Functionalities; Section 8, Integration Points
- **RPC Requirements**: Section 5, Supported Chains and RPCs
- **Scalability Strategies**: Section 7, Performance and Scalability
- **System Components**: Section 2, System Architecture
- **Testing Guidelines**: Section 9, Testing and Validation
