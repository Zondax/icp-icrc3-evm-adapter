package evm

import (
	"github.com/zondax/golem/pkg/zrouter"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp"
)

// NewEVMRouter creates and initializes a new EVM-compatible router
//
// Parameters:
//   - zr: The base router to add EVM routes to
//   - icpClients: The ICP clients (Logger and DEX) to use for operations
//
// The router will:
//  1. Initialize an evmRouter instance with the provided clients
//  2. Set up all supported RPC method handlers
//  3. Add the main RPC endpoint (/rpc/v1)
func NewEVMRouter(zr zrouter.ZRouter, icpClients *icp.Clients) {
	r := &evmRouter{
		icpClients: icpClients,
	}
	r.initMethodHandlers()

	zr.POST("/rpc/v1", r.HandleRPCRequest)
}
