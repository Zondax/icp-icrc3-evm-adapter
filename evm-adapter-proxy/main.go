package main

import (
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/commands"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/conf"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/version"
	"strings"
)

import (
	"github.com/zondax/golem/pkg/cli"
)

func main() {
	appName := "poc-icp-icrc3-evm-adapter"
	envPrefix := strings.ReplaceAll(appName, "-", "_")

	appSettings := cli.AppSettings{
		Name:        appName,
		Description: "Please override",
		ConfigPath:  "$HOME/.poc-icp-icrc3-evm-adapter/",
		EnvPrefix:   envPrefix,
		GitVersion:  version.GitVersion,
		GitRevision: version.GitRevision,
	}

	// Define application level features
	cli := cli.New[conf.Config](appSettings)
	defer cli.Close()

	cli.GetRoot().AddCommand(commands.GetStartCommand(cli))

	cli.Run()
}
