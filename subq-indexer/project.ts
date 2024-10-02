import {
  EthereumProject,
  EthereumDatasourceKind,
  EthereumHandlerKind,
} from "@subql/types-ethereum";

// Can expand the Datasource processor types via the generic param
const project: EthereumProject = {
  specVersion: "1.0.0",
  version: "0.0.1",
  name: "poc-subquery-icp-icrc3",
  description:
    "This project can be use as a starting point for developing your new Altlayer OP Demo Testnet SubQuery project",
  runner: {
    node: {
      name: "@subql/node-ethereum",
      version: ">=3.0.0",
    },
    query: {
      name: "@subql/query",
      version: "*",
    },
  },
  schema: {
    file: "./schema.graphql",
  },
  network: {
    /**
     * chainId is the EVM Chain ID, for Altlayer OP Demo Testnet this is 20240219
     * https://chainlist.org/chain/20240219
     */
    chainId: "CHAIN_ID_ENV_VAR_REF",
    /**
     * These endpoint(s) should be public non-pruned archive node
     * We recommend providing more than one endpoint for improved reliability, performance, and uptime
     * Public nodes may be rate limited, which can affect indexing speed
     * When developing your project we suggest getting a private API key
     * If you use a rate limited endpoint, adjust the --batch-size and --workers parameters
     * These settings can be found in your docker-compose.yaml, they will slow indexing but prevent your project being rate limited
     */
    endpoint: ["ENDPOINT_URL_ENV_VAR_REF"],
  },
  dataSources: [
    {
      kind: EthereumDatasourceKind.Runtime,
      // Block height at which the smart contract was deployed
      startBlock: -1,
      options: {
        abi: "erc20",
      },
      assets: new Map([["erc20", { file: "./abis/erc20.abi.json" }]]),
      mapping: {
        file: "./dist/index.js",
        handlers: [
          {
            handler: "handleLog",
            kind: EthereumHandlerKind.Event,
          },
          {
            handler: "handleLog",
            kind: EthereumHandlerKind.Call,
          },
          {
            handler: "handleLog",
            kind: EthereumHandlerKind.Block,
          },
        ],
      },
    },
  ],
  repository: "https://github.com/zondax/icp-icrc3-evm-adapter/subq-indexer",
};

// Must set default to the project instance
export default project;
