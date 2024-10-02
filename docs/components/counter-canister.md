# Counter Canister

The Counter Canister is a simple Rust-based canister that demonstrates basic functionality on the Internet Computer Protocol (ICP). It serves as an example of how to generate events that can be logged and later translated into EVM-compatible formats.

## Functionality

The Counter Canister provides the following functions:

1. Increment the counter
2. Get the current counter value

Both operations generate an event that is sent to the Logger Canister.

## Event Generation

When the counter is incremented or its value is retrieved, an event is generated and sent to the Logger Canister. The event includes:

- The action performed (increment or get)
- The current counter value
- A timestamp

These events are crucial for demonstrating the flow of data from ICP canisters to EVM-compatible systems through the EVM Adapter Proxy.

## Integration with Logger Canister

The Counter Canister interacts directly with the Logger Canister to store events. This interaction is typically done through inter-canister calls on the Internet Computer. When either an increment or get operation occurs, the Counter Canister creates an event and sends it to the Logger Canister for storage.
