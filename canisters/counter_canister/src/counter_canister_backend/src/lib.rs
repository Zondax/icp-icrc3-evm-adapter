mod state;
mod counter;
mod types;
mod icrc3;
mod logging;
mod constants;

pub use types::{
    LogEntry, Block, Value,
    GetArchivesArgs, GetArchivesResult, ArchiveInfo,
    DataCertificate, GetBlocksArgs, GetBlocksResult,
    BlockInfo, ArchivedBlocksRange, BlockTypeInfo
};

pub use counter::{increment, get};
pub use icrc3::{icrc3_get_archives, icrc3_get_tip_certificate, icrc3_get_blocks, icrc3_supported_block_types};

// Candid export
ic_cdk::export_candid!();

#[cfg(test)]
mod tests;
