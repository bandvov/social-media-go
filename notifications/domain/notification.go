package domain

type NotificationType string

const (
	NewFollower      NotificationType = "new_follower"
	NewReaction      NotificationType = "new_reaction"
	NewMention       NotificationType = "new_mention"
	NewDirectMessage NotificationType = "new_direct_message"
	NewPostComment   NotificationType = "new_post_comment"
	NewCommentReply  NotificationType = "new_comment_reply"
)

type EntityType string

const (
	Post    EntityType = "post"
	Comment EntityType = "comment"
	Reply   EntityType = "reply"
)

type Notification struct {
	ID         int              `json:"id"`
	UserID     int              `json:"user_id"`     // Who receives the notification
	Type       NotificationType `json:"type"`        // Type of notification
	Message    string           `json:"message"`     // Notification message
	EntityType EntityType       `json:"entity_type"` // Type of entity (post, comment, reply)
	EntityID   int              `json:"entity_id"`   // ID of the post, comment, or reply
	CreatedAt  string           `json:"created_at"`
}
