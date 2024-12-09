---
title: EVM Adapter Proxy
sidebar_position: 2
---

# EVM Adapter Proxy

The EVM Adapter Proxy is a crucial component in the ICP-EVM Proxy project. It serves as a bridge between the Internet Computer Protocol (ICP) and EVM-compatible systems, enabling interoperability and data translation.

## Important Considerations

This component is a proof of concept (PoC) and is not designed for production use. Its main purpose is to demonstrate the feasibility of integration between ICP and EVM-compatible systems.

### DEX Client Integration

The DEX client integration in this proxy is implemented solely for demonstration and testing purposes in this PoC. It provides a simplified interface to interact with basic DEX operations.

## Main Functions

1. **Log Retrieval**: Fetches ICRC-3 compatible logs from the Logger Canister on the Internet Computer.
2. **Log Translation**: Converts ICRC-3 logs into EVM-compatible event logs.
3. **EVM RPC Compatibility**: Exposes EVM RPC compatible methods, allowing interaction with ICP canisters using familiar Ethereum tooling.

## Implementation

The EVM Adapter Proxy is implemented in Go and consists of several main components:

- A structure to hold necessary configuration and state.
- Functions to fetch logs from the Logger Canister.
- Functions to translate ICRC-3 logs to EVM-compatible logs.
- An HTTP server to handle incoming RPC requests.

## Log Translation Process

The log translation process involves mapping ICRC-3 log fields to their EVM equivalents. This includes:

- Converting canister IDs to Ethereum-style addresses.
- Mapping event types to appropriate topic hashes.
- Formatting log data to match EVM expectations.

## EVM RPC Compatibility

The EVM Adapter Proxy implements a subset of Ethereum JSON-RPC methods to provide compatibility with Ethereum tooling. Some key methods include:

### Standard Ethereum Methods

- `eth_chainId`: Returns the chain ID.
- `eth_blockNumber`: Returns the number of the most recent block.
- `eth_getBlockByNumber`: Retrieves a block by its number.
- `eth_getBlockByHash`: Retrieves a block by its hash.
- `eth_getLogs`: Retrieves logs that match the specified filter criteria.
- `eth_accounts`: Returns a list of account addresses.
- `net_version`: Returns the current network version.
- `net_listening`: Indicates whether the client is actively listening for network connections.
- `net_peerCount`: Returns the number of peers currently connected to the client.
- `web3_clientVersion`: Returns the current client version.
- `web3_sha3`: Returns the Keccak-256 hash of the given input.

### DEX-Specific Methods

- `eth_getCurrencyPairs`: Returns all available currency pairs from the DEX.
- `eth_mintTokens`: Creates new tokens for a specified recipient.
- `eth_burnTokens`: Destroys tokens from an owner's balance.

These methods allow Ethereum tools and libraries to interact with ICP canisters as if they were EVM-compatible smart contracts.

## Example Usage

### Standard Methods

```bash
# Get Chain ID
curl -X POST -H "Content-Type: application/json" \
--data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
http://localhost:3030/rpc/v1

# Get Latest Block
curl -X POST -H "Content-Type: application/json" \
--data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
http://localhost:3030/rpc/v1
```

### DEX Methods

```bash
# Get Currency Pairs
curl -X POST -H "Content-Type: application/json" \
--data '{"jsonrpc":"2.0","method":"eth_getCurrencyPairs","params":[],"id":1}' \
http://localhost:3030/rpc/v1

# Mint Tokens
curl -X POST -H "Content-Type: application/json" \
--data '{
  "jsonrpc":"2.0",
  "method":"eth_mintTokens",
  "params":[{
    "currency": "ICP",
    "amount": "0x5f5e100",
    "recipient": "0x2vxsx-fae"
  }],
  "id":1
}' \
http://localhost:3030/rpc/v1

# Burn Tokens
curl -X POST -H "Content-Type: application/json" \
--data '{
  "jsonrpc":"2.0",
  "method":"eth_burnTokens",
  "params":[{
    "currency": "ICP",
    "amount": "0x2faf080",
    "owner": "0x2vxsx-fae"
  }],
  "id":1
}' \
http://localhost:3030/rpc/v1
```

Note: For this PoC, Ethereum addresses should be formatted as "0x<full_icp_principal>".
For example, if your ICP principal is "2vxsx-fae", use "0x2vxsx-fae" as the Ethereum address.
The proxy will simply remove the "0x" prefix to get the ICP principal.

## Integration with Other Components

- **Logger Canister**: The EVM Adapter Proxy retrieves logs from the Logger Canister on the Internet Computer.
- **DEX Canister**: Provides access to DEX functionality through EVM-compatible methods.
- **SubQuery Indexer**: Translated logs are made available to the SubQuery Indexer for further processing and indexing.

For more details on the interaction between components, refer to the project overview and other component-specific documentation.

## Technical Details

### Client Generation

The proxy uses automatically generated clients for both Logger and DEX canisters. These clients are generated using the Makefile command:

```bash
# Generate both Logger and DEX clients
make generate-client
```

This command internally uses the `goic` tool to generate Go clients from the Candid interface definitions. The process is defined in `Makefile.local.mk`

The generated clients provide type-safe interfaces for interacting with the canisters.

### Configuration

The proxy requires configuration for both Logger and DEX canister IDs in `config.yaml`:

```yaml
icp:
  loggerCanisterId: "ydpfi-uiaaa-aaaal-qjupa-cai"
  dexCanisterId: "7eo5f-eqaaa-aaaam-adqoq-cai"
  nodeUrl: "https://ic0.app"
```

### Development Workflow

1. After any changes to canister interfaces:

   ```bash
   make generate-client  # Regenerate Go clients
   ```

2. Build and run the proxy:

   ```bash
   make build
   ./output/poc-icp-icrc3-evm-adapter start -c config.yaml
   ```
