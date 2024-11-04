use std::cell::RefCell;
use crate::types::{Block, LogEntry};

thread_local! {
    pub static BLOCKS: RefCell<Vec<Block>> = RefCell::new(Vec::new());
    pub static LOG: RefCell<Vec<LogEntry>> = RefCell::new(Vec::new());
}
