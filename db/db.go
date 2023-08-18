package db

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

type DB struct {
	DB *sql.DB
}

func New() DB {
	username := viper.GetString("database.user")
	password := viper.GetString("database.password")
	databaseName := viper.GetString("database.name")
	databaseHost := viper.GetString("database.host")
	databasePort := viper.GetString("database.port")
	sslmode := viper.GetString("database.sslmode")

	dbDSN := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s port=%s", databaseHost, username, databaseName, sslmode, password, databasePort)

	config, err := pgx.ParseConfig(dbDSN)
	if err != nil {
		panic(err)
	}

	db := stdlib.OpenDB(*config)

	// db, err := sql.Open("postgres", dbDSN)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	return DB{
		DB: db,
	}
}
