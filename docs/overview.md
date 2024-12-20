---
title: Overview
sidebar_position: 1
---

# ICP-EVM Proxy Overview

The ICP-EVM Proxy is a proof-of-concept (PoC) project designed to enhance interoperability between the Internet Computer Protocol (ICP) and EVM-compatible platforms. This document provides a high-level overview of the project's components and functionality.

## Project Objective

The main goal of this project is to demonstrate a bridge between ICP's native capabilities and the widely adopted Ethereum Virtual Machine (EVM) ecosystem. By translating ICP events into EVM-compatible formats, we enable easier integration with existing Ethereum tooling and infrastructure.

## Key Components

1. **Counter Canister**: A simple Rust-based canister that serves as an example of basic ICP functionality. It generates events when incrementing or retrieving the counter value.

2. **DEX Canister**: A Rust-based canister that implements basic decentralized exchange functionality. It manages currency pairs, token operations (mint/burn), and generates ICRC-3 compatible events for all operations.

3. **Logger Canister**: A Rust-based canister that implements ICRC-3 compatible logging. It captures and stores events generated by other canisters (like the Counter and DEX Canisters) in a standardized format.

4. **EVM Adapter Proxy**: A Go-based service that acts as a bridge between ICP and EVM-compatible systems. It performs two main functions:
   - Retrieves ICRC-3 log data from the Logger Canister and translates it into EVM-compatible event logs.
   - Exposes EVM RPC compatible methods, allowing interaction with ICP canisters using familiar Ethereum tooling.

5. **SubQuery Indexer**: A Node.js-based service that indexes the translated EVM-compatible data. It provides a GraphQL API for efficient querying of the indexed data.

## Workflow

1. Users can interact with both Counter and DEX Canisters:
   - Counter operations: increment and get value
   - DEX operations: manage currency pairs, mint/burn tokens, check balances
2. All operations generate events that are captured and stored by the Logger Canister in ICRC-3 format
3. The EVM Adapter Proxy retrieves logs from the Logger Canister and translates them into EVM-compatible formats
4. The SubQuery Indexer processes the translated data, making it queryable through a GraphQL API

## Use Cases

This PoC demonstrates potential use cases such as:

- Enabling Ethereum developers to interact with ICP canisters using familiar tools and methods
- Allowing ICP dapps to integrate with existing Ethereum-based analytics or monitoring tools
- Facilitating cross-chain applications that span both ICP and EVM-compatible blockchains
- Demonstrating DEX functionality with standardized event logging and querying

## Limitations

As a proof-of-concept, this project has several limitations and is not intended for production use. Some key points to note:

- Simplified implementation of ICRC-3 logging and EVM event translation
- Limited error handling and edge case coverage
- Not optimized for high-throughput or large-scale deployments
- In-memory storage for logs in the Logger Canister, which is not suitable for long-term or high-volume data storage
- Basic DEX functionality without order matching or complex trading features

Component-specific documentation:

- [Counter Canister](./components/canisters/counter-canister.md)
- [DEX Canister](./components/canisters/dex-canister.md)
- [Logger Canister](./components/canisters/logger-canister.md)
- [EVM Adapter Proxy](./components/evm-adapter-proxy.md)
- [SubQuery Indexer](./components/subquery-indexer.md)
