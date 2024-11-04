mod state;
mod types;
mod icrc3;
mod logging;
mod utils;

use candid::Principal;

pub use types::{
    LogEntry, Block, Value,
    GetArchivesArgs, GetArchivesResult, ArchiveInfo,
    DataCertificate, GetBlocksArgs, GetBlocksResult,
    BlockInfo, ArchivedBlocksRange, BlockTypeInfo
};

pub use icrc3::{
    icrc3_get_archives,
    icrc3_get_tip_certificate,
    icrc3_get_blocks,
    icrc3_supported_block_types
};

pub use logging::*;

/// Gets the chain ID.
#[ic_cdk_macros::query]
fn chain_id() -> String {
    "314160".to_string()
}

/// Gets the network version.
#[ic_cdk_macros::query]
fn net_version() -> String {
    "314160".to_string()
}

// Candid export
ic_cdk::export_candid!();
