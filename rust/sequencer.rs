#![no_std]
#![no_main]

use sp_std::collections::btree_map::BTreeMap;
use sp_runtime::traits::Hash;
use frame_support::{decl_module, decl_storage, decl_event, dispatch::DispatchResult};
use frame_system::ensure_signed;

// Define the Sequencer Module
pub trait Config: frame_system::Config {
    type Event: From<Event<Self>> + Into<<Self as frame_system::Config>::Event>;
}

decl_storage! {
    trait Store for Module<T: Config> as SequencerModule {
        Transactions get(fn transactions): BTreeMap<T::Hash, u64>; // (Transaction Hash, Delay)
        OrderedTransactions get(fn ordered_transactions): Vec<T::Hash>;
    }
}

decl_event!(
    pub enum Event<T> where Hash = <T as frame_system::Config>::Hash {
        TransactionSequenced(Hash),
        BatchSubmitted,
    }
);

decl_module! {
    pub struct Module<T: Config> for enum Call where origin: T::Origin {
        fn deposit_event() = default;

        #[weight = 10_000]
        pub fn submit_transaction(origin, tx: Vec<u8>) -> DispatchResult {
            let sender = ensure_signed(origin)?;

            // Compute transaction hash
            let tx_hash = T::Hashing::hash(&tx);

            // Compute delay function (simplified for illustration)
            let delay = Self::compute_vdf(tx_hash);

            Transactions::insert(tx_hash, delay);
            Self::deposit_event(RawEvent::TransactionSequenced(tx_hash));

            Ok(())
        }

        #[weight = 50_000]
        pub fn order_transactions(origin) -> DispatchResult {
            let _sender = ensure_signed(origin)?;

            let mut transaction_list: Vec<_> = Transactions::iter().collect();
            transaction_list.sort_by_key(|&(_, delay)| delay);

            let ordered: Vec<T::Hash> = transaction_list.into_iter().map(|(hash, _)| hash).collect();
            OrderedTransactions::put(ordered);

            Self::deposit_event(RawEvent::BatchSubmitted);
            Ok(())
        }
    }
}

// Verifiable Delay Function (VDF) Simulation
impl<T: Config> Module<T> {
    fn compute_vdf(hash: T::Hash) -> u64 {
        let bytes = hash.as_ref();
        let mut delay: u64 = 0;
        for byte in bytes {
            delay += (*byte as u64) % 100;
        }
        delay
    }
}
