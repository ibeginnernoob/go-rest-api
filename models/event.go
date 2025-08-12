package models

import (
	"errors"
	"rest/goAPI/db"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Event struct {
	Id          int64
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Location    string `binding:"required"`
	DateTime    time.Time
	UserId      int
}

func (e *Event) Save() error {
	insertQuery := `
	INSERT INTO events (name, description, location, dateTime, user_id)
	VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := db.DB.Prepare(insertQuery)
	if err != nil {
		return errors.New("could not run the insert query")
	}

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserId)
	if err != nil {
		return errors.New("could not insert the new event")
	}

	defer stmt.Close()

	id, err := result.LastInsertId()
	if err != nil {
		return errors.New("some error occured")
	}

	e.Id = id
	return nil
}

func GetEvents() ([]Event, error) {
	getQuery := `
	SELECT * FROM events
	`

	rows, err := db.DB.Query(getQuery)
	if err != nil {
		return nil, errors.New("could not fetch events")
	}

	rows.Close()
	var events []Event

	for rows.Next() {
		var event Event

		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserId)
		if err != nil {
			return nil, errors.New("fetching a row from the DB failed")
		}

		events = append(events, event)
	}

	return events, nil
}

func GetEvent(id string) (Event, error) {
	query := `
	SELECT * FROM events WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return Event{}, errors.New("could not prepare query")
	}

	rows, err := stmt.Query(id)
	if err != nil {
		return Event{}, errors.New("could not fetch event")
	}

	stmt.Close()

	var event Event
	for rows.Next() {
		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserId)
		if err != nil {
			return Event{}, errors.New("some error in scanning the event row")
		}
	}

	return event, nil
}
