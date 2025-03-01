package application

import (
	"comments/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Helper struct for fetching data
type CommentFetcher struct {
	httpClient *http.Client
}

func NewCommentFetcher(client *http.Client) *CommentFetcher {
	return &CommentFetcher{
		httpClient: client,
	}
}

// Fetch user details
func (f *CommentFetcher) FetchUsersByID(ctx context.Context, userIDs []int) (map[int]domain.User, error) {
	url := fmt.Sprintf("http://users:8080/by-ids")
	body := strings.NewReader(fmt.Sprintf(`{"user_ids": %v}`, userIDs))

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []domain.User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	userMap := make(map[int]domain.User)
	for _, user := range users {
		userMap[user.ID] = user
	}
	return userMap, nil
}

// // Fetch user reactions for comments
// func (f *CommentFetcher) FetchReactions(ctx context.Context, commentIDs []int, userID int) (map[int][]domain.Reaction, error) {
// 	url := fmt.Sprintf("%s/reactions", f.reactionURL)
// 	body := strings.NewReader(fmt.Sprintf(`{"comment_ids": %v, "user_id": %d}`, commentIDs, userID))

// 	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := f.httpClient.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var reactions []domain.Reaction
// 	if err := json.NewDecoder(resp.Body).Decode(&reactions); err != nil {
// 		return nil, err
// 	}

// 	reactionMap := make(map[int]domain.Reaction)
// 	for _, reaction := range reactions {
// 		reactionMap[reaction.EntityID] = reaction
// 	}
// 	return reactionMap, nil
// }

// Fetch total reactions for comments
func (f *CommentFetcher) FetchTotalReactions(ctx context.Context, commentIDs []int) (map[int]int, error) {
	url := fmt.Sprintf("http://reactions:8080/total-reactions")
	body := strings.NewReader(fmt.Sprintf(`{"comment_ids": %v}`, commentIDs))

	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var totalReactions []domain.CommentCount
	if err := json.NewDecoder(resp.Body).Decode(&totalReactions); err != nil {
		return nil, err
	}

	totalReactionMap := make(map[int]int)
	for _, tr := range totalReactions {
		totalReactionMap[tr.EntityID] = tr.CommentCount + tr.ReplyCount
	}
	return totalReactionMap, nil
}
