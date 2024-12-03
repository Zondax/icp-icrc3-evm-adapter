use crate::state::{BLOCKS, LOG};
use crate::types::{LogEntry, Block, Value};
use crate::utils::calculate_hash;
use candid::{Principal, Nat};
use ic_cdk_macros::{update, query};

/// Adds a new log entry and creates a new block if needed.
/// 
/// This function creates a new log entry with the provided details and adds it to the log.
/// It also triggers the creation of a new block if necessary.
///
/// # Arguments
///
/// * `operation` - A string describing the operation being logged.
/// * `details` - A `Value` containing additional details about the operation.
/// * `caller` - The `Principal` of the entity calling this function.
///
/// # Panics
///
/// Panics if the `operation` string is empty.
///
/// # Note
///
/// For proof-of-concept purposes, this implementation creates a new block every 2 entries.
#[update]
pub fn add_entry(operation: String, details: Value, caller: Principal) {
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

/// Creates a new block if needed based on the current state of the blockchain.
/// 
/// This function is responsible for managing the creation and finalization of blocks.
/// It creates a new block when:
/// 1. The current block has reached the maximum number of entries (2 in this case).
/// 2. There are no existing blocks.
///
/// A block is finalized when:
/// 1. It reaches the maximum number of entries (2 in this case).
/// 2. A new block is created, which automatically finalizes the previous block.
///
/// # Arguments
///
/// * `current_time` - The current timestamp.
/// * `new_entry` - The new `LogEntry` to be added to the block.
///
/// # Note
///
/// This is a simplified implementation for proof-of-concept purposes.
pub fn create_block_if_needed(current_time: u64, new_entry: LogEntry) {
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

/// Retrieves all log entries.
///
/// This function is intended for testing purposes only. It returns a vector
/// containing all the log entries currently stored in the system.
///
/// # Returns
///
/// A `Vec<LogEntry>` containing all stored log entries.
#[query]
pub fn get_logs() -> Vec<LogEntry> {
    LOG.with(|log| log.borrow().clone())
}
