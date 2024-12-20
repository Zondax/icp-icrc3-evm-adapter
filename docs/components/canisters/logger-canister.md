---
title: Logger Canister
sidebar_position: 3
---

# Logger Canister

The Logger Canister is a Rust-based canister that implements ICRC-3 compatible logging on the Internet Computer Protocol (ICP). It plays a crucial role in capturing and storing events generated by other canisters, such as the Counter Canister, in a standardized format.

## Functionality

The Logger Canister provides the following main functions:

1. Log events in ICRC-3 compatible format
2. Retrieve logged events based on various criteria
3. Manage log storage

## ICRC-3 Compatibility

The Logger Canister implements the ICRC-3 standard for event logging. This ensures that logged events are in a standardized format that can be easily translated to EVM-compatible logs by the EVM Adapter Proxy. The ICRC-3 format includes essential information such as timestamps, event types, and detailed event data.

## Event Logging Process

When an event occurs in another canister (e.g., the Counter Canister), it sends the event details to the Logger Canister. The Logger Canister then:

1. Validates the incoming event data
2. Formats the event according to the ICRC-3 standard
3. Assigns a unique identifier and timestamp to the event
4. Stores the event in its internal log storage

## Log Retrieval

The Logger Canister provides flexible querying capabilities to retrieve logged events using standard ICRC-3 methods.

## Log Storage and Indexing

In this proof of concept, logs are stored in memory using efficient data structures for quick insertion and retrieval. For production use, consider implementing:

- Persistent storage using stable memory
- Advanced indexing mechanisms for faster querying
- Log rotation or archiving strategies for managing large volumes of data

## Integration with Other Components

- **Counter Canister**: The Counter Canister (and potentially other canisters) call the Logger Canister to log events.
- **EVM Adapter Proxy**: The EVM Adapter Proxy retrieves logs from the Logger Canister for translation into EVM-compatible format.

## Limitations

As a proof of concept, this implementation has several limitations:

- Limited storage capacity due to in-memory storage
- Basic querying capabilities
- Lack of advanced security features

For production use, these limitations should be addressed to ensure scalability, security, and performance.

## Conclusion

The Logger Canister is a critical component in the ICP-EVM Proxy system, providing standardized event logging that bridges the gap between ICP canisters and EVM-compatible systems. Its implementation of the ICRC-3 standard ensures compatibility and ease of integration with other components in the ecosystem.
