package datastore

import (
	"auth-users-service/pkg/models"
	"database/sql"
	"fmt"

	// sweet psql driver

	_ "github.com/lib/pq"
)

// DB is the sharable pool of connections for the app
var DB *sql.DB

// InitializeDB starts a connection to the database and exposes it to the application
func InitializePostgres(connectionString string) error {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return err
	}

	DB = db

	fmt.Println("Successful Postgres connection")
	return nil
}

func FindUserByUsername(username string) (models.User, error) {
	var user models.User

	// get user from db
	err := DB.QueryRow("SELECT uuid, email_address, password FROM users WHERE username = $1", username).Scan(
		&user.UserUUID,
		&user.EmailAddress,
		&user.Password,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}
