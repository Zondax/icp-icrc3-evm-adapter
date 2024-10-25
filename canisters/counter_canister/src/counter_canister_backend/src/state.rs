use candid::{CandidType, Deserialize, Principal, Nat};
use std::cell::RefCell;
use std::collections::{HashSet, HashMap};

#[derive(CandidType, Deserialize, Default)]
pub struct State {
    pub counter: u64,
    pub currency_pairs: HashSet<CurrencyPair>,
    pub token_balances: HashMap<(Principal, String), Nat>,
}

#[derive(CandidType, Deserialize, Clone, Debug, PartialEq, Eq, Hash)]
pub struct CurrencyPair {
    pub base_currency: String,
    pub quote_currency: String,
}

/// Thread-local storage for the canister state
thread_local! {
    pub static STATE: RefCell<State> = RefCell::default();
}
