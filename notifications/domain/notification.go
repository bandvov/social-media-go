package domain

import "time"

type Notification struct {
	ID        int
	UserID    string
	Message   string
	CreatedAt time.Time
}
