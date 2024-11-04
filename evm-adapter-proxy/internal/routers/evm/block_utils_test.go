package evm

import (
	"testing"

	"github.com/aviate-labs/agent-go/candid/idl"
	"github.com/stretchr/testify/assert"
	icpLogger "github.com/zondax/poc-icp-icrc3-evm-adapter/internal/icp/clients/logger"
)

func TestMapBlockToEVMBlock(t *testing.T) {
	tests := []struct {
		name        string
		blockValue  icpLogger.Value
		want        Block
		wantErr     bool
		errContains string
	}{
		{
			name: "Valid block",
			blockValue: icpLogger.Value{
				Map: &[]struct {
					Field0 string          `ic:"0" json:"0"`
					Field1 icpLogger.Value `ic:"1" json:"1"`
				}{
					{
						Field0: "id",
						Field1: icpLogger.Value{
							Nat: func() *idl.Nat {
								n := idl.NewNatFromString("1")
								return &n
							}(),
						},
					},
					{
						Field0: "hash",
						Field1: icpLogger.Value{
							Blob: &[]byte{1, 2, 3, 4},
						},
					},
					{
						Field0: "phash",
						Field1: icpLogger.Value{
							Blob: &[]byte{5, 6, 7, 8},
						},
					},
					{
						Field0: "ts",
						Field1: icpLogger.Value{
							Nat: func() *idl.Nat {
								n := idl.NewNatFromString("1000000000") // 1 second in nanoseconds
								return &n
							}(),
						},
					},
				},
			},
			want: Block{
				Number:           "0x1",
				Hash:             "0x01020304",
				ParentHash:       "0x05060708",
				Timestamp:        "0x1",
				Transactions:     []string{},
				TransactionsRoot: "0x0000000000000000000000000000000000000000000000000000000000000000",
				ReceiptsRoot:     "0x0000000000000000000000000000000000000000000000000000000000000000",
				StateRoot:        "0x0000000000000000000000000000000000000000000000000000000000000000",
			},
			wantErr: false,
		},
		{
			name: "Invalid block - missing Map",
			blockValue: icpLogger.Value{
				Map: nil,
			},
			wantErr:     true,
			errContains: "expected Map value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapBlockToEVMBlock(tt.blockValue)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.Number, got.Number)
			assert.Equal(t, tt.want.Hash, got.Hash)
			assert.Equal(t, tt.want.ParentHash, got.ParentHash)
			assert.Equal(t, tt.want.Timestamp, got.Timestamp)
			assert.Equal(t, tt.want.TransactionsRoot, got.TransactionsRoot)
			assert.Equal(t, tt.want.ReceiptsRoot, got.ReceiptsRoot)
			assert.Equal(t, tt.want.StateRoot, got.StateRoot)
		})
	}
}

func TestDecodeCertificateData(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		want        uint64
		wantErr     bool
		errContains string
	}{
		{
			name:    "Valid ULEB128 encoding",
			data:    []byte{0x2A}, // 42 in ULEB128
			want:    42,
			wantErr: false,
		},
		{
			name:    "Multi-byte ULEB128 encoding",
			data:    []byte{0xE5, 0x8E, 0x26}, // 624485 in ULEB128
			want:    624485,
			wantErr: false,
		},
		{
			name:        "Empty data",
			data:        []byte{},
			wantErr:     true,
			errContains: "failed to read byte",
		},
		{
			name:        "Incomplete ULEB128 encoding",
			data:        []byte{0x80}, // Incomplete encoding
			wantErr:     true,
			errContains: "failed to read byte",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeCertificateData(tt.data)
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
