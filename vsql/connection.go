package vsql

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	hashiVault "github.com/hashicorp/vault/api"
	"github.com/jacobbrewer1/vaulty"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func ConnectDB(
	ctx context.Context,
	l *slog.Logger,
	client vaulty.Client,
	vip *viper.Viper,
	currentSecrets *hashiVault.Secret,
) error {
	if vip == nil {
		return errors.New("no viper configuration provided")
	}
	if !vip.IsSet("vault") {
		return errors.New("no vault configuration found")
	} else if vip.IsSet("vault.database") {
		return errors.New("no vault database configuration found")
	}
	if client == nil {
		return errors.New("no vault client provided")
	}
	if currentSecrets == nil {
		return errors.New("no vault database credentials provided")
	}

	connectionString := connectionStringFromViperAndVaultSecret(vip, currentSecrets)
	dbConn, err := openDBConnection(ctx, connectionString)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	db := NewDatabase(dbConn)

	l.Debug("Initial database connection established")

	go func() {
		if err := vaulty.RenewLease(
			ctx,
			l,
			client,
			"database_connection",
			currentSecrets,
			func() (*hashiVault.Secret, error) {
				l.Warn("Vault lease expired, establishing new database connection")

				newDatabaseCredentials, err := client.Path(vip.GetString("vault.database.role"),
					vaulty.WithPrefix(vip.GetString("vault.database.path")),
				).GetSecret(ctx)
				if err != nil {
					return nil, fmt.Errorf("failed to get new database credentials: %w", err)
				}

				newConnectionString := connectionStringFromViperAndVaultSecret(vip, newDatabaseCredentials)
				newDB, err := openDBConnection(ctx, newConnectionString)
				if err != nil {
					return nil, fmt.Errorf("failed to open new database connection: %w", err)
				}

				l.Debug("New database connection established, replacing old connection")

				if err := db.ReplaceDB(ctx, newDB); err != nil {
					return nil, fmt.Errorf("failed to replace database connection: %w", err)
				}

				l.Info("Database connection renewed successfully")
				return newDatabaseCredentials, nil
			},
		); err != nil {
			l.Error("Failed to renew database connection lease",
				slog.String(loggingKeyError, err.Error()),
			)
		}
	}()

	return nil
}

func openDBConnection(ctx context.Context, connectionString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := pingDB(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func connectionStringFromViperAndVaultSecret(vip *viper.Viper, secrets *hashiVault.Secret) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		secrets.Data["username"],
		secrets.Data["password"],
		vip.GetString("vault.database.host"),
		vip.GetInt("vault.database.port"),
		vip.GetString("vault.database.name"),
	)
	return connectionString
}
