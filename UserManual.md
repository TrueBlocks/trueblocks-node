
# TrueBlocks Node User's Manual

## Title Page

**TrueBlocks Node User's Manual**  
Version: 1.0  
Published: [Insert Date]  
Author: [Your Company/Name]  

---

## Table of Contents

1. [Introduction](./UserManual_Chapter_1.md)
2. [Getting Started](./UserManual_Chapter_2.md)
3. [Understanding TrueBlocks Node](./UserManual_Chapter_3.md)
4. [Using TrueBlocks Node](./UserManual_Chapter_4.md)
5. [Advanced Operations](./UserManual_Chapter_5.md)
6. [Maintenance and Troubleshooting](./UserManual_Chapter_6.md)
7. [Appendices](./UserManual_Chapter_7.md)

# Chapter 1: Introduction

## Overview of TrueBlocks Node

TrueBlocks Node is a blockchain indexing and monitoring application designed to provide users with an efficient way to interact with and manage EVM-compatible blockchain data. It supports functionalities like transaction monitoring, blockchain indexing, and a REST API for accessing data.

## Purpose of the User's Manual

This User's Manual is designed to help users get started with TrueBlocks Node, understand its features, and operate the application effectively for both basic and advanced use cases.

## Intended Audience

This manual is for:
- End-users looking to index and monitor blockchain data.
- Developers integrating blockchain data into their applications.
- System administrators managing blockchain-related infrastructure.

# Chapter 2: Getting Started

## System Requirements

To run TrueBlocks Node, ensure your system meets the following requirements:
- Operating System: Linux or macOS
- Golang: Version 1.19 or later
- An active internet connection

## Installation Guide

1. Clone the repository:
   ```bash
   git clone https://github.com/TrueBlocks/trueblocks-node.git
   cd trueblocks-node
   ```

2. Build the application:
   ```bash
   go build -o trueblocks-node .
   ```

3. Configure the environment:
   - Create a `.env` file in the working directory using the provided `env.example` file.
   - Define required environment variables, including `TB_NODE_DATADIR` and `TB_NODE_MAINNETRPC`.

## Initial Configuration

Populate your `.env` file with the necessary parameters:
```env
TB_NODE_DATADIR="PWD/data"
TB_NODE_CHAINS="mainnet,sepolia"
TB_NODE_MAINNETRPC="https://mainnet.infura.io/v3/<your-api-key>"
TB_LOGLEVEL="Info"
```

## Starting the Application

Run the application with:
```bash
./trueblocks-node --init all --scrape on --api on
```

# Chapter 3: Understanding TrueBlocks Node

## Key Features

- **Blockchain Indexing**: Active indexing of EVM-compatible chains.
- **REST API**: Expose blockchain data and `chifra` commands via a RESTful interface.
- **Address Monitoring**: Track specific blockchain addresses for transactions.
- **IPFS Integration**: Pin indexed data to IPFS for decentralized storage.

## Application Interface Overview

TrueBlocks Node operates through:
- **Command-Line Interface (CLI)**: For configuration and command execution.
- **REST API**: For programmatic interaction with indexed data.

## Terminology and Concepts

- **Unchained Index**: An index of blockchain data optimized for querying.
- **Chains**: EVM-compatible blockchains (e.g., Ethereum mainnet, Sepolia).
- **Providers**: RPC endpoints for interacting with blockchains.

# Chapter 4: Using TrueBlocks Node

## Indexing Blockchains

To index a blockchain, ensure the required environment variables are set for your RPC endpoints, then run:
```bash
./trueblocks-node --init all --scrape on
```
This will initialize the blockchain index and start the scraping process.

## Accessing the REST API

Enable the REST API by running the application with:
```bash
./trueblocks-node --api on
```
Access the API through the default endpoint:
```
http://localhost:8080
```
Refer to the API documentation for available endpoints and usage.

## Monitoring Addresses

You can monitor specific blockchain addresses for transactions. Configure the monitored addresses in your `.env` file or through the API, and enable monitoring:
```bash
./trueblocks-node --monitor on
```

## Managing Configurations

TrueBlocks Node configurations can be managed using the `.env` file. Changes to the `.env` file require a restart of the application to take effect. 

# Chapter 5: Advanced Operations

## Integrating with IPFS

Enable IPFS support with:
```bash
./trueblocks-node --ipfs on
```
This will pin indexed blockchain data to IPFS, ensuring decentralized storage and retrieval.

## Customizing Chain Indexing

Specify additional chains by updating the `TB_NODE_CHAINS` environment variable. Example:
```env
TB_NODE_CHAINS="mainnet,sepolia,gnosis"
```
Ensure each chain has a valid RPC endpoint configured.

## Utilizing Command-Line Options

Key options include:
- `--init [all|blooms|none]`: Specify the type of index initialization.
- `--scrape [on|off]`: Enable or disable the scraper.
- `--api [on|off]`: Enable or disable the API.
- `--sleep [int]`: Set the sleep duration between updates in seconds.

# Chapter 6: Maintenance and Troubleshooting

## Updating TrueBlocks Node

To update the application, pull the latest changes from the repository and rebuild the binary:
```bash
git pull
go build -o trueblocks-node .
```

## Common Issues and Solutions

- **Missing RPC Provider**: Ensure your `.env` file contains valid RPC URLs.
- **Configuration Errors**: Use `--help` to validate command-line arguments.

## Log Files and Debugging

Logs are written to the standard output by default. Set the log level in the `.env` file:
```env
TB_LOGLEVEL="Debug"
```

## Contacting Support

If you encounter issues not covered in this guide, contact support at:
[TrueBlocks Support](mailto:support@trueblocks.io)

# Chapter 7: Appendices

## Glossary of Terms

- **EVM**: Ethereum Virtual Machine, the runtime environment for smart contracts in Ethereum and similar blockchains.
- **RPC**: Remote Procedure Call, a protocol allowing the application to communicate with blockchain nodes.
- **Indexing**: The process of organizing blockchain data for fast and efficient retrieval.
- **IPFS**: InterPlanetary File System, a decentralized storage system for sharing and retrieving data.

## Frequently Asked Questions (FAQ)

### 1. What chains are supported by TrueBlocks Node?
TrueBlocks Node supports Ethereum mainnet and other EVM-compatible chains such as Sepolia and Gnosis. Additional chains can be added by configuring the `TB_NODE_CHAINS` environment variable.

### 2. Do I need an RPC endpoint for every chain?
Yes, each chain you want to index or interact with requires a valid RPC endpoint specified in the `.env` file.

### 3. Can I run TrueBlocks Node without IPFS?
Yes, IPFS integration is optional and can be enabled or disabled using the `--ipfs` command-line option.

## References and Further Reading

- [TrueBlocks GitHub Repository](https://github.com/TrueBlocks/trueblocks-node)
- [TrueBlocks Official Website](https://trueblocks.io)
- [Ethereum Developer Documentation](https://ethereum.org/en/developers/)
- [IPFS Documentation](https://docs.ipfs.io)


# Index

- Address Monitoring: Chapter 4, Section "Monitoring Addresses"
- Advanced Operations: Chapter 5
- API Access: Chapter 4, Section "Accessing the REST API"
- Blockchain Indexing: Chapter 4, Section "Indexing Blockchains"
- Chains: Chapter 3, Section "Terminology and Concepts"
- Configuration Management: Chapter 4, Section "Managing Configurations"
- Glossary: Chapter 7, Section "Glossary of Terms"
- IPFS Integration: Chapter 5, Section "Integrating with IPFS"
- Logging and Debugging: Chapter 6, Section "Log Files and Debugging"
- RPC Endpoints: Chapter 2, Section "Initial Configuration"
- Troubleshooting: Chapter 6
