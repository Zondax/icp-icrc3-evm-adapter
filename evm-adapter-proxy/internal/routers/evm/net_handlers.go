package evm

// NetListening implements the net_listening RPC method
// Returns whether the client is actively listening for network connections
//
// Returns:
//   - bool: Always returns true for this PoC implementation
//   - error: Always returns nil for this implementation
func (r *evmRouter) NetListening(_ JSONRPCRequest) (interface{}, error) {
	// Always return true as we're always "listening"
	return true, nil
}

// NetPeerCount implements the net_peerCount RPC method
// Returns the number of peers currently connected to the client
//
// Returns:
//   - string: Always returns "0x1" for this PoC implementation
//   - error: Always returns nil for this implementation
//
// Note: This is a placeholder implementation that always returns 1 peer
func (r *evmRouter) NetPeerCount(_ JSONRPCRequest) (interface{}, error) {
	// Return 1 as we don't have peers in the traditional sense
	return "0x1", nil
}
