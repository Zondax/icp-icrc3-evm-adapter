package icp

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	icpDex "github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp/clients/dex"
	icpLogger "github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp/clients/logger"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/identity"
	"github.com/aviate-labs/agent-go/principal"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/conf"
)

// ICPClients holds both logger and dex clients
type Clients struct {
	Logger *icpLogger.Agent
	Dex    *icpDex.Agent
}

var (
	clients     *Clients
	clientsOnce sync.Once
)

// NewICPClient
// Parameters:
//   - cfg: Configuration containing canister IDs and node URL
//
// Returns:
//   - *ICPClients: A struct containing both Logger and DEX clients
//   - error: Any error that occurred during initialization
//
// The function will:
//  1. Parse and validate the canister IDs
//  2. Create a new identity for the agent
//  3. Initialize both Logger and DEX clients with the same configuration
//  4. Return a singleton instance of ICPClients
func NewICPClient(cfg *conf.ICPConfig) (*Clients, error) {
	var initErr error

	clientsOnce.Do(func() {
		loggerCanisterID, err := principal.Decode(cfg.LoggerCanisterID)
		if err != nil {
			initErr = fmt.Errorf("invalid Logger CanisterID '%s': %w", cfg.LoggerCanisterID, err)
			return
		}

		dexCanisterID, err := principal.Decode(cfg.DexCanisterID)
		if err != nil {
			initErr = fmt.Errorf("invalid DEX CanisterID '%s': %w", cfg.DexCanisterID, err)
			return
		}

		nodeURL, err := url.Parse(cfg.NodeURL)
		if err != nil {
			initErr = fmt.Errorf("invalid NodeURL '%s': %w", cfg.NodeURL, err)
			return
		}

		id, err := identity.NewRandomSecp256k1Identity()
		if err != nil {
			initErr = fmt.Errorf("failed to create identity: %w", err)
			return
		}

		timeOut, err := time.ParseDuration(cfg.Timeout)
		if err != nil {
			initErr = fmt.Errorf("failed to parse timeout: %w", err)
			return
		}

		agentConfig := agent.Config{
			ClientConfig: &agent.ClientConfig{
				Host: nodeURL,
			},
			FetchRootKey: true,
			Identity:     id,
			PollTimeout:  timeOut,
		}

		loggerAgent, err := icpLogger.NewAgent(loggerCanisterID, agentConfig)
		if err != nil {
			initErr = fmt.Errorf("failed to create logger agent: %w", err)
			return
		}

		dexAgent, err := icpDex.NewAgent(dexCanisterID, agentConfig)
		if err != nil {
			initErr = fmt.Errorf("failed to create dex agent: %w", err)
			return
		}

		clients = &Clients{
			Logger: loggerAgent,
			Dex:    dexAgent,
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return clients, nil
}
