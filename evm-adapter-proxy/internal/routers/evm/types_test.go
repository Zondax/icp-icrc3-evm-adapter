package evm

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestConvertHexAmountToBigInt(t *testing.T) {
	tests := []struct {
		name        string
		amount      hexutil.Big
		want        *big.Int
		wantErr     bool
		errContains string
	}{
		{
			name:   "Valid hex amount",
			amount: hexutil.Big(*new(big.Int).SetInt64(100)),
			want:   big.NewInt(100),
		},
		{
			name:   "Zero amount",
			amount: hexutil.Big(*new(big.Int).SetInt64(0)),
			want:   big.NewInt(0),
		},
		{
			name: "Large amount",
			amount: func() hexutil.Big {
				n, _ := new(big.Int).SetString("deadbeef", 16)
				return hexutil.Big(*n)
			}(),
			want: func() *big.Int {
				n, _ := new(big.Int).SetString("deadbeef", 16)
				return n
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertHexAmountToBigInt(tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestHexToDecimal(t *testing.T) {
	tests := []struct {
		name        string
		hex         string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name: "Valid hex with 0x prefix",
			hex:  "0x10",
			want: "16",
		},
		{
			name: "Valid hex without prefix",
			hex:  "10",
			want: "16",
		},
		{
			name: "Zero",
			hex:  "0x0",
			want: "0",
		},
		{
			name:        "Invalid hex",
			hex:         "0xZZ",
			wantErr:     true,
			errContains: "invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hexToDecimal(tt.hex)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
