package models

type Registration struct {
	id      int64
	userId  int64
	eventId int64
}

// func (r Registration) Save() {
// 	query := `
// 	INSERT INTO registrations (user_id, event_id) VALUES (?, ?)
// 	`
// }
