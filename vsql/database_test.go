package vsql

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func mockSqlxDB(t *testing.T, options ...func(mockedSQL sqlmock.Sqlmock)) *sqlx.DB {
	t.Helper()

	dbMock, sqlMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)

	for _, option := range options {
		option(sqlMock)
	}

	sqlxDB := &sqlx.DB{
		DB: dbMock,
	}

	return sqlxDB
}

func TestReplaceDB(t *testing.T) {
	ctx := context.Background()
	newDbPingTimeout = 10 * time.Millisecond

	t.Run("returns error when newDB is nil", func(t *testing.T) {
		db := NewDatabase(mockSqlxDB(t))
		err := db.ReplaceDB(ctx, nil)
		require.EqualError(t, err, "new database connection is nil")
	})

	t.Run("returns error when pinging newDB fails", func(t *testing.T) {
		db := NewDatabase(mockSqlxDB(t))
		err := db.ReplaceDB(ctx, mockSqlxDB(t, func(mockedSQL sqlmock.Sqlmock) {
			mockedSQL.ExpectPing().WillReturnError(errors.New("ping failed"))
		}))
		require.Error(t, err, "failed to ping new database: ping failed")
	})

	t.Run("successfully replaces database connection", func(t *testing.T) {
		db := NewDatabase(mockSqlxDB(t, func(mockedSQL sqlmock.Sqlmock) {
			mockedSQL.ExpectClose().WillReturnError(nil)
		}))
		err := db.ReplaceDB(ctx, mockSqlxDB(t, func(mockedSQL sqlmock.Sqlmock) {
			mockedSQL.ExpectPing().WillReturnError(nil)
		}))
		require.NoError(t, err)
	})
}

func TestPingDb(t *testing.T) {
	newDbPingTimeout = 100 * time.Millisecond

	t.Run("returns error when pinging times out", func(t *testing.T) {
		db := mockSqlxDB(t, func(mockedSQL sqlmock.Sqlmock) {
			mockedSQL.ExpectPing().WillDelayFor(20 * time.Millisecond).WillReturnError(errors.New("ping failed"))
		})

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		err := pingDB(ctx, db)
		require.EqualError(t, err, "failed to successfully ping database before context timeout: ping failed")
	})

	t.Run("successfully pings database", func(t *testing.T) {
		db := mockSqlxDB(t, func(mockedSQL sqlmock.Sqlmock) {
			mockedSQL.ExpectPing().WillReturnError(nil)
		})

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		err := pingDB(ctx, db)
		require.NoError(t, err)
	})
}
