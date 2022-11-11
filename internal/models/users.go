package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	conn, err := m.DB.Acquire(context.Background())
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), "insert into users(name, email, hashed_password, created) values ($1, $2, $3, now())", name, email, string(hashedPassword))
	if err != nil {
		if strings.Contains(err.Error(), "users_uc_email") {
			return ErrDuplicateEmail
		}
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	conn, err := m.DB.Acquire(context.Background())
	if err != nil {
		return 0, err
	}

	row := conn.QueryRow(context.Background(), "SELECT id, hashed_password from users where email = $1", email)
	err = row.Scan(&id, &hashedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
