package vaulty

import (
	"testing"

	hashiVault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		want  string
		input *hashiVault.Secret
	}{
		{
			name:  "nil secret",
			want:  "",
			input: nil,
		},
		{
			name:  "nil data",
			want:  "",
			input: &hashiVault.Secret{},
		},
		{
			name: "nil cipher text",
			want: "",
			input: &hashiVault.Secret{
				Data: make(map[string]any),
			},
		},
		{
			name: "invalid cipher text",
			want: "",
			input: &hashiVault.Secret{
				Data: map[string]any{
					TransitKeyCipherText: 1,
				},
			},
		},
		{
			name: "valid cipher text",
			want: "cipher text",
			input: &hashiVault.Secret{
				Data: map[string]any{
					TransitKeyCipherText: "cipher text",
				},
			},
		},
		{
			name: "valid cipher text: empty",
			want: "",
			input: &hashiVault.Secret{
				Data: map[string]any{
					TransitKeyCipherText: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := CipherTextFromSecret(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_UintToInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   uint
		out  int
		err  bool
	}{
		{
			name: "max int",
			in:   uint(maxInt),
			out:  maxInt,
			err:  false,
		},
		{
			name: "max int + 1",
			in:   uint(maxInt) + 1,
			out:  0,
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := uintToInt(tt.in)
			if tt.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.out, got)
		})
	}
}
