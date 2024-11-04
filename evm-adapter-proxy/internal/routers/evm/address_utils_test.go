package evm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertEthAddressToICPPrincipal(t *testing.T) {
	tests := []struct {
		name        string
		ethAddr     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "Valid ICP principal with 0x prefix",
			ethAddr: "0x2vxsx-fae",
			wantErr: false,
		},
		{
			name:    "Valid ICP principal without 0x prefix",
			ethAddr: "2vxsx-fae",
			wantErr: false,
		},
		{
			name:        "Invalid ICP principal",
			ethAddr:     "0xinvalid-principal",
			wantErr:     true,
			errContains: "failed to decode principal",
		},
		{
			name:        "Empty address",
			ethAddr:     "",
			wantErr:     true,
			errContains: "failed to decode principal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			principal, err := ConvertEthAddressToICPPrincipal(tt.ethAddr)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, principal)
			}
		})
	}
}
