package evm

func (r *evmRouter) NetListening(_ JSONRPCRequest) (interface{}, error) {
	// Always return true as we're always "listening"
	return true, nil
}

func (r *evmRouter) NetPeerCount(_ JSONRPCRequest) (interface{}, error) {
	// Return 1 as we don't have peers in the traditional sense
	return "0x1", nil
}
