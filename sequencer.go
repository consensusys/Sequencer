package sequencer

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Transaction represents a Layer 2 transaction
type Transaction struct {
	Hash       string
	Data       []byte
	Commitment string
	Delay      int64
	Timestamp  time.Time
}

// Sequencer manages the ordering of transactions
type Sequencer struct {
	transactions []Transaction
}

// NewSequencer initializes a new sequencer
func NewSequencer() *Sequencer {
	return &Sequencer{
		transactions: []Transaction{},
	}
}

// SubmitTransaction adds a new transaction to the sequencer
func (s *Sequencer) SubmitTransaction(data []byte) {
	hash := computeHash(data)
	commitment := computeCommitment(hash)
	delay := computeVDF(hash)

	tx := Transaction{
		Hash:       hash,
		Data:       data,
		Commitment: commitment,
		Delay:      delay,
		Timestamp:  time.Now(),
	}

	s.transactions = append(s.transactions, tx)
	fmt.Printf("Transaction submitted: %s\n", hash)
}

// OrderTransactions orders transactions based on VDF delay
func (s *Sequencer) OrderTransactions() []Transaction {
	sort.Slice(s.transactions, func(i, j int) bool {
		return s.transactions[i].Delay < s.transactions[j].Delay
	})

	fmt.Println("Transactions ordered based on delay")
	return s.transactions
}

// Compute Hash (Cryptographic Commitment)
func computeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// Compute Commitment using a secondary hash function
func computeCommitment(hash string) string {
	commitmentHash := sha256.Sum256([]byte(hash))
	return hex.EncodeToString(commitmentHash[:])
}

// Compute Verifiable Delay Function (VDF) Simulation
func computeVDF(hash string) int64 {
	var delay int64
	for _, b := range hash {
		delay += int64(b) % 100
	}
	return delay
}

// Finalize and Batch Transactions
func (s *Sequencer) SubmitBatchToLayer1() {
	orderedTxs := s.OrderTransactions()

	var batch bytes.Buffer
	for _, tx := range orderedTxs {
		batch.Write(tx.Data)
		batch.WriteString("\n")
	}

	fmt.Println("Batch submitted to Layer 1:", batch.String())
}
