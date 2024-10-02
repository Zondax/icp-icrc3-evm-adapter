# EVM Adapter Proxy

The EVM Adapter Proxy is a crucial component in the ICP-EVM Proxy project. It serves as a bridge between the Internet Computer Protocol (ICP) and EVM-compatible systems, enabling interoperability and data translation.

## Important Considerations

This component is a proof of concept (PoC) and is not designed for production use. Its main purpose is to demonstrate the feasibility of integration between ICP and EVM-compatible systems.

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

These methods allow Ethereum tools and libraries to interact with ICP canisters as if they were EVM-compatible smart contracts.

## Integration with Other Components

- **Logger Canister**: The EVM Adapter Proxy retrieves logs from the Logger Canister on the Internet Computer.
- **SubQuery Indexer**: Translated logs are made available to the SubQuery Indexer for further processing and indexing.

For more details on the interaction between components, refer to the project overview and other component-specific documentation.
