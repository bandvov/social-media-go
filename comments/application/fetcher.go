package application

import (
	"bytes"
	"comments/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type UserRequest struct {
	Data []int `json:"data"`
}

type StatsRequest struct {
	Data []domain.Entity `json:"data"`
}

type userResponse struct {
	Data []domain.User `json:"data"`
}

type ReactionsStatsResponse struct {
	Data []domain.ReactionStat `json:"data"`
}

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

	requestBody, err := json.Marshal(UserRequest{Data: userIDs})
	if err != nil {
		return nil, err
	}
	body := strings.NewReader(string(requestBody))
	fmt.Fprintln(os.Stdout, "user body", string(requestBody))
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

	var res userResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, "user response", res)
	return sliceToMap(res.Data, func(u domain.User) int {
		return u.ID
	}), nil
}

// Fetch total reactions for comments
func (f *CommentFetcher) FetchReactionStats(ctx context.Context, entities []domain.Entity) (map[int]domain.ReactionStat, error) {
	url := fmt.Sprintf("http://reactions:8080/stats")

	body, err := json.Marshal(StatsRequest{Data: entities})
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, "reaction body", string(body))
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var reactionStats ReactionsStatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&reactionStats); err != nil {
		return nil, err
	}
	return sliceToMap(reactionStats.Data, func(stat domain.ReactionStat) int {
		return stat.EntityId
	}), nil
}

func sliceToMap[T any](data []T, getID func(T) int) map[int]T {
	dataMap := make(map[int]T)
	for _, tr := range data {
		data[getID(tr)] = tr
	}
	return dataMap
}
