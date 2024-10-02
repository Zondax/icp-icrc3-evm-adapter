package service

import (
	"github.com/zondax/golem/pkg/logger"
	gm "github.com/zondax/golem/pkg/metrics"
	"github.com/zondax/golem/pkg/runner"
	"github.com/zondax/golem/pkg/zrouter"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/conf"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/routers/evm"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/version"
	"go.uber.org/zap"
)

const (
	appName = "evm-adapter-proxy"
)

func Start(c *conf.Config) {
	logger.InitLogger(logger.Config{Level: c.Logging.Level})
	metricServer := gm.NewTaskMetrics(c.Metrics.Path, c.Metrics.Port, appName)

	zrouterConfig := &c.RouterConfig
	zrouterConfig.AppRevision = version.GitRevision
	zrouterConfig.AppVersion = version.GitVersion

	zr := zrouter.New(metricServer, zrouterConfig)

	tr := runner.NewRunner()
	tr.AddTask(metricServer)
	tr.Start()

	icpClient, err := icp.NewICPClient(c.ICP)
	if err != nil {
		zap.S().Fatalf("Error initializing icpClient: %x", err)
	}
	evm.NewEVMRouter(zr, icpClient)

	zap.S().Fatal(zr.Run(c.ServerPort))
}
