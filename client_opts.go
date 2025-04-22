package vaulty

import (
	"context"
	"errors"
	"log/slog"
	"os"

	hashiVault "github.com/hashicorp/vault/api"
	kubernetesAuth "github.com/hashicorp/vault/api/auth/kubernetes"
)

type ClientOption func(c *client) error

// WithContext sets the context for the client.
func WithContext(ctx context.Context) ClientOption {
	return func(c *client) error {
		c.ctx = ctx
		return nil
	}
}

// WithLogger sets the logger for the client.
func WithLogger(l *slog.Logger) ClientOption {
	return func(c *client) error {
		c.l = l
		return nil
	}
}

// WithGeneratedVaultClient creates a vault client with the given address.
//
// Deprecated: Use WithAddr instead for the same effect.
func WithGeneratedVaultClient(vaultAddress string) ClientOption {
	return WithAddr(vaultAddress)
}

func WithAddr(addr string) ClientOption {
	return func(c *client) error {
		c.config.Address = addr
		return nil
	}
}

func WithConfig(config *hashiVault.Config) ClientOption {
	return func(c *client) error {
		c.config = config
		return nil
	}
}

func WithTokenAuth(token string) ClientOption {
	return func(c *client) error {
		if token == "" {
			return errors.New("token is empty")
		}

		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			return tokenLogin(v, token)
		}
		return nil
	}
}

func WithAppRoleAuth(roleID, secretID string) ClientOption {
	return func(c *client) error {
		if roleID == "" {
			return errors.New("roleID is empty")
		} else if secretID == "" {
			return errors.New("secretID is empty")
		}

		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			sec, err := appRoleLogin(v, roleID, secretID)
			if err != nil {
				return nil, err
			}
			go c.renewAuthInfo()
			return sec, nil
		}
		return nil
	}
}

func WithUserPassAuth(username, password string) ClientOption {
	return func(c *client) error {
		if username == "" {
			return errors.New("username is empty")
		} else if password == "" {
			return errors.New("password is empty")
		}

		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			sec, err := userPassLogin(v, username, password)
			if err != nil {
				return nil, err
			}
			go c.renewAuthInfo()
			return sec, nil
		}
		return nil
	}
}

func WithKvv2Mount(mount string) ClientOption {
	return func(c *client) error {
		c.kvv2Mount = mount
		return nil
	}
}

func WithKubernetesAuthDefault(roleName string) ClientOption {
	return func(c *client) error {
		if roleName == "" {
			return errors.New("role name is empty")
		}

		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			sec, err := kubernetesLogin(v, roleName, kubernetesAuth.WithServiceAccountTokenPath(kubernetesServiceAccountTokenPath))
			if err != nil {
				return nil, err
			}

			go c.renewAuthInfo()

			return sec, nil
		}
		return nil
	}
}

func WithKubernetesAuthFromEnv() ClientOption {
	return func(c *client) error {
		roleFromEnv := os.Getenv(envServiceAccountName)
		if roleFromEnv == "" {
			return errors.New("role name is not set in environment variable")
		}

		return WithKubernetesAuthDefault(roleFromEnv)(c)
	}
}

func WithKubernetesAuth(role, token string) ClientOption {
	return func(c *client) error {
		if role == "" {
			return errors.New("role name is empty")
		} else if token == "" {
			return errors.New("token is empty")
		}

		c.auth = func(v *hashiVault.Client) (*hashiVault.Secret, error) {
			sec, err := kubernetesLogin(v, role, kubernetesAuth.WithServiceAccountToken(token))
			if err != nil {
				return nil, err
			}

			go c.renewAuthInfo()

			return sec, nil
		}
		return nil
	}
}
