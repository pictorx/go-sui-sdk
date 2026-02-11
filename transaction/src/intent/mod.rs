use crate::Argument;
use crate::TransactionBuilder;

mod coin_with_balance;
pub use coin_with_balance::CoinWithBalance;

const MAX_GAS_OBJECTS: usize = 250; // 256
#[allow(unused)]
const MAX_COMMANDS: usize = 1000; // 1024
#[allow(unused)]
const MAX_INPUT_OBJECTS: usize = 2000; // 2048
const MAX_ARGUMENTS: usize = 500; // 512

pub(crate) type BoxError = Box<dyn std::error::Error + Send + Sync + 'static>;

// BUG FIX: Was `#[cfg(feature = "rpc")]` but the feature is named "intents" in
// Cargo.toml.  Using the wrong name silently compiled out the async `resolve`
// method, making `IntentResolver` a no-op trait and breaking all intent
// resolution at runtime (the intents map would never be drained, causing
// `try_build` to return `Err(Input("unable to resolve intents offline"))`).
#[cfg(feature = "intents")]
#[async_trait::async_trait]
pub(crate) trait IntentResolver: std::any::Any + std::fmt::Debug + Send + Sync {
    async fn resolve(
        &self,
        builder: &mut TransactionBuilder,
        client: &mut sui_rpc::Client,
    ) -> Result<(), BoxError>;
}

// Stub for builds where the "intents" feature is disabled.
#[cfg(not(feature = "intents"))]
pub(crate) trait IntentResolver: std::any::Any + std::fmt::Debug + Send + Sync {}