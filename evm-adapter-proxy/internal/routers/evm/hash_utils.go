package evm

import (
	"encoding/hex"
	"fmt"

	icpLogger "github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp/clients/logger"
	"golang.org/x/crypto/sha3"
)

// Note: All methods in this file are placeholders and not suitable for production use.
// They provide simplified implementations for demonstration purposes only.

// getBlockHash extracts the hash from an ICRC-3 block value
//
// Parameters:
//   - block: The ICRC-3 block value containing the hash
//
// Returns:
//   - string: The block hash in EVM format (0x-prefixed hex)
//   - error: Any error that occurred during extraction
//
// Note: This is a PoC implementation that assumes a specific block structure
func getBlockHash(block icpLogger.Value) (string, error) {
	blockMap := block.Map
	if blockMap == nil {
		return "", fmt.Errorf("invalid block format: Map is nil")
	}

	for _, field := range *blockMap {
		if field.Field0 == "hash" {
			if field.Field1.Blob != nil {
				return fmt.Sprintf("0x%x", *field.Field1.Blob), nil
			}
		}
	}

	return "", fmt.Errorf("block hash not found")
}

// calculateTransactionHash generates a Keccak-256 hash for a transaction
//
// Parameters:
//   - tx: The ICRC-3 value representing the transaction
//
// Returns:
//   - string: The transaction hash in EVM format (0x-prefixed hex)
//
// Note: This is a simplified implementation for PoC purposes.
// In production, proper transaction serialization would be required.
func calculateTransactionHash(tx icpLogger.Value) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(fmt.Sprintf("%v", tx)))
	return "0x" + hex.EncodeToString(hash.Sum(nil))
}
