package postgres

import (
	"context"
	"cs01.com/snippetbox/pkg/models"
	"database/sql"
	"errors"

	//"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
	//"strings"
	"time"
)

type UserModel struct {
	Pool *pgxpool.Pool
}

const (
	insertSQL = "INSERT INTO users (name, email, hashed_password, created) VALUES($1, $2, $3, $4) "
	authSQL   = "SELECT id, hashed_password FROM users WHERE email = $1 AND active = TRUE"
)

func (m *UserModel) Insert(name, email, password string) error {
	var id uint64
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	row := m.Pool.QueryRow(context.Background(), insertSQL, name, email, string(hashedPassword), time.Now())
	err = row.Scan(&id)
	if err != nil {
		return nil
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := m.Pool.QueryRow(context.Background(), authSQL, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Otherwise, the password is correct. Return the user ID.
	return id, nil

}

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
