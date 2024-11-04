# DEX Canister

The DEX (Decentralized Exchange) Canister is a Rust-based canister that implements basic decentralized exchange functionality on the Internet Computer Protocol (ICP). It provides trading pair management, token operations, and generates events that are logged using the ICRC-3 standard.

## Functionality

The DEX Canister provides the following functions:

1. Add and manage currency pairs with exchange rates
2. Query available currency pairs
3. Mint tokens for recipients
4. Burn tokens from owners
5. Query token balances
6. ICRC-3 compatible logging and querying

All operations generate events that are sent to the Logger Canister.

## Currency Pair Operations

### Add Currency Pair

Adds a new currency pair with its exchange rate to the DEX. The operation validates that:

- The rate is not zero
- The currencies are valid
- The pair doesn't already exist

Example:

```shell
dfx canister call dex_canister add_currency_pair '(record { 
    base_currency = "ICP"; 
    quote_currency = "BTC";
    rate = 31337;
})'
```

### Get Currency Pairs

Retrieves all currently supported currency pairs in the DEX, including their current exchange rates.

Example:

```shell
dfx canister call dex_canister get_currency_pairs
```

## Token Operations

### Mint Tokens

Creates new tokens for a specified recipient. The operation:

- Validates the currency exists in a trading pair
- Updates the recipient's balance
- Logs the minting event

Example:

```shell
dfx canister call dex_canister mint_tokens '(record { 
    currency = "ICP"; 
    amount = 100000000; 
    recipient = principal "2vxsx-fae" 
})'
```

### Burn Tokens

Removes tokens from an owner's balance. The operation:

- Verifies sufficient balance exists
- Validates the currency
- Updates the owner's balance
- Logs the burning event

Example:

```shell
dfx canister call dex_canister burn_tokens '(record { 
    currency = "ICP"; 
    amount = 50000000; 
    owner = principal "2vxsx-fae" 
})'
```

### Get Token Balance

Retrieves the current token balance for a specific user and currency combination.

Example:

```shell
dfx canister call dex_canister get_token_balance '(principal "2vxsx-fae", "ICP")'
```

## Integration with Logger Canister

The DEX Canister integrates with the Logger Canister through inter-canister calls to maintain a complete audit trail of all operations. Each operation is logged with detailed information to ensure traceability and compliance.

## ICRC-3 Compatibility

The Counter Canister implements ICRC-3 compatible functions to support standardized logging and querying of events. These functions include:

- `icrc3_get_archives`
- `icrc3_get_tip_certificate`
- `icrc3_get_blocks`
- `icrc3_supported_block_types`

These functions allow for standardized interaction with the canister's event log, facilitating integration with other systems and services.
