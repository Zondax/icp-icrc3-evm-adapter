package evm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWeb3ClientVersion(t *testing.T) {
	router := &evmRouter{}
	result, err := router.Web3ClientVersion(JSONRPCRequest{})

	assert.NoError(t, err)
	assert.Equal(t, "EVM-Adapter/v0.1.0", result)
}

func TestWeb3Sha3(t *testing.T) {
	tests := []struct {
		name        string
		input       []interface{}
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:  "Valid hex input",
			input: []interface{}{"0x68656c6c6f20776f726c64"}, // "hello world" in hex
			want:  "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad",
		},
		{
			name:        "Invalid hex input",
			input:       []interface{}{"invalid"},
			wantErr:     true,
			errContains: "invalid hex string",
		},
		{
			name:        "Empty params",
			input:       []interface{}{},
			wantErr:     true,
			errContains: "invalid params",
		},
	}

	router := &evmRouter{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := router.Web3Sha3(JSONRPCRequest{
				Params: tt.input,
			})

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}
