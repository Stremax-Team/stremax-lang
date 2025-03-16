package blockchain

import (
	"errors"
	"fmt"
	"reflect"
)

// ContractState represents the state of a smart contract
type ContractState map[string]interface{}

// Event represents a contract event
type Event struct {
	Name   string
	Params map[string]interface{}
}

// Contract represents a smart contract
type Contract struct {
	Name       string
	Owner      Address
	State      ContractState
	Functions  map[string]ContractFunction
	Events     map[string]EventDefinition
	EventLog   []Event
	Blockchain *Blockchain
}

// ContractFunction represents a function in a smart contract
type ContractFunction func(ctx *ContractContext, args ...interface{}) (interface{}, error)

// EventDefinition represents an event definition
type EventDefinition struct {
	Name       string
	Parameters []ParameterDefinition
}

// ParameterDefinition represents a parameter definition
type ParameterDefinition struct {
	Name string
	Type reflect.Type
}

// ContractContext represents the context for a contract function call
type ContractContext struct {
	Contract   *Contract
	Sender     Address
	Value      int64
	Blockchain *Blockchain
}

// NewContract creates a new contract
func NewContract(name string, owner Address, bc *Blockchain) *Contract {
	return &Contract{
		Name:       name,
		Owner:      owner,
		State:      make(ContractState),
		Functions:  make(map[string]ContractFunction),
		Events:     make(map[string]EventDefinition),
		EventLog:   []Event{},
		Blockchain: bc,
	}
}

// RegisterFunction registers a function with the contract
func (c *Contract) RegisterFunction(name string, fn ContractFunction) {
	c.Functions[name] = fn
}

// RegisterEvent registers an event with the contract
func (c *Contract) RegisterEvent(name string, params []ParameterDefinition) {
	c.Events[name] = EventDefinition{
		Name:       name,
		Parameters: params,
	}
}

// Call calls a function on the contract
func (c *Contract) Call(sender Address, functionName string, value int64, args ...interface{}) (interface{}, error) {
	fn, exists := c.Functions[functionName]
	if !exists {
		return nil, fmt.Errorf("function %s does not exist", functionName)
	}

	ctx := &ContractContext{
		Contract:   c,
		Sender:     sender,
		Value:      value,
		Blockchain: c.Blockchain,
	}

	return fn(ctx, args...)
}

// EmitEvent emits an event
func (c *Contract) EmitEvent(name string, params map[string]interface{}) error {
	event, exists := c.Events[name]
	if !exists {
		return fmt.Errorf("event %s does not exist", name)
	}

	// Validate parameters
	for _, param := range event.Parameters {
		value, exists := params[param.Name]
		if !exists {
			return fmt.Errorf("missing parameter %s for event %s", param.Name, name)
		}

		if reflect.TypeOf(value) != param.Type {
			return fmt.Errorf("invalid type for parameter %s in event %s", param.Name, name)
		}
	}

	c.EventLog = append(c.EventLog, Event{
		Name:   name,
		Params: params,
	})

	return nil
}

// Require checks if a condition is true, and reverts the transaction if not
func (ctx *ContractContext) Require(condition bool, message string) error {
	if !condition {
		return errors.New(message)
	}
	return nil
}

// Transfer transfers tokens from the contract to an address
func (ctx *ContractContext) Transfer(to Address, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Create a transaction from the contract to the recipient
	ctx.Blockchain.CreateTransaction(
		Address(ctx.Contract.Name),
		to,
		amount,
		[]byte(fmt.Sprintf("Transfer from contract %s", ctx.Contract.Name)),
	)

	return nil
}

// GetSender returns the sender of the current transaction
func (ctx *ContractContext) GetSender() Address {
	return ctx.Sender
}

// GetValue returns the value sent with the current transaction
func (ctx *ContractContext) GetValue() int64 {
	return ctx.Value
}

// GetState gets a value from the contract state
func (ctx *ContractContext) GetState(key string) interface{} {
	return ctx.Contract.State[key]
}

// SetState sets a value in the contract state
func (ctx *ContractContext) SetState(key string, value interface{}) {
	ctx.Contract.State[key] = value
}

// EmitEvent emits an event
func (ctx *ContractContext) EmitEvent(name string, params map[string]interface{}) error {
	return ctx.Contract.EmitEvent(name, params)
}
