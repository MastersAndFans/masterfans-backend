package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDBConnection(t *testing.T) {
	db, err := InitDB()
	require.NoError(t, err, "Failed to connect to database")

	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get generic database object from GORM")
	defer sqlDB.Close()

	err = sqlDB.Ping()
	require.NoError(t, err, "Failed to ping database")

	t.Log("Database connection is healthy")
}
