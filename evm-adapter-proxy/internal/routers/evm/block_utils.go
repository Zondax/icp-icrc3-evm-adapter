package evm

import (
	"bytes"
	"fmt"

	"reflect"

	icpLogger "github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp/clients/logger"
)

// mapBlockToEVMBlock converts an ICRC-3 block to an EVM-compatible block format
//
// Parameters:
//   - blockValue: The ICRC-3 block value to convert
//
// Returns:
//   - Block: The converted EVM-compatible block
//   - error: Any error that occurred during conversion
//
// Note: This is a PoC implementation that fills many fields with placeholder values
// as they don't have direct equivalents in ICRC-3
func mapBlockToEVMBlock(blockValue icpLogger.Value) (Block, error) {
	var icrcBlock icpLogger.Block
	err := decodeValue(blockValue, &icrcBlock)
	if err != nil {
		return Block{}, fmt.Errorf("failed to decode block: %w", err)
	}

	return Block{
		Number:       fmt.Sprintf("0x%x", icrcBlock.Id.BigInt()),
		Hash:         fmt.Sprintf("0x%x", icrcBlock.Hash),
		ParentHash:   fmt.Sprintf("0x%x", icrcBlock.Phash),
		Timestamp:    fmt.Sprintf("0x%x", icrcBlock.Ts/1e9), // Convert nanoseconds to seconds to prevent overflow
		Transactions: []string{},                            // we don't have transactions in this PoC
		// Fill other fields with placeholder values
		TransactionsRoot: "0x0000000000000000000000000000000000000000000000000000000000000000",
		ReceiptsRoot:     "0x0000000000000000000000000000000000000000000000000000000000000000",
		StateRoot:        "0x0000000000000000000000000000000000000000000000000000000000000000",
		Nonce:            "0x0000000000000000",
		Sha3Uncles:       "",
		LogsBloom:        "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Miner:            "0x0000000000000000000000000000000000000000",
		Difficulty:       "0x0",
		TotalDifficulty:  "0x0",
		ExtraData:        "0x",
		Size:             "0x0",
		GasLimit:         "0x0",
		GasUsed:          "0x0",
		Uncles:           []string{},
	}, nil
}

// decodeValue decodes an ICRC-3 Value into a Block structure using reflection
//
// Parameters:
//   - value: The ICRC-3 Value to decode
//   - target: Pointer to the target Block structure
//
// Returns:
//   - error: Any error that occurred during decoding
//
// The function uses reflection to map ICRC-3 fields to their corresponding
// Block structure fields. It performs type validation for each field.
func decodeValue(value icpLogger.Value, target interface{}) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer")
	}

	block, ok := v.Interface().(*icpLogger.Block)
	if !ok {
		return fmt.Errorf("target must be a pointer to icp.Block")
	}

	if value.Map == nil {
		return fmt.Errorf("expected Map value")
	}

	for _, field := range *value.Map {
		switch field.Field0 {
		case "id":
			if field.Field1.Nat == nil {
				return fmt.Errorf("invalid id field")
			}
			block.Id = *field.Field1.Nat
		case "hash":
			if field.Field1.Blob == nil {
				return fmt.Errorf("invalid hash field")
			}
			block.Hash = *field.Field1.Blob
		case "phash":
			if field.Field1.Blob == nil {
				return fmt.Errorf("invalid phash field")
			}
			block.Phash = *field.Field1.Blob
		case "ts":
			if field.Field1.Nat == nil {
				return fmt.Errorf("invalid ts field")
			}
			block.Ts = field.Field1.Nat.BigInt().Uint64()
		case "entries":
			if field.Field1.Array == nil {
				return fmt.Errorf("invalid entries field")
			}
			block.Entries = *field.Field1.Array
		case "finalized":
			if field.Field1.Text == nil {
				return fmt.Errorf("invalid finalized field")
			}
			block.Finalized = *field.Field1.Text == "true"
		}
	}

	return nil
}

// decodeCertificateData decodes a ULEB128-encoded block number from certificate data
//
// Parameters:
//   - data: The raw certificate data containing the encoded block number
//
// Returns:
//   - uint64: The decoded block number
//   - error: Any error that occurred during decoding
//
// This function implements ULEB128 (Unsigned Little-Endian Base 128) decoding
// which is used in ICRC-3 certificates to encode block numbers.
func decodeCertificateData(data []byte) (uint64, error) {
	reader := bytes.NewReader(data)
	var blockNumber uint64
	shift := uint(0)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("failed to read byte: %w", err)
		}

		blockNumber |= uint64(b&0x7F) << shift

		if (b & 0x80) == 0 {
			break
		}
		shift += 7

		if shift >= 64 {
			return 0, fmt.Errorf("ULEB128 encoding is too large")
		}
	}

	return blockNumber, nil
}
