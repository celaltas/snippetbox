package postgresql

import (
	"database/sql"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"snippetbox/pkg/models"
	"strings"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES($1, $2, $3, CURRENT_TIMESTAMP)`

	_, err = u.DB.Exec(stmt, name, email, hashedPassword)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if strings.Contains(pqError.Message, "unique") {
				return models.ErrDuplicateEmail
			}
		}

	}
	return err
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	row := u.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
