package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Address represents a blockchain address
type Address string

// Transaction represents a blockchain transaction
type Transaction struct {
	From      Address
	To        Address
	Amount    int64
	Timestamp time.Time
	Data      []byte
	Hash      string
}

// Block represents a block in the blockchain
type Block struct {
	Index        int
	Timestamp    time.Time
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Nonce        int
}

// Blockchain represents a blockchain
type Blockchain struct {
	Chain               []*Block
	PendingTransactions []Transaction
	Difficulty          int
	Contracts           map[Address]*SmartContract
}

// SmartContract represents a smart contract
type SmartContract struct {
	Address    Address
	Owner      Address
	Code       []byte
	State      map[string]interface{}
	Functions  map[string]func([]interface{}) interface{}
	Events     map[string]func([]interface{})
	Deployed   bool
	DeployTime time.Time
}

// New creates a new blockchain
func New() *Blockchain {
	bc := &Blockchain{
		Chain:               []*Block{},
		PendingTransactions: []Transaction{},
		Difficulty:          4, // Arbitrary difficulty
		Contracts:           make(map[Address]*SmartContract),
	}
	
	// Create the genesis block
	bc.createGenesisBlock()
	
	return bc
}

// createGenesisBlock creates the genesis block
func (bc *Blockchain) createGenesisBlock() {
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now(),
		Transactions: []Transaction{},
		PrevHash:     "0",
		Nonce:        0,
	}
	
	genesisBlock.Hash = bc.calculateHash(genesisBlock)
	bc.Chain = append(bc.Chain, genesisBlock)
}

// calculateHash calculates the hash of a block
func (bc *Blockchain) calculateHash(block *Block) string {
	record := fmt.Sprintf("%d%s%s%d", 
		block.Index, 
		block.Timestamp.String(), 
		block.PrevHash, 
		block.Nonce,
	)
	
	// Add transaction data to the record
	for _, tx := range block.Transactions {
		record += tx.Hash
	}
	
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	
	return hex.EncodeToString(hashed)
}

// CreateTransaction creates a new transaction
func (bc *Blockchain) CreateTransaction(from, to Address, amount int64, data []byte) Transaction {
	tx := Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Timestamp: time.Now(),
		Data:      data,
	}
	
	// Calculate transaction hash
	h := sha256.New()
	record := fmt.Sprintf("%s%s%d%s%s", 
		from, 
		to, 
		amount, 
		tx.Timestamp.String(),
		data,
	)
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	tx.Hash = hex.EncodeToString(hashed)
	
	bc.PendingTransactions = append(bc.PendingTransactions, tx)
	
	return tx
}

// MineBlock mines a new block
func (bc *Blockchain) MineBlock(minerAddress Address) *Block {
	lastBlock := bc.GetLastBlock()
	newBlock := &Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    time.Now(),
		Transactions: bc.PendingTransactions,
		PrevHash:     lastBlock.Hash,
		Nonce:        0,
	}
	
	// Add mining reward
	bc.CreateTransaction(Address("SYSTEM"), minerAddress, 1, []byte("Mining Reward"))
	
	// Mine the block (find a hash with the required difficulty)
	bc.mineBlockWithProofOfWork(newBlock)
	
	// Add the block to the chain
	bc.Chain = append(bc.Chain, newBlock)
	
	// Reset pending transactions
	bc.PendingTransactions = []Transaction{}
	
	return newBlock
}

// mineBlockWithProofOfWork mines a block with proof of work
func (bc *Blockchain) mineBlockWithProofOfWork(block *Block) {
	target := ""
	for i := 0; i < bc.Difficulty; i++ {
		target += "0"
	}
	
	for {
		block.Hash = bc.calculateHash(block)
		if block.Hash[:bc.Difficulty] == target {
			fmt.Printf("Block mined: %s\n", block.Hash)
			break
		}
		block.Nonce++
	}
}

// GetLastBlock returns the last block in the chain
func (bc *Blockchain) GetLastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

// IsChainValid checks if the blockchain is valid
func (bc *Blockchain) IsChainValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		prevBlock := bc.Chain[i-1]
		
		// Check if the current block's hash is valid
		if currentBlock.Hash != bc.calculateHash(currentBlock) {
			return false
		}
		
		// Check if the current block points to the previous block's hash
		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}
	
	return true
}

// GetBalance returns the balance of an address
func (bc *Blockchain) GetBalance(address Address) int64 {
	balance := int64(0)
	
	for _, block := range bc.Chain {
		for _, tx := range block.Transactions {
			if tx.From == address {
				balance -= tx.Amount
			}
			if tx.To == address {
				balance += tx.Amount
			}
		}
	}
	
	return balance
}

// DeployContract deploys a smart contract to the blockchain
func (bc *Blockchain) DeployContract(owner Address, code []byte) (Address, error) {
	// Generate a new address for the contract
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s%s%d", owner, code, time.Now().UnixNano())))
	hashed := h.Sum(nil)
	contractAddress := Address(hex.EncodeToString(hashed)[:40])
	
	// Create the contract
	contract := &SmartContract{
		Address:    contractAddress,
		Owner:      owner,
		Code:       code,
		State:      make(map[string]interface{}),
		Functions:  make(map[string]func([]interface{}) interface{}),
		Events:     make(map[string]func([]interface{})),
		Deployed:   true,
		DeployTime: time.Now(),
	}
	
	// Add the contract to the blockchain
	bc.Contracts[contractAddress] = contract
	
	// Create a deployment transaction
	bc.CreateTransaction(owner, contractAddress, 0, code)
	
	return contractAddress, nil
}

// Contracts stores the smart contracts deployed on the blockchain
var Contracts = make(map[Address]*SmartContract)

// GetContract retrieves a contract by its address
func (bc *Blockchain) GetContract(address Address) (*SmartContract, bool) {
	contract, ok := bc.Contracts[address]
	return contract, ok
}

// CallContract calls a function on a smart contract
func (bc *Blockchain) CallContract(from Address, to Address, functionName string, args []interface{}) (interface{}, error) {
	// Get the contract
	contract, ok := bc.GetContract(to)
	if !ok {
		return nil, fmt.Errorf("contract not found at address %s", to)
	}
	
	// Check if the function exists
	function, ok := contract.Functions[functionName]
	if !ok {
		return nil, fmt.Errorf("function %s not found in contract", functionName)
	}
	
	// Call the function
	result := function(args)
	
	// Create a transaction for the function call
	data := []byte(fmt.Sprintf("%s(%v)", functionName, args))
	bc.CreateTransaction(from, to, 0, data)
	
	return result, nil
}

// EmitEvent emits an event from a smart contract
func (bc *Blockchain) EmitEvent(contract Address, eventName string, args []interface{}) error {
	// Get the contract
	c, ok := bc.GetContract(contract)
	if !ok {
		return fmt.Errorf("contract not found at address %s", contract)
	}
	
	// Check if the event exists
	event, ok := c.Events[eventName]
	if !ok {
		return fmt.Errorf("event %s not found in contract", eventName)
	}
	
	// Emit the event
	event(args)
	
	// Log the event
	fmt.Printf("Event emitted: %s from contract %s with args %v\n", eventName, contract, args)
	
	return nil
} 