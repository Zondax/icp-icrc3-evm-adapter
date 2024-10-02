use ic_cdk_macros::{query, update};
use candid::{CandidType, Principal, Nat};
use sha2::{Sha256, Digest};
use std::cell::RefCell;
use num_traits::{ToBytes, ToPrimitive};
use serde::{Serialize, Deserialize};
use leb128;

/// Represents a value that can be stored in a log entry or block.
#[derive(CandidType, Deserialize, Serialize, Clone, Debug)]
pub enum Value {
    Blob(Vec<u8>),
    Text(String),
    Nat(Nat),
    Int(i128),
    Array(Vec<Value>),
    Map(Vec<(String, Value)>),
}

/// Represents a block in the ICRC-3 log.
#[derive(CandidType, Deserialize, Serialize, Clone, Debug)]
struct Block {
    id: Nat,
    hash: Vec<u8>,
    phash: Vec<u8>,
    btype: String,
    ts: u64,
    entries: Vec<LogEntry>,
    finalized: bool,
}

#[derive(CandidType, Deserialize, Serialize, Clone, Debug)]
struct LogEntry {
    timestamp: u64,
    operation: String,
    details: Value,
    caller: String,
}

thread_local! {
    static BLOCKS: RefCell<Vec<Block>> = RefCell::new(Vec::new());
    static LOG: RefCell<Vec<LogEntry>> = RefCell::new(Vec::new());
}

/// ICRC-3 implementation: Get archives information.
/// 
/// NOTE: mock implementation for PoC purposes.
#[query]
fn icrc3_get_archives(args: GetArchivesArgs) -> GetArchivesResult {
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

/// ICRC-3 implementation: Get blocks from the log.
/// 
/// NOTE: This is a simplified implementation for PoC purposes.
/// In a production environment, this would handle archived blocks and callbacks properly.
#[query]
fn icrc3_get_blocks(args: GetBlocksArgs) -> GetBlocksResult {
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
                block: Value::Map(vec![
                    ("id".to_string(), Value::Nat(block.id.clone())),
                    ("hash".to_string(), Value::Blob(block.hash.clone())),
                    ("phash".to_string(), Value::Blob(block.phash.clone())),
                    ("btype".to_string(), Value::Text(block.btype.clone())),
                    ("ts".to_string(), Value::Nat(Nat::from(block.ts))),
                    ("finalized".to_string(), Value::Text(block.finalized.to_string())),
                    ("entries".to_string(), Value::Array(block.entries.iter().map(|entry| {
                        Value::Map(vec![
                            ("timestamp".to_string(), Value::Nat(Nat::from(entry.timestamp))),
                            ("operation".to_string(), Value::Text(entry.operation.clone())),
                            ("details".to_string(), entry.details.clone()),
                            ("caller".to_string(), Value::Text(entry.caller.clone())),
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

/// ICRC-3 implementation: Get tip certificate.
#[query]
fn icrc3_get_tip_certificate() -> Option<DataCertificate> {
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

/// Adds a new log entry and creates a new block if needed.
/// 
/// NOTE: For PoC purposes, this implementation creates a new block every 2 entries.
#[update]
fn add_entry(operation: String, details: Value, caller: Principal) {
    if operation.is_empty() {
        ic_cdk::trap("Operation cannot be empty");
    }

    let timestamp = ic_cdk::api::time();
    let entry = LogEntry {
        timestamp,
        operation,
        details,
        caller: caller.to_string(),
    };
    LOG.with(|log| log.borrow_mut().push(entry.clone()));
    
    create_block_if_needed(timestamp, entry);
}

/// Creates a new block if needed.
/// 
/// NOTE: simplified implementation for PoC purposes.
/// A block is finalized when:
/// 1. It reaches the maximum number of entries (2 in this case).
/// 2. A new block is created, which automatically finalizes the previous block.
fn create_block_if_needed(current_time: u64, new_entry: LogEntry) {
    BLOCKS.with(|blocks| {
        let mut blocks = blocks.borrow_mut();
        let entries_count = blocks.last().map_or(0, |b| b.entries.len());
        
        if entries_count >= 1 || blocks.is_empty() {
            // Finalize the previous block if it exists
            if let Some(last_block) = blocks.last_mut() {
                last_block.finalized = true;
            }

            let phash = if let Some(last_block) = blocks.last() {
                last_block.hash.clone()
            } else {
                // Genesis block pHash can be a fixed or empty hash
                vec![0u8; 32] // Represents 0x000...000 (32 bytes)
            };

            let mut new_block = Block {
                id: Nat::from(blocks.len()),
                hash: vec![],
                phash,
                btype: "log_entry".to_string(),
                ts: current_time,
                entries: vec![new_entry],
                finalized: false,  // New block is not finalized initially
            };

            new_block.hash = calculate_hash(&new_block);
            blocks.push(new_block);
        } else {
            if let Some(current_block) = blocks.last_mut() {
                current_block.entries.push(new_entry);
                // If this entry makes the block full, finalize it
                if current_block.entries.len() >= 2 {
                    current_block.finalized = true;
                }
            }
        }
    });
}

/// Gets the chain ID
#[query]
fn chain_id() -> String {
    "314160".to_string()
}

/// Gets the network version.
/// 
/// NOTE: mock implementation for PoC purposes.
#[query]
fn net_version() -> String {
    "314160".to_string()
}

fn calculate_hash(block: &Block) -> Vec<u8> {
    let mut hasher = Sha256::new();

    hasher.update(&block.id.0.to_bytes_le());
    hasher.update(&block.phash);
    hasher.update(block.btype.as_bytes());
    hasher.update(block.ts.to_be_bytes());

    for entry in &block.entries {
        hasher.update(entry.timestamp.to_be_bytes());
        hasher.update(entry.operation.as_bytes());
        hasher.update(&serialize_value(&entry.details));
        hasher.update(entry.caller.as_bytes());
    }

    hasher.finalize().to_vec()
}

fn serialize_value(value: &Value) -> Vec<u8> {
    let mut hasher = Sha256::new();
    match value {
        Value::Nat(n) => {
            hasher.update(&n.0.to_bytes_le());
        },
        Value::Int(i) => {
            hasher.update(&i.to_le_bytes());
        },
        Value::Text(s) => {
            hasher.update(s.as_bytes());
        },
        Value::Blob(b) => {
            hasher.update(b);
        },
        Value::Array(arr) => {
            for elem in arr {
                hasher.update(&serialize_value(elem));
            }
        },
        Value::Map(map) => {
            let mut hashes: Vec<(Vec<u8>, Vec<u8>)> = map.iter()
                .map(|(k, v)| (Sha256::digest(k.as_bytes()).to_vec(), serialize_value(v)))
                .collect();
            hashes.sort();
            for (key_hash, val_hash) in hashes {
                hasher.update(&key_hash);
                hasher.update(&val_hash);
            }
        },
    }
    hasher.finalize().to_vec()
}

// Just for testing purposes
#[query]
fn get_logs() -> Vec<LogEntry> {
    LOG.with(|log| log.borrow().clone())
}

// Type definitions for ICRC-3 interface
#[derive(CandidType, Deserialize, Serialize, Clone, Debug)]
struct GetArchivesArgs {
    from: Option<Principal>,
}

#[derive(CandidType, Serialize, Clone, Debug)]
struct GetArchivesResult(Vec<ArchiveInfo>);

#[derive(CandidType, Serialize, Clone, Debug)]
struct ArchiveInfo {
    canister_id: Principal,
    start: Nat,
    end: Nat,
}

#[derive(CandidType, Deserialize, Serialize, Clone, Debug)]
struct GetBlocksArgs {
    start: Nat,
    length: Nat,
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

#[derive(CandidType, Serialize, Clone, Debug)]
struct DataCertificate {
    certificate: Vec<u8>,
    hash_tree: Vec<u8>,
}

#[derive(CandidType, Serialize, Clone, Debug)]
struct BlockTypeInfo {
    block_type: String,
    url: String,
}

// Candid export
ic_cdk::export_candid!();