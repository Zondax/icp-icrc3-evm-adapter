---
title: SubQuery Indexer
sidebar_position: 3
---

# SubQuery Indexer

The SubQuery Indexer is a crucial component of the ICP-EVM Proxy project, designed to index and make queryable the EVM-compatible data translated by the EVM Adapter Proxy. This Node.js-based service provides a GraphQL API for efficient data retrieval and analysis.

## Key Features

1. **Data Indexing**: Processes and indexes EVM-compatible event logs translated from ICRC-3 logs.
2. **GraphQL API**: Exposes a flexible GraphQL API for querying indexed data.
3. **PostgreSQL Integration**: Utilizes PostgreSQL for persistent storage and efficient querying of indexed data.

## Implementation Overview

The SubQuery Indexer is built using the SubQuery framework, which provides a robust infrastructure for blockchain data indexing. The main components include:

1. **Schema Definition**: Defines the structure of the data to be indexed.
2. **Mapping Functions**: Processes incoming events and transforms them into the defined schema.
3. **Project Configuration**: Specifies the data source, network, and other indexing parameters.

## Data Model

The project defines a `Log` entity to represent the events from the EVM-compatible logs:
