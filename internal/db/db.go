package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func InitDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// configure ssl mode in production and set sslmode=verify-full when deploying
	// do not hardcode the paths, just use environment variables
	// dsn := "host=host user=user password=password dbname=dbname port=port sslmode=verify-full sslrootcert=/path/to/root.crt sslcert=/path/to/postgresql.crt sslkey=/path/to/postgresql.key"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		FullSaveAssociations: true,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
