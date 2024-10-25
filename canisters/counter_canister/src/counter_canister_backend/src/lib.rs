mod state;
mod counter;
mod dex;
mod icrc3;
mod logging;
mod constants;

use candid::{Principal, Nat};
use state::CurrencyPair;
use dex::{MintOperation, BurnOperation};
use icrc3::{GetArchivesArgs, ArchiveInfo, GetBlocksArgs, GetBlocksResult, DataCertificate, BlockTypeInfo};

pub use counter::{increment, get};
pub use dex::{add_currency_pair, get_currency_pairs, mint_tokens, burn_tokens, get_token_balance};
pub use icrc3::{icrc3_get_archives, icrc3_get_tip_certificate, icrc3_get_blocks, icrc3_supported_block_types};

// Candid export
ic_cdk::export_candid!();
