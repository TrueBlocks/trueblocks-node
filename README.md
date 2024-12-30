
# TrueBlocks Node

TrueBlocks Node is a powerful tool for indexing and managing Ethereum and EVM-compatible blockchain data. It enables efficient data retrieval and monitoring by indexing blockchain transactions, providing a REST API, and supporting multiple features like IPFS integration, address monitoring, and more.

## Key Features

- **Blockchain Indexing**: Actively indexes EVM-compatible blockchains based on user configuration.
- **Unchained Index**: Pins portions of the Unchained Index to IPFS and publishes its manifest hash to a smart contract.
- **REST API**: Exposes the full functionality of the TrueBlocks `chifra` toolset via a REST API.
- **Address Monitoring**: Allows tracking transactions for specific addresses.
- **Highly Configurable**: Customizable chains, RPC endpoints, and runtime behavior.

## Documentation

For complete details on using and configuring TrueBlocks Node, refer to the following documents (located in the same directory as this `README.md`):

- [**User's Manual**](./UsersManual.md): Step-by-step instructions on installation, configuration, and basic usage.
- [**Technical Specification**](./TechnicalSpecification.md): Detailed technical documentation for developers and advanced users.

## Installation without verification

```bash
go install github.com/TrueBlocks/trueblocks-node/v4@latest
```

This will install the `trueblocks-node` binary in your `GOBIN` directory which is likely in your `$PATH`.

## Getting Started

### Prerequisites

1. **System Requirements**:
   - Operating system: Linux or macOS
   - Golang: Version 1.19 or later
   - Internet connectivity for blockchain data synchronization

2. **Dependencies**:
   - A valid RPC endpoint for the Ethereum mainnet is mandatory.
   - Additional RPC endpoints for other chains (e.g., Sepolia, Gnosis) as required.

3. **Environment Setup**:
   - Create a `.env` file in the working directory. You can use the provided `env.example` file as a reference.

### Installation

1. **Clone the Repository**:
   
   ```bash
   git clone https://github.com/TrueBlocks/trueblocks-node.git
   cd trueblocks-node
   ```

2. **Build the Application**:
   Ensure that Go is installed, then run:
   
   ```bash
   go build -o trueblocks-node .
   ```

3. **Environment Configuration**:
   Create a `.env` file in the directory where you will run the application. Populate it as follows:
   
   ```env
   # Data directory for storing blockchain data
   TB_NODE_DATADIR="PWD/data"

   # Comma-separated list of chains to index (default: mainnet)
   TB_NODE_CHAINS="mainnet,sepolia"

   # RPC providers for each chain
   TB_NODE_MAINNETRPC="https://mainnet.infura.io/v3/<your-api-key>"
   TB_NODE_SEPOLIARPC="https://sepolia.infura.io/v3/<your-api-key>"

   # Logging level (Options: Debug, Info, Warn, Error)
   TB_LOGLEVEL="Info"
   ```

## Usage

### Running the Application

Execute the binary with the desired configuration:

```bash
./trueblocks-node <options>
```

### Command-Line Options

| Option      | Description                                                          | Default |
| ----------- | -------------------------------------------------------------------- | ------- |
| `--init`    | Initialize the Unchained Index. Options: `all`, `blooms`, or `none`. | `none`  |
| `--scrape`  | Enable or disable the scraper. Options: `on`, `off`.                 | `off`   |
| `--api`     | Enable or disable the REST API. Options: `on`, `off`.                | `off`   |
| `--ipfs`    | Enable or disable IPFS. Options: `on`, `off`.                        | `off`   |
| `--monitor` | Enable or disable address monitoring. Options: `on`, `off`.          | `off`   |
| `--sleep`   | Time in seconds between updates when scraping.                       | `30`    |
| `--version` | Display the current version of TrueBlocks Node.                      |         |
| `--help`    | Display the help text with all available options.                    |         |

### Example Command

To index mainnet and enable the API:

```bash
./trueblocks-node --init all --scrape on --api on
```

## Configuration

### Required Environment Variables

- `TB_NODE_DATADIR`: Path to the data directory where blockchain data is stored.
- `TB_NODE_MAINNETRPC`: A valid Ethereum mainnet RPC endpoint.

### Optional Environment Variables

- `TB_NODE_CHAINS`: Comma-separated list of chains to index. Example: `mainnet,sepolia`.
- `TB_NODE_<CHAIN>RPC`: RPC endpoints for each chain (e.g., `TB_NODE_SEPOLIARPC` for Sepolia).
- `TB_LOGLEVEL`: Logging level for the application.

## Development and Testing

### Testing

Run unit tests using the following command:
```bash
go test ./...
```

### Development Workflow

- All configurations and templates are located in the `config` package.
- The application logic is in the `app` package.
- Testing covers edge cases and expected inputs, particularly for chain validation, configuration parsing, and argument handling.

## Troubleshooting

1. **Missing RPC Provider**:
   Ensure that you have defined `TB_NODE_MAINNETRPC` and any additional required RPC endpoints in your `.env` file.

2. **Configuration Errors**:
   Run the app with `--help` to ensure your command-line arguments are correctly formatted.

3. **Logs**:
   Logs are printed to standard output by default. Adjust the logging level using `TB_LOGLEVEL`.

## License

This project is licensed under the MIT License. See `LICENSE` for details.

## Contributing

We love contributors. Please see information about our [workflow](https://github.com/TrueBlocks/trueblocks-core/blob/develop/docs/BRANCHING.md) before proceeding.

1. Fork this repository into your own repo.
2. Create a branch: `git checkout -b <branch_name>`.
3. Make changes to your local branch and commit them to your forked repo: `git commit -m '<commit_message>'`
4. Push back to the original branch: `git push origin TrueBlocks/trueblocks-core`
5. Create the pull request.

## Contact

If you have questions, comments, or complaints, please join the discussion on our discord server which is [linked from our website](https://trueblocks.io).

## List of Contributors

Thanks to the following people who have contributed to this project:

- [@tjayrush](https://github.com/tjayrush)
- [@dszlachta](https://github.com/dszlachta)
- many others
