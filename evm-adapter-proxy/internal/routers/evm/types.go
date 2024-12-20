package evm

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      interface{} `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      interface{}   `json:"id"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Block struct {
	Number           string   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	Nonce            string   `json:"nonce"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	LogsBloom        string   `json:"logsBloom"`
	TransactionsRoot string   `json:"transactionsRoot"`
	StateRoot        string   `json:"stateRoot"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Miner            string   `json:"miner"`
	Difficulty       string   `json:"difficulty"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	ExtraData        string   `json:"extraData"`
	Size             string   `json:"size"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Timestamp        string   `json:"timestamp"`
	Transactions     []string `json:"transactions"`
	Uncles           []string `json:"uncles"`
}

type Log struct {
	Address     string   `json:"address"`
	Topics      []string `json:"topics"`
	Data        string   `json:"data"`
	BlockNumber string   `json:"blockNumber"`
	TxHash      string   `json:"transactionHash"`
	TxIndex     string   `json:"transactionIndex"`
	BlockHash   string   `json:"blockHash"`
	LogIndex    string   `json:"logIndex"`
	Removed     bool     `json:"removed"`
}

type LogData struct {
	Operation string         `json:"operation"`
	Detail    LogDataDetails `json:"detail"`
}

type LogDataDetails struct {
	Map []LogDataField `json:"map"`
}

type LogDataField struct {
	Field0 string     `json:"field0"`
	Field1 FieldValue `json:"field1"`
}

type FieldValue struct {
	Text *string `json:"text,omitempty"`
	Nat  *string `json:"nat,omitempty"`
}

type MintRequest struct {
	Currency  string      `json:"currency"`
	Amount    hexutil.Big `json:"amount"`
	Recipient string      `json:"recipient"`
}

type BurnRequest struct {
	Currency string      `json:"currency"`
	Amount   hexutil.Big `json:"amount"`
	Owner    string      `json:"owner"`
}

func hexToDecimal(hex string) (string, error) {
	hex = strings.TrimPrefix(hex, "0x")

	switch hex {
	case "latest", "pending", "earliest", "finalized", "safe":
		return hex, nil
	}

	n, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(n, 10), nil
}

// ConvertHexAmountToBigInt converts a hex amount to a big.Int
// Returns nil if the conversion fails
func ConvertHexAmountToBigInt(amount hexutil.Big) (*big.Int, error) {
	amountStr := strings.TrimPrefix(amount.String(), "0x")
	amountInt, ok := new(big.Int).SetString(amountStr, 16)
	if !ok {
		return nil, fmt.Errorf("failed to parse amount: %s", amount.String())
	}
	return amountInt, nil
}
