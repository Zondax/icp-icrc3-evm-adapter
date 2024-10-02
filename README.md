[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GithubActions](https://github.com/Zondax/icp-icrc3-evm-adapter/actions/workflows/checks.golem.yml/badge.svg)](https://github.com/Zondax/icp-icrc3-evm-adapter/blob/master/.github/workflows/checks.golem.yml)

![zondax_light](docs/assets/zondax_light.png#gh-light-mode-only)
![zondax_dark](docs/assets/zondax_dark.png#gh-dark-mode-only)

_Please visit our website at [zondax.ch](https://www.zondax.ch)_

---

# ICP-EVM Proxy

This project aims to enhance the interoperability of Internet Computer Protocol (ICP) with EVM-compatible platforms.

## Components

1. **Counter Canister**: A simple Rust-based canister that demonstrates basic functionality.
2. **Logger Canister**: A Rust-based canister implementing ICRC-3 compatible logging.
3. **EVM Adapter Proxy**: A Go-based service that:
   - Translates ICRC-3 log data into EVM-compatible formats.
   - Exposes EVM RPC compatible methods for interaction with ICP canisters.
4. **SubQuery Indexer**: Indexes and provides queryable access to the translated EVM-compatible data.

## Prerequisites

- [DFINITY Canister SDK](https://sdk.dfinity.org/)
- Docker and Docker Compose
- Node.js (version 16 or later)
- Go (version 1.22 or later)
- Rust (latest stable version)

## Setup and Deployment

1. Clone the repository:

   ```
   git clone https://github.com/your-username/icp-evm-proxy.git
   cd icp-evm-proxy
   ```

2. Deploy Logger Canister:

   ```
   cd canisters/logger_canister
   dfx deploy --network=<chosen_icp_network>
   ```

   Note down the deployed canister ID.

3. Configure Counter Canister:
   Edit `canisters/counter_canister/src/lib.rs` and update the `LOGGER_CANISTER_ID` constant with the ID obtained in the previous step.

4. Deploy Counter Canister:

   ```
   cd ../counter_canister
   dfx deploy --network=<chosen_icp_network>
   ```

5. Configure and deploy EVM Adapter Proxy:

   ```
   cd ../../evm-adapter-proxy
   ```

   Edit `config.yaml` and update the ICP configuration with the deployed canister IDs.

   ```
   make run
   ```

6. Configure and run SubQuery Indexer:

   ```
   cd ../subq-indexer
   ```

   Edit `project.yaml` and update the URL of the deployed EVM Adapter Proxy.

   ```
   npm install
   npm run codegen
   npm run build
   docker-compose up
   ```

## Usage

1. Interact with the Counter Canister to generate log entries.
2. The Logger Canister automatically records these interactions in ICRC-3 format.
3. The EVM Adapter Proxy translates ICRC-3 logs to EVM-compatible formats and exposes EVM RPC methods.
4. Query the translated data using one of the following methods:
   a. Use the SubQuery Indexer's GraphQL endpoint (typically at `http://localhost:3000`) to execute GraphQL queries.
   b. Directly query the PostgreSQL database.

## Stopping and Cleaning Up

To stop all services:

```
dfx stop  # Stops local canisters
docker-compose down  # Stops indexer and database
```

To clean up build artifacts and Docker volumes:

```
dfx clean
docker-compose down -v
```
