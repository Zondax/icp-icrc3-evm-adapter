use crate::types::{Block, Value};
use sha2::{Sha256, Digest};
use num_traits::ToBytes;

/// Calculates the hash of a Block.
///
/// This function computes a SHA-256 hash of the block's contents, including its ID, previous hash,
/// type, timestamp, and all entries within the block.
///
/// # Arguments
///
/// * `block` - A reference to the Block to be hashed.
///
/// # Returns
///
/// * `Vec<u8>` - The resulting hash as a vector of bytes.
pub fn calculate_hash(block: &Block) -> Vec<u8> {
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

/// Serializes a Value into a vector of bytes.
///
/// This function takes a Value enum and converts it into a byte representation,
/// which is then hashed using SHA-256. The method of serialization depends on the
/// variant of the Value enum.
///
/// # Arguments
///
/// * `value` - A reference to the Value to be serialized.
///
/// # Returns
///
/// * `Vec<u8>` - The resulting serialized and hashed value as a vector of bytes.
pub fn serialize_value(value: &Value) -> Vec<u8> {
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
