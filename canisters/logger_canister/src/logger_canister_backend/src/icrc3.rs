use candid::{Principal, Nat};
use ic_cdk_macros::*;
use num_traits::ToPrimitive;
use sha2::{Sha256, Digest};
use crate::state::BLOCKS;
use crate::types::{
    GetArchivesArgs, GetArchivesResult, ArchiveInfo,
    DataCertificate, GetBlocksArgs, GetBlocksResult,
    BlockInfo, BlockTypeInfo
};

// Note: These are stub implementations for Proof of Concept (PoC) purposes only.

// ICRC-3 Interface Implementations

/// Retrieves information about archives.
#[query]
pub fn icrc3_get_archives(args: GetArchivesArgs) -> GetArchivesResult {
    // Mock implementation
    let archives = vec![
        ArchiveInfo {
            canister_id: Principal::from_text("aaaaa-aa").unwrap_or_else(|_| Principal::anonymous()),
            start: Nat::from(0u64),
            end: Nat::from(999u64),
        },
        ArchiveInfo {
            canister_id: Principal::from_text("bbbbb-bb").unwrap_or_else(|_| Principal::anonymous()),
            start: Nat::from(1000u64),
            end: Nat::from(1999u64),
        },
    ];

    // If 'from' is provided, return archives starting from that canister
    let filtered_archives = match args.from {
        Some(from) => archives.into_iter().skip_while(|info| info.canister_id != from).collect(),
        None => archives,
    };

    GetArchivesResult(filtered_archives)
}

/// Retrieves the tip certificate for data integrity verification.
///
/// # Returns
///
/// * `Option<DataCertificate>` - The tip certificate if available, or None.
#[query]
pub fn icrc3_get_tip_certificate() -> Option<DataCertificate> {
    BLOCKS.with(|blocks| {
        let blocks = blocks.borrow();
        if blocks.is_empty() {
            return None;
        }
        
        let last_block = blocks.last().unwrap();
        let last_block_index = Nat::from(blocks.len() - 1);
        
        let mut hasher = Sha256::new();
        let mut certificate_data = Vec::new();
        
        let mut leb_encoded = Vec::new();
        leb128::write::unsigned(&mut leb_encoded, last_block_index.0.to_u64().unwrap_or(0)).unwrap();
        certificate_data.extend_from_slice(&leb_encoded);
        
        certificate_data.extend_from_slice(&last_block.hash);
        
        hasher.update(&certificate_data);
        let hash_tree = hasher.finalize().to_vec();
        
        Some(DataCertificate {
            certificate: certificate_data,
            hash_tree,
        })
    })
}

/// Retrieves blocks within a specified range.
///
/// # Arguments
///
/// * `args` - A `GetBlocksArgs` struct specifying the start and length of the block range.
///
/// # Returns
///
/// * `GetBlocksResult` - The result containing log length, blocks, and archived block ranges.
#[query]
pub fn icrc3_get_blocks(args: GetBlocksArgs) -> GetBlocksResult {
    let start = args.start.0.to_u64().unwrap_or(0);
    let length = args.length.0.to_u64().unwrap_or(0);

    BLOCKS.with(|blocks| {
        let blocks = blocks.borrow();
        let log_length = blocks.len() as u64;
        let blocks = blocks.iter()
            .skip(start as usize)
            .take(length as usize)
            .enumerate()
            .map(|(index, block)| BlockInfo {
                id: Nat::from(start + index as u64),
                block: crate::types::Value::Map(vec![
                    ("id".to_string(), crate::types::Value::Nat(block.id.clone())),
                    ("hash".to_string(), crate::types::Value::Blob(block.hash.clone())),
                    ("phash".to_string(), crate::types::Value::Blob(block.phash.clone())),
                    ("btype".to_string(), crate::types::Value::Text(block.btype.clone())),
                    ("ts".to_string(), crate::types::Value::Nat(Nat::from(block.ts))),
                    ("finalized".to_string(), crate::types::Value::Text(block.finalized.to_string())),
                    ("entries".to_string(), crate::types::Value::Array(block.entries.iter().map(|entry| {
                        crate::types::Value::Map(vec![
                            ("timestamp".to_string(), crate::types::Value::Nat(Nat::from(entry.timestamp))),
                            ("operation".to_string(), crate::types::Value::Text(entry.operation.clone())),
                            ("details".to_string(), entry.details.clone()),
                            ("caller".to_string(), crate::types::Value::Text(entry.caller.clone())),
                        ])
                    }).collect())),
                ])
            })
            .collect();

        // NOTE: For PoC purposes, we're not actually implementing archived blocks.
        // In a real implementation, this would check for and return actual archived blocks.
        let archived_blocks = vec![];

        GetBlocksResult {
            log_length: Nat::from(log_length),
            blocks,
            archived_blocks,
        }
    })
}

/// Returns supported block types.
///
/// # Returns
///
/// * `Vec<BlockTypeInfo>` - A vector of supported block types and their corresponding URLs.
#[query]
pub fn icrc3_supported_block_types() -> Vec<BlockTypeInfo> {
    vec![
        BlockTypeInfo {
            block_type: "log_entry".to_string(),
            url: "https://github.com/dfinity/ICRC-1/blob/main/standards/ICRC-3/README.md".to_string(),
        }
    ]
}