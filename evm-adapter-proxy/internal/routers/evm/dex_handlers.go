package evm

import (
	"encoding/json"
	"fmt"

	"github.com/aviate-labs/agent-go/candid/idl"
	icpDex "github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp/clients/dex"
)

// GetCurrencyPairs handles the eth_getCurrencyPairs RPC method
// Returns all available currency pairs from the DEX canister
func (r *evmRouter) GetCurrencyPairs(_ JSONRPCRequest) (interface{}, error) {
	pairs, err := r.icpClients.Dex.GetCurrencyPairs()
	if err != nil {
		return nil, fmt.Errorf("failed to get currency pairs: %w", err)
	}
	return pairs, nil
}

// MintTokens handles the eth_mintTokens RPC method
// Mints new tokens for a specified recipient
//
// The request must include:
// - currency: The token to mint
// - amount: Amount to mint in hex format
// - recipient: Ethereum address of the recipient
func (r *evmRouter) MintTokens(request JSONRPCRequest) (interface{}, error) {
	var mintReq MintRequest
	params, ok := request.Params.([]interface{})
	if !ok || len(params) == 0 {
		return nil, fmt.Errorf("invalid params for eth_mintTokens")
	}

	paramBytes, err := json.Marshal(params[0])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mint params: %w", err)
	}

	if err := json.Unmarshal(paramBytes, &mintReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mint request: %w", err)
	}

	principalRecipient, err := ConvertEthAddressToICPPrincipal(mintReq.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to convert eth address to principal: %w", err)
	}

	amountInt, err := ConvertHexAmountToBigInt(mintReq.Amount)
	if err != nil {
		return nil, err
	}

	operation := icpDex.MintOperation{
		Currency:  mintReq.Currency,
		Amount:    idl.NewNatFromString(amountInt.String()),
		Recipient: principalRecipient,
	}

	result, err := r.icpClients.Dex.MintTokens(operation)
	if err != nil {
		return nil, fmt.Errorf("failed to mint tokens: %w", err)
	}

	if result.Err != nil {
		return nil, fmt.Errorf("failed to mint tokens: %s", *result.Err)
	}

	return true, nil
}

// BurnTokens handles the eth_burnTokens RPC method
// Burns tokens from a specified owner's balance
//
// The request must include:
// - currency: The token to burn
// - amount: Amount to burn in hex format
// - owner: Ethereum address of the token owner
func (r *evmRouter) BurnTokens(request JSONRPCRequest) (interface{}, error) {
	var burnReq BurnRequest
	params, ok := request.Params.([]interface{})
	if !ok || len(params) == 0 {
		return nil, fmt.Errorf("invalid params for eth_burnTokens")
	}

	paramBytes, err := json.Marshal(params[0])
	if err != nil {
		return nil, fmt.Errorf("failed to marshal burn params: %w", err)
	}

	if err := json.Unmarshal(paramBytes, &burnReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal burn request: %w", err)
	}

	principalOwner, err := ConvertEthAddressToICPPrincipal(burnReq.Owner)
	if err != nil {
		return nil, fmt.Errorf("failed to convert eth address to principal: %w", err)
	}

	amountInt, err := ConvertHexAmountToBigInt(burnReq.Amount)
	if err != nil {
		return nil, err
	}

	operation := icpDex.BurnOperation{
		Currency: burnReq.Currency,
		Amount:   idl.NewNatFromString(amountInt.String()),
		Owner:    principalOwner,
	}

	result, err := r.icpClients.Dex.BurnTokens(operation)
	if err != nil {
		return nil, fmt.Errorf("failed to burn tokens: %w", err)
	}

	if result.Err != nil {
		return nil, fmt.Errorf("failed to burn tokens: %s", *result.Err)
	}

	return true, nil
}
