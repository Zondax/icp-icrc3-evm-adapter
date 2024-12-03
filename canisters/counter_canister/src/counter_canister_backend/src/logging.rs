use candid::{Principal};
use crate::types::Value;

/// Logs an operation to the logger canister.
///
/// # Arguments
///
/// * `operation` - A string slice that holds the name of the operation.
/// * `details` - An object implementing the `std::fmt::Debug` trait, containing details of the operation.
///
/// # Returns
///
/// * `Result<(), String>` - Ok if logging was successful, Err with an error message otherwise.
pub async fn log_operation(operation: &str, details: impl std::fmt::Debug) -> Result<(), String> {
    let logger_id = "ydpfi-uiaaa-aaaal-qjupa-cai"; // only for PoC purposes
    let logger_id_principal = Principal::from_text(logger_id)
        .map_err(|e| format!("Failed to parse logger ID: {:?}", e))?;

    let details = Value::Map(vec![
        ("value".to_string(), Value::Text(format!("{:?}", details))),
    ]);

    ic_cdk::call(logger_id_principal, "add_entry", (operation.to_string(), details, ic_cdk::caller()))
        .await
        .map_err(|e| format!("Failed to log operation: {:?}", e))?;

    Ok(())
}
