package evm

import (
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	icpLogger "github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp/clients/logger"

	"github.com/aviate-labs/agent-go/candid/idl"
)

// EthChainID implements the eth_chainId RPC method
// Returns the current chain ID in hexadecimal format
func (r *evmRouter) EthChainID(_ JSONRPCRequest) (interface{}, error) {
	chainID, err := r.icpClients.Logger.ChainId()
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}
	if chainID == nil {
		return nil, fmt.Errorf("chain ID is nil")
	}

	chainIDInt, err := strconv.ParseUint(*chainID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid chain ID format: %w", err)
	}

	evmChainID := fmt.Sprintf("0x%x", chainIDInt)
	return evmChainID, nil
}

// EthBlockNumber implements the eth_blockNumber RPC method
// Returns the latest block number in hexadecimal format
func (r *evmRouter) EthBlockNumber(_ JSONRPCRequest) (interface{}, error) {
	tipCert, err := r.icpClients.Logger.Icrc3GetTipCertificate()
	if err != nil {
		return nil, fmt.Errorf("failed to get tip certificate: %w", err)
	}
	if tipCert == nil || *tipCert == nil {
		return nil, fmt.Errorf("no tip certificate found")
	}

	certificate := (*tipCert).Certificate
	if certificate == nil {
		return nil, fmt.Errorf("certificate data is nil")
	}

	blockNumber, err := decodeCertificateData(certificate)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tip certificate: %w", err)
	}

	ethBlockNumber := fmt.Sprintf("0x%x", blockNumber)
	return ethBlockNumber, nil
}

// EthAccounts Return an empty array as we don't manage accounts
func (r *evmRouter) EthAccounts(_ JSONRPCRequest) (interface{}, error) {
	return []string{}, nil
}

// EthNetVersion implements the net_version RPC method
// Returns the current network ID
func (r *evmRouter) EthNetVersion(_ JSONRPCRequest) (interface{}, error) {
	netVersion, err := r.icpClients.Logger.NetVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get net version: %w", err)
	}
	return netVersion, nil
}

// EthGetBlockByNumber implements the eth_getBlockByNumber RPC method
// Retrieves a block by its number
//
// Parameters:
// - blockNumberHex: Block number in hex or "latest"
// Note: This is a PoC implementation and should be enhanced for production use
func (r *evmRouter) EthGetBlockByNumber(request JSONRPCRequest) (interface{}, error) {
	params, ok := request.Params.([]interface{})
	if !ok || len(params) < 1 {
		return nil, fmt.Errorf("invalid params for eth_getBlockByNumber")
	}

	blockNumberHex, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid block number")
	}

	var blockNumber string
	var err error
	switch blockNumberHex {
	case "latest", "safe", "finalized":
		latestBlockHex, err := r.EthBlockNumber(JSONRPCRequest{})
		if err != nil {
			return nil, fmt.Errorf("failed to get latest block number: %w", err)
		}
		blockNumberHex = latestBlockHex.(string)
	}

	blockNumber, err = hexToDecimal(blockNumberHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse block number: %w", err)
	}

	blocksArgs := icpLogger.GetBlocksArgs{
		Start:  idl.NewNatFromString(blockNumber),
		Length: idl.NewNatFromString("1"),
	}
	result, err := r.icpClients.Logger.Icrc3GetBlocks(blocksArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %w", err)
	}

	if result.LogLength.BigInt().Uint64() == 0 {
		return nil, fmt.Errorf("block not found")
	}

	if len(result.Blocks) == 0 {
		return nil, fmt.Errorf("block not found")
	}

	firstBlock := result.Blocks[0]
	evmBlock, err := mapBlockToEVMBlock(firstBlock.Block)
	if err != nil {
		return nil, fmt.Errorf("failed to map ICRC3 block to EVM block: %w", err)
	}

	return evmBlock, nil
}

// EthGetBlockByHash implements the eth_getBlockByHash RPC method
// Retrieves a block by its hash
//
// Parameters:
// - blockHash: The hash of the block to retrieve
func (r *evmRouter) EthGetBlockByHash(request JSONRPCRequest) (interface{}, error) {
	params, ok := request.Params.([]interface{})
	if !ok || len(params) < 1 {
		return nil, fmt.Errorf("invalid params for eth_getBlockByHash")
	}

	requestedBlockHash, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid block hash")
	}

	blocksArgs := icpLogger.GetBlocksArgs{
		Start:  idl.NewNatFromString("0"), // Initialize as needed
		Length: idl.NewNatFromString("1"),
	}
	latestBlockResult, err := r.icpClients.Logger.Icrc3GetBlocks(blocksArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}

	if len(latestBlockResult.Blocks) == 0 {
		return nil, fmt.Errorf("no blocks found")
	}

	currentBlockNumber := latestBlockResult.Blocks[0].Id.BigInt().Uint64()

	for {
		blocksArgs := icpLogger.GetBlocksArgs{
			Start:  idl.NewNatFromString(fmt.Sprintf("%d", currentBlockNumber)),
			Length: idl.NewNatFromString("1"),
		}
		blocksResult, err := r.icpClients.Logger.Icrc3GetBlocks(blocksArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to get block %d: %w", currentBlockNumber, err)
		}

		if len(blocksResult.Blocks) == 0 {
			return nil, fmt.Errorf("block %d not found", currentBlockNumber)
		}

		blockResult := blocksResult.Blocks[0]

		blockHash, err := getBlockHash(blockResult.Block)
		if err != nil {
			return nil, fmt.Errorf("failed to get block hash: %w", err)
		}

		if blockHash == requestedBlockHash {
			return mapBlockToEVMBlock(blockResult.Block)
		}

		if currentBlockNumber <= 0 {
			break
		}

		currentBlockNumber--
	}

	return nil, fmt.Errorf("block with hash %s not found", requestedBlockHash)
}

// EthGetLogs implements the eth_getLogs RPC method
// Retrieves logs matching the provided filter criteria
//
// The filter can include:
// - fromBlock: Start block number
// - toBlock: End block number
// - address: Contract address to filter
// - blockHash: Specific block to get logs from
func (r *evmRouter) EthGetLogs(request JSONRPCRequest) (interface{}, error) {
	filter, err := extractFilterFromParams(request.Params)
	if err != nil {
		return nil, err
	}

	fromBlock, toBlock, err := extractBlockRange(filter)
	if err != nil {
		return nil, err
	}

	address, err := extractAddress(filter)
	if err != nil {
		return nil, err
	}

	filterBlockHash := extractBlockHash(filter)

	return r.getLogsByFilter(fromBlock, toBlock, address, filterBlockHash)
}

// Helper functions

// extractFilterFromParams extracts and validates the filter parameters from the RPC request
func extractFilterFromParams(params interface{}) (map[string]interface{}, error) {
	paramSlice, ok := params.([]interface{})
	if !ok || len(paramSlice) < 1 {
		return nil, fmt.Errorf("invalid params for eth_getLogs")
	}

	filter, ok := paramSlice[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid filter parameters")
	}

	return filter, nil
}

// extractBlockHash extracts the blockHash parameter from the filter
func extractBlockHash(filter map[string]interface{}) string {
	if blockHash, ok := filter["blockHash"].(string); ok {
		return blockHash
	}
	return ""
}

// extractBlockRange extracts and validates the block range from the filter
func extractBlockRange(filter map[string]interface{}) (uint64, uint64, error) {
	fromBlock, err := parseBlockParam(filter["fromBlock"], 0)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid fromBlock: %w", err)
	}

	toBlock, err := parseBlockParam(filter["toBlock"], ^uint64(0)) // Max uint64 value
	if err != nil {
		return 0, 0, fmt.Errorf("invalid toBlock: %w", err)
	}

	return fromBlock, toBlock, nil
}

// parseBlockParam parses a block parameter which can be a number or "latest"
func parseBlockParam(blockParam interface{}, defaultValue uint64) (uint64, error) {
	if blockParam == nil {
		return defaultValue, nil
	}

	switch v := blockParam.(type) {
	case string:
		if v == "latest" {
			return defaultValue, nil
		}
		return strconv.ParseUint(v, 0, 64)
	case float64:
		return uint64(v), nil
	default:
		return 0, fmt.Errorf("invalid block number format")
	}
}

// getLogsByFilter retrieves logs matching the specified filter criteria
//
// Parameters:
// - fromBlock: Start of block range
// - toBlock: End of block range
// - address: Optional address to filter logs
// - filterBlockHash: Optional specific block hash to get logs from
func (r *evmRouter) getLogsByFilter(fromBlock, toBlock uint64, address, filterBlockhash string) ([]Log, error) {
	var logs []Log
	// TODO: Just for PoC
	batchSize := uint64(10)

	latestBlock, err := r.getLatestBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block number: %w", err)
	}
	if toBlock > latestBlock {
		toBlock = latestBlock
	}

	if fromBlock > toBlock {
		return nil, fmt.Errorf("fromBlock (%d) is greater than toBlock (%d)", fromBlock, toBlock)
	}

	for start := fromBlock; start <= toBlock; start += batchSize {
		end := start + batchSize - 1
		if end > toBlock {
			end = toBlock
		}

		blocksArgs := icpLogger.GetBlocksArgs{
			Start:  idl.NewNatFromString(fmt.Sprintf("%d", start)),
			Length: idl.NewNatFromString(fmt.Sprintf("%d", end-start+1)),
		}

		result, err := r.icpClients.Logger.Icrc3GetBlocks(blocksArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to get blocks: %w", err)
		}

		for _, blockInfo := range result.Blocks {
			blockLogs, err := extractLogsFromBlock(blockInfo.Block, address, filterBlockhash)
			if err != nil {
				return nil, fmt.Errorf("failed to extract logs from block: %w", err)
			}

			logs = append(logs, blockLogs...)
		}

		if uint64(len(result.Blocks)) < batchSize {
			break
		}
	}

	return logs, nil
}
func (r *evmRouter) getLatestBlockNumber() (uint64, error) {
	latestBlockHex, err := r.EthBlockNumber(JSONRPCRequest{})
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block number: %w", err)
	}

	latestBlockStr, ok := latestBlockHex.(string)
	if !ok {
		return 0, fmt.Errorf("invalid block number format")
	}

	latestBlock, err := strconv.ParseUint(latestBlockStr[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block number: %w", err)
	}

	return latestBlock, nil
}

// extractLogsFromBlock extracts EVM-compatible logs from an ICRC-3 block
//
// Parameters:
// - block: The ICRC-3 block to extract logs from
// - address: Optional address filter
// - filterBlockHash: Optional block hash filter
func extractLogsFromBlock(block icpLogger.Value, address, filterBlockHash string) ([]Log, error) {
	var logs []Log

	blockMap := block.Map
	if blockMap == nil {
		return nil, fmt.Errorf("invalid block format: Map is nil")
	}

	var id *idl.Nat
	var entries []icpLogger.Value
	for _, entry := range *blockMap {
		switch entry.Field0 {
		case "id":
			if nat := entry.Field1.Nat; nat != nil {
				id = nat
			}
		case "entries":
			if entriesArray := entry.Field1.Array; entriesArray != nil {
				entries = *entriesArray
			}
		}
	}

	for i, entryValue := range entries {
		entryMap := entryValue.Map
		if entryMap == nil {
			continue
		}

		var logEntry struct {
			Timestamp uint64
			Operation string
			Details   icpLogger.Value
			Caller    string
		}

		for _, field := range *entryMap {
			switch field.Field0 {
			case "timestamp":
				if nat := field.Field1.Nat; nat != nil {
					logEntry.Timestamp = nat.BigInt().Uint64()
				}
			case "operation":
				if text := field.Field1.Text; text != nil {
					logEntry.Operation = *text
				}
			case "details":
				logEntry.Details = field.Field1
			case "caller":
				if text := field.Field1.Text; text != nil {
					logEntry.Caller = *text
				}
			}
		}

		ethAddress, err := convertICPToEthAddress(logEntry.Caller) //nolint
		if err != nil {
			return nil, fmt.Errorf("failed to convert ICP address to ETH address: %w", err)
		}

		if address != "" && ethAddress != address {
			continue
		}

		blockHash, err := getBlockHash(block)
		if err != nil {
			return nil, fmt.Errorf("failed to get block hash: %w", err)
		}

		blockNumber := "0x0"
		if id != nil {
			if bigInt := id.BigInt(); bigInt != nil {
				blockNumber = fmt.Sprintf("0x%x", bigInt.Uint64())
			}
		}

		logData := LogData{
			Operation: logEntry.Operation,
			Detail: LogDataDetails{
				Map: extractDetailMap(logEntry.Details),
			},
		}

		logDataJSON, err := json.Marshal(logData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal log data: %w", err)
		}

		log := Log{
			Address:     ethAddress,
			Topics:      []string{"0x0000000000000000000000000000000000000000000000000000000000000000"},
			Data:        fmt.Sprintf("0x%x", logDataJSON),
			BlockNumber: blockNumber,
			BlockHash:   blockHash,
			TxHash:      calculateTransactionHash(block),
			TxIndex:     fmt.Sprintf("0x%x", i),
			LogIndex:    fmt.Sprintf("0x%x", i),
			Removed:     false,
		}

		if filterBlockHash != "" && blockHash != filterBlockHash {
			continue
		}

		logs = append(logs, log)
	}

	return logs, nil
}

// extractDetailMap converts ICRC-3 log details to EVM log format
func extractDetailMap(details icpLogger.Value) []LogDataField {
	var fields []LogDataField
	if details.Map == nil {
		return fields
	}

	for _, entry := range *details.Map {
		field := LogDataField{
			Field0: entry.Field0,
			Field1: FieldValue{},
		}

		switch {
		case entry.Field1.Text != nil:
			field.Field1.Text = entry.Field1.Text
		case entry.Field1.Nat != nil:
			natStr := entry.Field1.Nat.String()
			field.Field1.Nat = &natStr
		}

		fields = append(fields, field)
	}

	return fields
}

// extractAddress extracts the address filter from the filter parameters
func extractAddress(filter map[string]interface{}) (string, error) {
	if addr, ok := filter["address"]; ok {
		if addrStr, ok := addr.(string); ok {
			return addrStr, nil
		}
		return "", fmt.Errorf("invalid address format")
	}
	return "", nil // No address filter
}

// Only for POC purposes
func convertICPToEthAddress(icpAddress string) (string, error) {
	icpAddress = strings.TrimSuffix(icpAddress, "-fae")

	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(icpAddress))
	if err != nil {
		return "", fmt.Errorf("failed to decode ICP address: %w", err)
	}

	if len(decoded) < 20 {
		padding := make([]byte, 20-len(decoded))
		decoded = append(padding, decoded...) //nolint
	} else if len(decoded) > 20 {
		decoded = decoded[len(decoded)-20:]
	}

	ethAddress := "0x" + hex.EncodeToString(decoded)

	return ethAddress, nil
}
