package domain

import "time"

type Notification struct {
	ID        string
	UserID    string
	Message   string
	CreatedAt time.Time
}
