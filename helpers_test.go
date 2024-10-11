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
				Data: map[string]interface{}{},
			},
		},
		{
			name: "invalid cipher text",
			want: "",
			input: &hashiVault.Secret{
				Data: map[string]interface{}{
					TransitKeyCipherText: 1,
				},
			},
		},
		{
			name: "valid cipher text",
			want: "cipher text",
			input: &hashiVault.Secret{
				Data: map[string]interface{}{
					TransitKeyCipherText: "cipher text",
				},
			},
		},
		{
			name: "valid cipher text: empty",
			want: "",
			input: &hashiVault.Secret{
				Data: map[string]interface{}{
					TransitKeyCipherText: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CipherTextFromSecret(tt.input)
			require.Equal(t, tt.want, got)
		})
	}
}
