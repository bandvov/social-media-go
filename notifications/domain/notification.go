package domain

import "fmt"

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

// GenerateMessage generates a notification message based on the type
func (n Notification) GenerateMessage() string {
	switch n.Type {
	case NewFollower:
		return fmt.Sprintf("You have a new follower!")

	case NewMention:
		return fmt.Sprintf("You were mentioned in a %s.", n.EntityType)

	case NewDirectMessage:
		return fmt.Sprintf("You have a new direct message.")

	case NewPostComment:
		return fmt.Sprintf("Someone commented on your post.")

	case NewCommentReply:
		return fmt.Sprintf("Someone replied to your comment.")

	default:
		// Handle reactions separately
		if isReaction(n.Type) {
			return generateReactionMessage(n.ActorIDs, n.Type, n.EntityType)
		}
	}

	return "You have a new notification."
}

// isReaction checks if the notification type is a reaction
func isReaction(notificationType NotificationType) bool {
	reactionTypes := map[NotificationType]bool{
		NewReactionLike:    true,
		NewReactionDislike: true,
		NewReactionLove:    true,
		NewReactionLaugh:   true,
		NewReactionAngry:   true,
		NewReactionWow:     true,
	}
	return reactionTypes[notificationType]
}

// getReactionVerb returns the correct verb based on the reaction type
func getReactionVerb(notificationType NotificationType) string {
	switch notificationType {
	case NewReactionLike:
		return "liked"
	case NewReactionDislike:
		return "disliked"
	case NewReactionLove:
		return "loved"
	case NewReactionLaugh:
		return "laughed at"
	case NewReactionAngry:
		return "reacted angrily to"
	case NewReactionWow:
		return "were amazed by"
	default:
		return "reacted to"
	}
}

// generateReactionMessage creates a message for reactions with formatted names
func generateReactionMessage(ActorIDs []int, reactionType NotificationType, entityType EntityType) string {
	count := len(ActorIDs)
	verb := getReactionVerb(reactionType)

	if count == 0 {
		return fmt.Sprintf("Someone %s your %s.", verb, entityType)
	}

	switch count {
	case 1:
		return fmt.Sprintf("%v %v your %v.", ActorIDs[0], verb, entityType)
	case 2:
		return fmt.Sprintf("%v and %v %v your %s.", ActorIDs[0], ActorIDs[1], verb, entityType)
	default:
		return fmt.Sprintf("%v, %v and %v more %s your %s.", ActorIDs[0], ActorIDs[1], count-2, verb, entityType)
	}
}
