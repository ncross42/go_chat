package db

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser 새 사용자 만듦
func CreateUser(username, password string) (id int, err error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	err = db.QueryRow(`INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id`, username, string(passwordHash)).Scan(&id)
	return
}

// ErrUnauthorized 권한 없음 에러
var ErrUnauthorized = errors.New("db: unauthorized")

// FindUser 닉네임, 비번 가지고 사용자 찾기
func FindUser(username, password string) (id int, err error) {
	var passwordHash string
	err = db.QueryRow(`SELECT id, password_hash FROM users
		WHERE username = $1`, username).
		Scan(&id, &passwordHash)

	if err == sql.ErrNoRows {
		return 0, ErrUnauthorized
	} else if err != nil {
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash), []byte(password)); err != nil {
		return 0, ErrUnauthorized
	}

	return
}
