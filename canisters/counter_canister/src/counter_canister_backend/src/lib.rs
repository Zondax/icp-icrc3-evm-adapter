use candid::{CandidType, Deserialize, Principal, Nat};
use ic_cdk_macros::*;
use std::cell::RefCell;
use serde::{Serialize};

/// Represents the state of the counter canister.
#[derive(CandidType, Deserialize, Default)]
struct State {
    /// The current value of the counter.
    counter: u64,
}

thread_local! {
    static STATE: RefCell<State> = RefCell::default();
}

/// Increments the counter and logs the operation.
///
/// # Returns
///
/// * `u64` - The new value of the counter after incrementing.
#[update]
async fn increment() -> u64 {
    let new_value = STATE.with(|state| {
        let mut state = state.borrow_mut();
        state.counter += 1;
        state.counter
    });

    // Log the increment operation
    let _ = log_operation("increment", new_value).await;

    new_value
}

/// Gets the current value of the counter and logs the operation.
///
/// # Returns
///
/// * `u64` - The current value of the counter.
#[update]
async fn get() -> u64 {
    let value = STATE.with(|state| state.borrow().counter);

    // Log the get operation
    let _ = log_operation("get", value).await;

    value
}

/// Logs an operation to the logger canister.
///
/// # Arguments
///
/// * `operation` - A string slice that holds the name of the operation.
/// * `value` - The value of the counter after the operation.
///
/// # Returns
///
/// * `Result<(), String>` - Ok if logging was successful, Err with an error message otherwise.
async fn log_operation(operation: &str, value: u64) -> Result<(), String> {
    let logger_id = "ydpfi-uiaaa-aaaal-qjupa-cai"; // only for PoC purposes
    let logger_id_principal = Principal::from_text(logger_id)
        .map_err(|e| format!("Failed to parse logger ID: {:?}", e))?;

    let details = Value::Map(vec![
        ("value".to_string(), Value::Nat(Nat::from(value))),
    ]);

    // Call the logger canister to log the message
    ic_cdk::call(logger_id_principal, "add_entry", (operation.to_string(), details, ic_cdk::caller()))
        .await
        .map_err(|e| format!("Failed to log operation: {:?}", e))?;

    Ok(())
}

// ICRC-3 interface implementations (stubs)

#[derive(CandidType, Deserialize)]
struct GetArchivesArgs {
    from: Option<Principal>,
}

#[derive(CandidType, Serialize, Clone, Debug)]
struct ArchiveInfo {
    canister_id: Principal,
    start: Nat,
    end: Nat,
}

/// Stub implementation of the ICRC-3 `get_archives` function.
#[query]
fn icrc3_get_archives(args: GetArchivesArgs) -> ArchiveInfo {
    // Mock implementation
    let archive = ArchiveInfo {
        canister_id: Principal::from_text("aaaaa-aa").unwrap_or_else(|_| Principal::anonymous()),
        start: Nat::from(0u64),
        end: Nat::from(999u64),
    };

    archive
}

#[derive(CandidType, Deserialize)]
struct DataCertificate {
    certificate: Vec<u8>,
    hash_tree: Vec<u8>,
}

/// Stub implementation of the ICRC-3 `get_tip_certificate` function.
#[query]
fn icrc3_get_tip_certificate() -> Option<DataCertificate> {
    None
}

#[derive(CandidType, Deserialize, Serialize, Clone, Debug)]
struct GetBlocksArgs {
    start: Nat,
    length: Nat,
}

#[derive(CandidType, Serialize, Deserialize, Clone, Debug)]
enum Value {
    Blob(Vec<u8>),
    Text(String),
    Nat(Nat),
    Int(i128),
    Array(Vec<Value>),
    Map(Vec<(String, Value)>),
}

#[derive(CandidType, Serialize, Clone, Debug)]
struct GetBlocksResult {
    log_length: Nat,
    blocks: Vec<BlockInfo>,
    archived_blocks: Vec<ArchivedBlocksRange>,
}

#[derive(CandidType, Serialize, Clone, Debug)]
struct BlockInfo {
    id: Nat,
    block: Value,
}

#[derive(CandidType, Clone, Debug)]
struct ArchivedBlocksRange {
    args: GetBlocksArgs,
    // NOTE: For PoC purposes, we're using a String to represent the callback function.
    callback: String,
}

impl Serialize for ArchivedBlocksRange {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut state = serializer.serialize_struct("ArchivedBlocksRange", 2)?;
        state.serialize_field("args", &self.args)?;
        state.serialize_field("callback", &self.callback)?;
        state.end()
    }
}

/// Stub implementation of the ICRC-3 `get_blocks` function.
#[query]
fn icrc3_get_blocks(_args: GetBlocksArgs) -> GetBlocksResult {
    GetBlocksResult {
        log_length: Nat::from(0u64),
        blocks: vec![],
        archived_blocks: vec![],
    }
}

/// ICRC-3 implementation: Get supported block types.
/// 
/// NOTE: simplified implementation for PoC purposes.
#[query]
fn icrc3_supported_block_types() -> Vec<BlockTypeInfo> {
    vec![
        BlockTypeInfo {
            block_type: "log_entry".to_string(),
            url: "https://github.com/dfinity/ICRC-1/blob/main/standards/ICRC-3/README.md".to_string(),
        }
    ]
}

#[derive(CandidType, Serialize, Clone, Debug)]
struct BlockTypeInfo {
    block_type: String,
    url: String,
}

// Candid export
ic_cdk::export_candid!();
