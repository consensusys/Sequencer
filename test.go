package sequencer_test

import (
	"testing"
	"github.com/consensusys/Sequencer/sequencer"
	"github.com/stretchr/testify/assert"
)

func TestInitializeSequencer(t *testing.T) {
	config := sequencer.Config{
		RpcEndpoint: "https://l2-network.example.com",
		ChainID: 420,
		MaxTxsPerBlock: 100,
		VdfDifficulty: 2048,
	}
	seq, err := sequencer.NewSequencer(config)
	assert.NoError(t, err, "Sequencer should initialize without error")
	assert.NotNil(t, seq, "Sequencer instance should not be nil")
}

func TestTransactionSubmission(t *testing.T) {
	config := sequencer.Config{
		RpcEndpoint: "https://l2-network.example.com",
		ChainID: 420,
		MaxTxsPerBlock: 100,
		VdfDifficulty: 2048,
	}
	seq, _ := sequencer.NewSequencer(config)

tx := sequencer.Transaction{
		ID:   "tx123",
		Data: "sample transaction",
	}

	err := seq.SubmitTransaction(tx)
	assert.NoError(t, err, "Transaction submission should not return an error")
	assert.Equal(t, 1, len(seq.Mempool()), "Mempool should contain the submitted transaction")
}

func TestBlockCreation(t *testing.T) {
	config := sequencer.Config{
		RpcEndpoint: "https://l2-network.example.com",
		ChainID: 420,
		MaxTxsPerBlock: 2,
		VdfDifficulty: 2048,
	}
	seq, _ := sequencer.NewSequencer(config)
	
	// Submit two transactions
	seq.SubmitTransaction(sequencer.Transaction{ID: "tx1", Data: "data1"})
	seq.SubmitTransaction(sequencer.Transaction{ID: "tx2", Data: "data2"})
	
	block, err := seq.CreateBlock()
	assert.NoError(t, err, "Block creation should not return an error")
	assert.Equal(t, 2, len(block.Transactions), "Block should contain two transactions")
}

func TestFraudProofMechanism(t *testing.T) {
	config := sequencer.Config{
		RpcEndpoint: "https://l2-network.example.com",
		ChainID: 420,
		MaxTxsPerBlock: 100,
		VdfDifficulty: 2048,
		EnableFraudProofs: true,
	}
	seq, _ := sequencer.NewSequencer(config)
	
	// Simulate fraudulent transaction
	tx := sequencer.Transaction{ID: "tx_fraud", Data: "malicious_data"}
	seq.SubmitTransaction(tx)
	
	fraudulent := seq.DetectFraud(tx)
	assert.True(t, fraudulent, "Fraud detection should mark transaction as fraudulent")
}
