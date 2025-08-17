package models

import (
	"errors"
	"fmt"
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
	UserId      int64
}

type UpdateEvent struct {
	Name        string
	Description string
	Location    string
	DateTime    time.Time
	UserId      int64
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

	defer rows.Close()
	var events []Event

	for rows.Next() {
		var event Event

		err := rows.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserId)
		if err != nil {
			return nil, errors.New("fetching a row from the DB failed")
		}

		events = append(events, event)
	}

	fmt.Println(events)

	return events, nil
}

func GetEventByID(id string) (Event, error) {
	query := `
	SELECT * FROM events WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return Event{}, errors.New("could not prepare query")
	}

	row := stmt.QueryRow(id)

	stmt.Close()

	var event Event
	err = row.Scan(&event.Id, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserId)
	if err != nil {
		return Event{}, errors.New("some error in scanning the event row")
	}

	return event, nil
}

func UpdateEventById(id string, details UpdateEvent) error {
	oldEvent, err := GetEventByID(id)

	if err != nil {
		return errors.New("could not update event, could not fetch old event")
	}

	query := `
	UPDATE events
	SET name=?, description=?, location=?, dateTime=?, user_id=?
	WHERE id=?
	`

	if len(details.Name) > 0 {
		oldEvent.Name = details.Name
	}
	if len(details.Description) > 0 {
		oldEvent.Description = details.Description
	}
	if len(details.Location) > 0 {
		oldEvent.Location = details.Location
	}
	if !details.DateTime.IsZero() {
		oldEvent.DateTime = details.DateTime
	}
	if details.UserId != 0 {
		oldEvent.UserId = details.UserId
	}

	_, err = db.DB.Exec(query, oldEvent.Name, oldEvent.Description, oldEvent.Location, oldEvent.DateTime, oldEvent.UserId, oldEvent.Id)

	if err != nil {
		return errors.New("could not update event")
	}

	return nil
}

func DeleteEventById(id string) error {
	query := `
	DELETE FROM events WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return errors.New("some error occured during query prep")
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return errors.New("could not delete event")
	}

	return nil
}
