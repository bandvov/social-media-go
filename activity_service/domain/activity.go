package domain

import (
	"encoding/json"
	"time"
)

// Activity represents an event in the activity feed.
type Activity struct {
	ID        int             `json:"id"`
	UserID    int             `json:"user_id"`
	Action    string          `json:"action"`               // e.g., "liked", "commented"
	TargetID  int             `json:"target_id"`            // ID of the entity being interacted with
	EventData json.RawMessage `json:"event_data,omitempty"` // JSON field for event details
	CreatedAt time.Time       `json:"created_at"`
}
