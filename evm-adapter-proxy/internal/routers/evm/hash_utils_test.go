package evm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	icpLogger "github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp/clients/logger"
)

func TestGetBlockHash(t *testing.T) {
	tests := []struct {
		name        string
		block       icpLogger.Value
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name: "Valid block with hash",
			block: icpLogger.Value{
				Map: &[]struct {
					Field0 string          `ic:"0" json:"0"`
					Field1 icpLogger.Value `ic:"1" json:"1"`
				}{
					{
						Field0: "hash",
						Field1: icpLogger.Value{
							Blob: &[]byte{1, 2, 3, 4},
						},
					},
				},
			},
			want:    "0x01020304",
			wantErr: false,
		},
		{
			name: "Block without Map",
			block: icpLogger.Value{
				Map: nil,
			},
			wantErr:     true,
			errContains: "invalid block format: Map is nil",
		},
		{
			name: "Block without hash field",
			block: icpLogger.Value{
				Map: &[]struct {
					Field0 string          `ic:"0" json:"0"`
					Field1 icpLogger.Value `ic:"1" json:"1"`
				}{
					{
						Field0: "other",
						Field1: icpLogger.Value{},
					},
				},
			},
			wantErr:     true,
			errContains: "block hash not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBlockHash(tt.block)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
