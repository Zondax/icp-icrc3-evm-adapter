package evm

import (
	"github.com/zondax/golem/pkg/zrouter"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp"
)

func NewEVMRouter(zr zrouter.ZRouter, icpClient *icp.Agent) {
	r := &evmRouter{
		icpClient: icpClient,
	}
	r.initMethodHandlers()

	zr.POST("/rpc/v1", r.HandleRPCRequest)
}
