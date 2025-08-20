package models

import (
	"errors"
	"fmt"
	"rest/goAPI/db"
)

type Registration struct {
	Id      int64
	UserId  int64
	EventId int64
}

func (r *Registration) Save() (int64, error) {
	query := `
	INSERT INTO registrations (user_id, event_id) VALUES (?, ?)
	`

	stmt, err := db.DB.Prepare(query)
	fmt.Println(err)
	if err != nil {
		return 0, errors.New("some unknown error occured")
	}

	defer stmt.Close()

	result, err := stmt.Exec(r.UserId, r.EventId)
	if err != nil {
		return 0, errors.New("registration failed, please try again later")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.New("some unknown error occured")
	}

	return id, nil
}

func GetRegistrationId(registrationId int64) (Registration, error) {
	query := `
	SELECT * FROM registrations WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return Registration{}, errors.New("some unknown error occured")
	}

	row := stmt.QueryRow(registrationId)
	var registration Registration
	err = row.Scan(&registration.Id, &registration.UserId, &registration.EventId)
	if err != nil {
		return Registration{}, errors.New("registration deletion failed, please try again later")
	}
	return registration, nil
}

func DeleteRegistration(registrationId, userId int64) error {
	registration, err := GetRegistrationId(registrationId)
	if err != nil {
		return errors.New("registration with id does not exist")
	}

	if registration.UserId != userId {
		return errors.New("user unauthorized")
	}

	query := `
	DELETE FROM registrations WHERE id = ?
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return errors.New("some unknown error occured")
	}

	_, err = stmt.Exec(registrationId)
	if err != nil {
		return errors.New("registration deletion failed, please try again later")
	}
	return err
}
