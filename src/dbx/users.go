package dbx

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/unicode/norm"
)

type User struct {
	Email      string `db:"email"`
	Username   string `db:"username"`
	Password   string `db:"password"`
	AdminP     bool   `db:"adminp"`
	Translates string `db:"translates"`
}

var LoginFailed = errors.New("No user with the given username and password")

func CreateUser(user User) error {
	user.Username = norm.NFC.String(user.Username)
	user.Password = norm.NFC.String(user.Password)

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO users (
			email,
			username,
			password,
			adminp,
			translates
		) VALUES (?, ?, ?, ?, ?)
	`, user.Email, user.Username, string(hash), user.AdminP, user.Translates)
	return err
}

func Login(username, password string) (User, error) {
	username = norm.NFC.String(username)
	password = norm.NFC.String(password)

	/* TODO: Pass a context here? */
	rs, err := db.Queryx(`SELECT * FROM users WHERE username = ?`, username)
	if err != nil {
		return User{}, err
	}

	var u User
	switch err = rs.Scan(&u); {
	case errors.Is(err, sql.ErrNoRows):
		return User{}, LoginFailed
	case err != nil:
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return User{}, LoginFailed
	case err != nil:
		return User{}, err
	}
	return u, nil
}
