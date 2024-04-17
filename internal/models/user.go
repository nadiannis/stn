package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         string
	Email      string
	Password   []byte
	DateJoined time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	id := uuid.New()

	stmt := `INSERT INTO users (id, email, password, date_joined)
					 VALUES (?, ?, ?, UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, id, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (uuid.UUID, error) {
	var id uuid.UUID
	var hashedPassword []byte

	stmt := "SELECT id, password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, ErrInvalidCredentials
		} else {
			return uuid.Nil, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return uuid.Nil, ErrInvalidCredentials
		} else {
			return uuid.Nil, err
		}
	}

	return id, nil
}
