package evm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zondax/golem/pkg/logger"
	"github.com/zondax/golem/pkg/zrouter"
	"github.com/zondax/golem/pkg/zrouter/domain"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp"
)

// Note: This pkg contains implementations for Proof of Concept (POC) purposes.
// These methods should be optimized and properly implemented for production use.

type methodHandler func(JSONRPCRequest) (interface{}, error)

type EVMRouter interface {
	HandleRPCRequest(ctx zrouter.Context) (domain.ServiceResponse, error)
}

type evmRouter struct {
	methodHandlers       map[string]methodHandler
	icpClients           *icp.Clients
	arrayResponseMethods map[string]bool
}

func (r *evmRouter) initMethodHandlers() {
	r.methodHandlers = map[string]methodHandler{
		// Standard Ethereum JSON-RPC methods
		"eth_chainId":          r.EthChainID,
		"net_version":          r.EthNetVersion,
		"eth_getBlockByNumber": r.EthGetBlockByNumber,
		"eth_getBlockByHash":   r.EthGetBlockByHash,
		"eth_getLogs":          r.EthGetLogs,
		"eth_blockNumber":      r.EthBlockNumber,
		"eth_accounts":         r.EthAccounts,
		"web3_clientVersion":   r.Web3ClientVersion,
		"web3_sha3":            r.Web3Sha3,
		"net_listening":        r.NetListening,
		"net_peerCount":        r.NetPeerCount,

		// Custom DEX methods
		"eth_getCurrencyPairs": r.GetCurrencyPairs,
		"eth_mintTokens":       r.MintTokens,
		"eth_burnTokens":       r.BurnTokens,
	}

	r.arrayResponseMethods = map[string]bool{
		"net_version": true,
	}
}

// HandleRPCRequest processes incoming JSON-RPC requests and returns appropriate responses
//
// The function:
// 1. Reads and parses the request body
// 2. Validates the request format
// 3. Routes to appropriate handler
// 4. Formats and returns the response
//
// Returns:
//   - domain.ServiceResponse: The formatted response
//   - error: Any error that occurred during processing
func (r *evmRouter) HandleRPCRequest(ctx zrouter.Context) (domain.ServiceResponse, error) {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	var request JSONRPCRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		var requests []JSONRPCRequest
		err = json.Unmarshal(body, &requests)
		if err != nil {
			return nil, fmt.Errorf("failed to decode request: %w", err)
		}
		if len(requests) == 0 {
			return nil, fmt.Errorf("empty request array")
		}
		request = requests[0]
	}

	handler, ok := r.methodHandlers[request.Method]
	if !ok {
		return nil, fmt.Errorf("unsupported method: %s", request.Method)
	}

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
	}

	headers := http.Header{}
	headers.Set("Content-Type", "text/plain; charset=utf-8")

	result, err := handler(request)
	if err != nil {
		logger.GetLoggerFromContext(ctx.Context()).Errorf("error with method %s and details %v", request.Method, err)
		response.Error = &JSONRPCError{
			Code:    1,
			Message: err.Error(),
		}
		return domain.NewServiceResponseWithHeader(http.StatusOK, response, headers), nil
	}

	response.Result = result
	if r.arrayResponseMethods[request.Method] {
		return domain.NewServiceResponseWithHeader(http.StatusOK, []JSONRPCResponse{response}, headers), nil
	}

	return domain.NewServiceResponseWithHeader(http.StatusOK, response, headers), nil
}
