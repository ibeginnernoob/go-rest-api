package models

import (
	"errors"
	"rest/goAPI/db"
	"rest/goAPI/utils"
)

type User struct {
	Id       int64
	Name     string `binding:"required"`
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func GetUsers() ([]User, error) {
	query := `SELECT * FROM users`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, errors.New("could not fetch users")
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User

		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, errors.New("fetching a row from the DB failed")
		}

		users = append(users, user)
	}

	return users, nil
}

func (u User) SaveUser() (int64, error) {
	query := `
	INSERT INTO users (name, email, password)
	VALUES (?, ?, ?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return -1, errors.New("could not prepare query")
	}

	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return -1, errors.New("could not generated hashed password")
	}

	result, err := stmt.Exec(u.Name, u.Email, hashedPassword)
	if err != nil {
		return -1, errors.New("could not save user")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, errors.New("could not retrieve userId")
	}

	return id, nil
}

func GetUserByEmail(email string) (User, error) {
	query := `
	SELECT * FROM users WHERE email = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return User{}, errors.New("could not prepare query")
	}

	defer stmt.Close()

	row := stmt.QueryRow(email)

	var user User
	err = row.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return User{}, errors.New("could not scan user details")
	}

	return user, nil
}

func DeleteUsers() error {
	query := `
	DELETE FROM users
	`
	_, err := db.DB.Exec(query)
	if err != nil {
		return errors.New("could not delete all users")
	}

	return nil
}
