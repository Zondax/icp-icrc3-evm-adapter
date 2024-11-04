package evm

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// Web3ClientVersion implements the web3_clientVersion RPC method
// Returns the current client version
//
// Returns:
//   - string: The client version string
//   - error: Always returns nil for this implementation
func (r *evmRouter) Web3ClientVersion(_ JSONRPCRequest) (interface{}, error) {
	return "EVM-Adapter/v0.1.0", nil
}

// Web3Sha3 implements the web3_sha3 RPC method
// Returns Keccak-256 hash of the given data
//
// Parameters in request:
//   - input: The data to hash (hex string)
//
// Returns:
//   - string: The hash in hex format
//   - error: Any error that occurred during processing
func (r *evmRouter) Web3Sha3(request JSONRPCRequest) (interface{}, error) {
	params, ok := request.Params.([]interface{})
	if !ok || len(params) == 0 {
		return nil, fmt.Errorf("invalid params for web3_sha3")
	}

	input, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid input for web3_sha3")
	}

	if len(input) > 2 && input[:2] == "0x" {
		input = input[2:]
	}

	data, err := hex.DecodeString(input)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}

	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return "0x" + hex.EncodeToString(hash.Sum(nil)), nil
}
