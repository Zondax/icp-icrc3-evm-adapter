use crate::state::{STATE, CurrencyPair};
use crate::logging::log_operation;
use candid::{CandidType, Deserialize, Principal, Nat};
use ic_cdk_macros::*;
use crate::constants::{ADD_CURRENCY_PAIR_OPERATION, MINT_TOKENS_OPERATION, BURN_TOKENS_OPERATION};

#[derive(CandidType, Deserialize, Debug)]
pub struct MintOperation {
    pub currency: String,
    pub amount: Nat,
    pub recipient: Principal,
}

#[derive(CandidType, Deserialize, Debug)]
pub struct BurnOperation {
    pub currency: String,
    pub amount: Nat,
    pub owner: Principal,
}

/// Adds a new currency pair to the DEX.
///
/// # Arguments
///
/// * `pair` - A `CurrencyPair` struct representing the currency pair to be added.
#[update]
pub async fn add_currency_pair(pair: CurrencyPair) {
    STATE.with(|state| {
        state.borrow_mut().currency_pairs.insert(pair.clone());
    });
    let _ = log_operation(ADD_CURRENCY_PAIR_OPERATION, format!("{:?}", pair)).await;
}

/// Retrieves all currency pairs currently available in the DEX.
///
/// # Returns
///
/// * `Vec<CurrencyPair>` - A vector containing all currency pairs in the DEX.
#[query]
pub fn get_currency_pairs() -> Vec<CurrencyPair> {
    STATE.with(|state| state.borrow().currency_pairs.iter().cloned().collect())
}

/// Mints new tokens for a specified recipient.
///
/// # Arguments
///
/// * `operation` - A `MintOperation` struct containing the details of the minting operation.
///
/// # Returns
///
/// * `Result<(), String>` - Ok if the minting was successful, Err with an error message otherwise.
#[update]
pub async fn mint_tokens(operation: MintOperation) -> Result<(), String> {
    STATE.with(|state| {
        let mut state = state.borrow_mut();
        
        if !state.currency_pairs.iter().any(|pair| pair.base == operation.currency || pair.quote == operation.currency) {
            return Err(format!("Currency {} is not listed", operation.currency));
        }
        
        let balance = state.token_balances
            .entry((operation.recipient, operation.currency.clone()))
            .or_insert(Nat::from(0u64));
        *balance += operation.amount.clone();
        Ok(())
    })?;

    let _ = log_operation(MINT_TOKENS_OPERATION, format!("{:?}", operation)).await;
    Ok(())
}

/// Burns tokens from a specified owner's balance.
///
/// # Arguments
///
/// * `operation` - A `BurnOperation` struct containing the details of the burning operation.
///
/// # Returns
///
/// * `Result<(), String>` - Ok if the burning was successful, Err with an error message otherwise.
///
/// # Errors
///
/// Returns an error if:
/// * The owner has no balance for the specified currency.
/// * The owner has insufficient balance to burn the requested amount.
#[update]
pub async fn burn_tokens(operation: BurnOperation) -> Result<(), String> {
    STATE.with(|state| {
        let mut state = state.borrow_mut();
        let balance = state.token_balances
            .get_mut(&(operation.owner, operation.currency.clone()))
            .ok_or("No balance for this user and currency")?;
        
        if *balance < operation.amount {
            return Err("Insufficient balance".to_string());
        }
        
        *balance -= operation.amount.clone();
        Ok(())
    })?;

    let _ = log_operation(BURN_TOKENS_OPERATION, format!("{:?}", operation)).await;
    Ok(())
}

/// Retrieves the token balance for a specific user and currency.
///
/// # Arguments
///
/// * `user` - The `Principal` of the user whose balance is being queried.
/// * `currency` - A `String` representing the currency for which to check the balance.
///
/// # Returns
///
/// * `Nat` - The balance of the specified currency for the given user. Returns 0 if no balance is found.
#[query]
pub fn get_token_balance(user: Principal, currency: String) -> Nat {
    STATE.with(|state| {
        state.borrow().token_balances
            .get(&(user, currency))
            .cloned()
            .unwrap_or_else(|| Nat::from(0u64))
    })
}
