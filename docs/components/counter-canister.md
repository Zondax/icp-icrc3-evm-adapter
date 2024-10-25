# Counter Canister

The Counter Canister is a Rust-based canister that demonstrates basic functionality on the Internet Computer Protocol (ICP). It now includes both counter operations and basic DEX-like functionality, serving as an example of how to generate events that can be logged and later translated into EVM-compatible formats.

## Functionality

The Counter Canister provides the following functions:

1. Increment the counter
2. Get the current counter value
3. Add currency pairs
4. Get currency pairs
5. Mint tokens
6. Burn tokens
7. Get token balance

All operations generate events that are sent to the Logger Canister.

## Counter Operations

### Increment

Increments the counter and logs the operation.

Example:

```shell
dfx canister call counter_canister increment
```

### Get

Retrieves the current counter value and logs the operation.

Example:

```shell
dfx canister call counter_canister get
```

## DEX-like Operations

### Add Currency Pair

Adds a new currency pair to the list of supported pairs.

Example:

```shell
dfx canister call counter_canister add_currency_pair '(record { base_currency = "ICP"; quote_currency = "USD" })'
```

### Get Currency Pairs

Retrieves the list of all supported currency pairs.

Example:

```shell
dfx canister call counter_canister get_currency_pairs
```

### Mint Tokens

Mints new tokens for a specified currency and recipient.

Example:

```shell
dfx canister call counter_canister mint_tokens '(record { currency = "ICP"; amount =    ; recipient = principal "2vxsx-fae" })'
```

### Burn Tokens

Burns existing tokens for a specified currency and owner.

Example:

```shell
dfx canister call counter_canister burn_tokens '(record { currency = "ICP"; amount = 50000000; owner = principal "2vxsx-fae" })'
```

### Get Token Balance

Retrieves the token balance for a specific user and currency.

Example:

```shell
dfx canister call counter_canister get_token_balance '(principal "2vxsx-fae", "ICP")'
```

## Event Generation

When any operation is performed, an event is generated and sent to the Logger Canister. The event includes:

- The action performed (e.g., increment, get, add_currency_pair, mint_tokens, burn_tokens)
- Relevant details of the operation
- A timestamp

These events are crucial for demonstrating the flow of data from ICP canisters to EVM-compatible systems through the EVM Adapter Proxy.

## Integration with Logger Canister

The Counter Canister interacts directly with the Logger Canister to store events. This interaction is typically done through inter-canister calls on the Internet Computer. When any operation occurs, the Counter Canister creates an event and sends it to the Logger Canister for storage.

## ICRC-3 Compatibility

The Counter Canister implements ICRC-3 compatible functions to support standardized logging and querying of events. These functions include:

- `icrc3_get_archives`
- `icrc3_get_tip_certificate`
- `icrc3_get_blocks`
- `icrc3_supported_block_types`

These functions allow for standardized interaction with the canister's event log, facilitating integration with other systems and services.

## Notes on Usage

- The amounts in mint and burn operations are represented in the smallest unit of the currency. For example, if ICP has 8 decimal places, 100000000 represents 1 ICP.
- The principal used in the examples ("2vxsx-fae") is a placeholder. Replace it with an actual principal ID when making calls.
- Make sure to replace "counter_canister" with the actual canister ID if it's different in your deployment.
