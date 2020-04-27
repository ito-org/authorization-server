package main

import (
	"fmt"
	"strings"
	"math/big"
	"crypto/rand"
	
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewDBConnection creates and tests a new db connection and returns it.
func NewDBConnection(dbHost, dbUser, dbPassword, dbName string) (*DBConnection, error) {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		fmt.Printf("Failed to connect to Postgres database: %s\n", err.Error())
		return nil, err
	}
	return &DBConnection{db}, err
}

// DBConnection implements several functions for fetching and manipulation
// of reports in the database.
type DBConnection struct {
	*sqlx.DB
}

const tokenRunes = "abcdefghijklmnopqrstuvwxyz0123456789"
const tokenLength = 8 // keep in sync with db.sql

func (db *DBConnection) createToken() (string, error) {
	// TODO handle collisions (unique constraint violation)
	var tokenBuilder strings.Builder
	for i := 0; i < tokenLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(tokenRunes))))
		if err != nil {
			fmt.Printf("Failed to generate random number: %s\n", err.Error())
			return "", err
		}
		tokenBuilder.WriteRune(rune(tokenRunes[num.Int64()]))
	}
	
	token := tokenBuilder.String()
	
	_, err := db.Exec(
		`
		INSERT INTO
		tokens(token)
		VALUES($1);
		`,
		token,
	)
	if err != nil {
		fmt.Printf("Failed to insert token into database: %s\n", err.Error())
		return "", err
	}
	return token, nil
}

func (db *DBConnection) checkAndRemoveToken(token string) (bool, error) {
	result, err := db.Exec(
		`
		DELETE FROM
		tokens
		WHERE
		token = $1;
		`,
		token,
	)
	if err != nil {
		fmt.Printf("Failed to delete token from database: %s\n", err.Error())
		return false, err
	}
	
	rows, err := result.RowsAffected()
	
	if err != nil {
		fmt.Printf("Failed to get rows affected: %s\n", err.Error())
		return false, err
	}
	return rows == 1, nil
}
