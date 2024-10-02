# Setup and Deployment

This guide provides instructions for setting up and deploying the ICP-EVM Proxy project components.

## Prerequisites

- [DFINITY Canister SDK](https://sdk.dfinity.org/)
- Docker and Docker Compose
- Node.js (version 16 or later)
- Go (version 1.22 or later)
- Rust (latest stable version)

## Deployment Steps

1. **Clone the Repository**

   ```shell
   git clone https://github.com/your-username/icp-evm-proxy.git
   cd icp-evm-proxy
   ```

2. **Deploy Logger Canister**

   ```shell
   cd canisters/logger_canister
   dfx deploy --network=<chosen_icp_network>
   ```

   Note down the deployed canister ID.

3. **Configure Counter Canister**

   Edit `canisters/counter_canister/src/lib.rs` and update the `LOGGER_CANISTER_ID` constant with the ID obtained in the previous step.

4. **Deploy Counter Canister**

   ```shell
   cd ../counter_canister
   dfx deploy --network=<chosen_icp_network>
   ```

5. **Configure and Deploy EVM Adapter Proxy**

   ```shell
   cd ../../evm-adapter-proxy
   ```

   Edit `config.yaml` and update the ICP configuration with the deployed canister IDs.

   ```shell
   make run
   ```

6. **Configure and Run SubQuery Indexer**

   ```shell
   cd ../subq-indexer
   ```

   Edit `project.yaml` and update the URL of the deployed EVM Adapter Proxy.

   ```shell
   npm install
   npm run codegen
   npm run build
   docker-compose up
   ```

## Verification

To verify the deployment:

1. Interact with the Counter Canister to generate log entries.
2. Check the EVM Adapter Proxy's functionality:

   ```shell
   curl -X POST -H "Content-Type: application/json" \
   --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
   http://localhost:3030/rpc/v1
   ```

3. Query the translated data using one of the following methods:
   a. Use the SubQuery Indexer's GraphQL endpoint (typically at `http://localhost:3000`) to execute GraphQL queries.
   b. Directly query the PostgreSQL database.

## Deployed Canister URLs

You can interact with the deployed canisters using the following URLs:

- Counter Canister: <https://a4gq6-oaaaa-aaaab-qaa4q-cai.raw.ic0.app/?id=5kyqu-qyaaa-aaaak-qitna-cai>
- Logger Canister: <https://a4gq6-oaaaa-aaaab-qaa4q-cai.raw.ic0.app/?id=ydpfi-uiaaa-aaaal-qjupa-cai>

These URLs provide access to the canister interfaces, allowing you to interact with them directly through the Internet Computer's web interface.
