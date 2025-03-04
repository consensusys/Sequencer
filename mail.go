func main() {
	sequencer := NewSequencer()

	// Submit Transactions
	sequencer.SubmitTransaction([]byte("Tx1: Alice -> Bob 10 ETH"))
	sequencer.SubmitTransaction([]byte("Tx2: Charlie -> Dave 5 ETH"))
	sequencer.SubmitTransaction([]byte("Tx3: Eve -> Frank 2 ETH"))

	// Order and Submit to Layer 1
	sequencer.SubmitBatchToLayer1()
}
