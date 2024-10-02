package icp

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/aviate-labs/agent-go"
	"github.com/aviate-labs/agent-go/identity"
	"github.com/aviate-labs/agent-go/principal"
	"github.com/zondax/poc-icp-icrc3-evm-adapter/internal/conf"
)

var (
	client     *Agent
	clientOnce sync.Once
)

// NewICPClient initializes and returns a new ICP client using the provided configuration.
func NewICPClient(cfg *conf.ICPConfig) (*Agent, error) {
	var initErr error

	clientOnce.Do(func() {
		canisterID, err := principal.Decode(cfg.CanisterID)
		if err != nil {
			initErr = fmt.Errorf("invalid CanisterID '%s': %w", cfg.CanisterID, err)
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

		agentConfig := agent.Config{
			ClientConfig: &agent.ClientConfig{Host: nodeURL},
			FetchRootKey: true,
			Identity:     id,
		}

		agentInstance, err := NewAgent(canisterID, agentConfig)
		if err != nil {
			initErr = fmt.Errorf("failed to create agent: %w", err)
			return
		}

		client = agentInstance
	})

	// Return any error that occurred during initialization
	if initErr != nil {
		return nil, initErr
	}

	return client, nil
}
