package db

import (
	"fmt"
	"testing"
)

func TestDBConnection(t *testing.T) {
	db, err := InitDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get generic database object from GORM: %v", err)
	}
	defer sqlDB.Close()

	err = sqlDB.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Database connection is healthy")
}
