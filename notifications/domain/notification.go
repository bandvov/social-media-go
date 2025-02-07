package domain

type NotificationType string

const (
	NewFollower        NotificationType = "new_follower"
	NewReactionLike    NotificationType = "new_reaction_like"
	NewReactionDislike NotificationType = "new_reaction_dislike"
	NewReactionLove    NotificationType = "new_reaction_love"
	NewReactionLaugh   NotificationType = "new_reaction_laugh"
	NewReactionAngry   NotificationType = "new_reaction_angry"
	NewReactionWow     NotificationType = "new_reaction_wow"
	NewMention         NotificationType = "new_mention"
	NewDirectMessage   NotificationType = "new_direct_message"
	NewPostComment     NotificationType = "new_post_comment"
	NewCommentReply    NotificationType = "new_comment_reply"
)

type EntityType string

const (
	Post    EntityType = "post"
	Comment EntityType = "comment"
	Reply   EntityType = "reply"
)

type NotificationRequest struct {
	BaseNotification
	SenderId int `json:"sender_id"`
}

type BaseNotification struct {
	ID         int              `json:"id"`
	UserID     int              `json:"user_id"`     // Who receives the notification
	Type       NotificationType `json:"type"`        // Type of notification
	EntityType EntityType       `json:"entity_type"` // Type of entity (post, comment, reply)
	EntityID   int              `json:"entity_id"`   // ID of the post, comment, or reply
}

type Notification struct {
	BaseNotification
	ActorIDs  []int  `json:"actor_ids"` // array of users who reacted during half hour
	CreatedAt string `json:"created_at"`
}
