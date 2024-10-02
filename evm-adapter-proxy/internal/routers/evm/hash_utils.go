package evm

import (
	"encoding/hex"
	"fmt"

	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp"
	"golang.org/x/crypto/sha3"
)

// Note: All methods in this file are placeholders and not suitable for production use.
// They provide simplified implementations for demonstration purposes only.

func getBlockHash(block icp.Value) (string, error) {
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

func calculateTransactionHash(tx icp.Value) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(fmt.Sprintf("%v", tx)))
	return "0x" + hex.EncodeToString(hash.Sum(nil))
}
