package evm

import (
	"fmt"
	"strings"

	"github.com/aviate-labs/agent-go/principal"
)

// ConvertEthAddressToICPPrincipal converts an Ethereum address to an ICP principal
// For PoC purposes, we expect the Ethereum address to be in format "0x<full_icp_principal>"
// Example: if ICP principal is "2vxsx-fae", the Ethereum address should be "0x2vxsx-fae"
//
// Parameters:
//   - ethAddr: Ethereum address (with 0x prefix)
//
// Returns:
//   - principal.Principal: The converted ICP principal
//   - error: Any error that occurred during conversion
func ConvertEthAddressToICPPrincipal(ethAddr string) (principal.Principal, error) {
	principalStr := strings.TrimPrefix(ethAddr, "0x")

	principalID, err := principal.Decode(principalStr)
	if err != nil {
		return principal.Principal{}, fmt.Errorf("failed to decode principal: %w", err)
	}

	return principalID, nil
}
