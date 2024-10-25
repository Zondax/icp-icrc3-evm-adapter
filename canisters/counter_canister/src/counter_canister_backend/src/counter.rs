use crate::state::STATE;
use crate::logging::log_operation;
use crate::constants::{INCREMENT_OPERATION, GET_OPERATION};
use ic_cdk_macros::*;

/// Increments the counter and logs the operation.
///
/// # Returns
///
/// * `u64` - The new value of the counter after incrementing.
#[update]
pub async fn increment() -> u64 {
    let new_value = STATE.with(|state| {
        let mut state = state.borrow_mut();
        state.counter += 1;
        state.counter
    });

    let _ = log_operation(INCREMENT_OPERATION, new_value).await;
    new_value
}

/// Gets the current value of the counter and logs the operation.
///
/// # Returns
///
/// * `u64` - The current value of the counter.
#[update]
pub async fn get() -> u64 {
    let value = STATE.with(|state| state.borrow().counter);
    let _ = log_operation(GET_OPERATION, value).await;
    value
}
