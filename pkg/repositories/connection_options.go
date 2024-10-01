package repositories

import (
	"context"

	hashiVault "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

type ConnectionOption func(c *databaseConnector)

func WithVaultClient(client vault.Client) ConnectionOption {
	return func(c *databaseConnector) {
		c.client = client
	}
}

func WithViper(v *viper.Viper) ConnectionOption {
	return func(c *databaseConnector) {
		c.vip = v
	}
}

func WithCurrentSecrets(secrets *hashiVault.Secret) ConnectionOption {
	return func(c *databaseConnector) {
		c.currentSecrets = secrets
	}
}

func WithContext(ctx context.Context) ConnectionOption {
	return func(c *databaseConnector) {
		c.ctx = ctx
	}
}
