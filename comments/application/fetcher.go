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

// Fetch total reactions for comments
func (f *CommentFetcher) FetchReactionStats(ctx context.Context, entities []domain.Entity) (map[int]domain.ReactionStat, error) {
	url := fmt.Sprintf("http://reactions:8080/stats")
	body := strings.NewReader(fmt.Sprintf(`{"data": %v}`, entities))

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

	var reactionStats []domain.ReactionStat
	if err := json.NewDecoder(resp.Body).Decode(&reactionStats); err != nil {
		return nil, err
	}
	return mapReactionsStats(reactionStats), nil
}

func mapReactionsStats(reactionStats []domain.ReactionStat) map[int]domain.ReactionStat {
	reactionStatsMap := make(map[int]domain.ReactionStat)
	for _, tr := range reactionStats {
		reactionStatsMap[tr.EntityId] = tr
	}
	return reactionStatsMap
}
