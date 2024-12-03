use candid::{CandidType, Deserialize, Nat, Principal};
use serde::Serialize;

#[derive(CandidType, Deserialize, Serialize, Clone, Debug)]
pub struct GetArchivesArgs {
    pub from: Option<Principal>,
}

#[derive(CandidType, Serialize, Clone, Debug)]
pub struct GetArchivesResult(pub Vec<ArchiveInfo>);

#[derive(CandidType, Serialize, Clone, Debug)]
pub struct ArchiveInfo {
    pub canister_id: Principal,
    pub start: Nat,
    pub end: Nat,
}

#[derive(CandidType, Serialize, Clone, Debug)]
pub struct DataCertificate {
    pub certificate: Vec<u8>,
    pub hash_tree: Vec<u8>,
}

#[derive(CandidType, Deserialize, Serialize, Clone, Debug)]
pub struct GetBlocksArgs {
    pub start: Nat,
    pub length: Nat,
}

#[derive(CandidType, Serialize, Clone, Debug)]
pub struct GetBlocksResult {
    pub log_length: Nat,
    pub blocks: Vec<BlockInfo>,
    pub archived_blocks: Vec<ArchivedBlocksRange>,
}

#[derive(CandidType, Serialize, Clone, Debug)]
pub struct BlockInfo {
    pub id: Nat,
    pub block: crate::types::Value,
}

#[derive(CandidType, Serialize, Clone, Debug)]
pub struct ArchivedBlocksRange {
    pub args: GetBlocksArgs,
    pub callback: String,
}

#[derive(CandidType, Serialize, Clone, Debug)]
pub struct BlockTypeInfo {
    pub block_type: String,
    pub url: String,
}

#[derive(CandidType, Deserialize, Serialize, Clone, Debug)]
pub struct LogEntry {
    pub timestamp: u64,
    pub operation: String,
    pub details: Value,
    pub caller: String,
}

#[derive(CandidType, Deserialize, Serialize, Clone, Debug)]
pub struct Block {
    pub id: Nat,
    pub hash: Vec<u8>,
    pub phash: Vec<u8>,
    pub btype: String,
    pub ts: u64,
    pub entries: Vec<LogEntry>,
    pub finalized: bool,
}

#[derive(CandidType, Serialize, Deserialize, Clone, Debug)]
pub enum Value {
    Blob(Vec<u8>),
    Text(String),
    Nat(Nat),
    Int(i128),
    Array(Vec<Value>),
    Map(Vec<(String, Value)>),
} 