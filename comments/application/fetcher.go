package application

import (
	"bytes"
	"comments/domain"
	"comments/internal"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Request[T any] struct {
	Data T `json:"data"`
}

type Response[T any] struct {
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"` // For error response
}

// Helper struct for fetching data
type CommentFetcher struct {
	userHttpClient      internal.ClientInterface
	reactionsHttpClient internal.ClientInterface
}

func NewCommentFetcher(usersClient internal.ClientInterface, reactionsHttpClient internal.ClientInterface) *CommentFetcher {
	return &CommentFetcher{
		userHttpClient:      usersClient,
		reactionsHttpClient: reactionsHttpClient,
	}
}

// Fetch user details
func (f *CommentFetcher) FetchUsersByID(ctx context.Context, userIDs []int) (map[int]domain.User, error) {

	url := fmt.Sprintf("/by-ids")

	requestBody, err := json.Marshal(Request[[]int]{Data: userIDs})
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(os.Stdout, "user body", string(requestBody))
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	var res Response[[]domain.User]
	if err := f.userHttpClient.GetJSON(req, &res); err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, fmt.Sprintf("user response: %+v", res))
	if res.Message != "" {
		return nil, fmt.Errorf(res.Message)
	}

	return sliceToMap(res.Data, func(u domain.User) int {
		return u.ID
	}), nil
}

// Fetch total reactions for comments
func (f *CommentFetcher) FetchReactionStats(ctx context.Context, entities []domain.Entity) (map[int]domain.ReactionStat, error) {
	url := fmt.Sprintf("/stats")

	body, err := json.Marshal(Request[[]domain.Entity]{Data: entities})
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, "reaction body", string(body))
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var reactionStats Response[[]domain.ReactionStat]
	if err := f.reactionsHttpClient.GetJSON(req, &reactionStats); err != nil {
		return nil, err
	}

	if reactionStats.Message != "" {
		return nil, fmt.Errorf(reactionStats.Message)
	}
	return sliceToMap(reactionStats.Data, func(stat domain.ReactionStat) int {
		return stat.EntityId
	}), nil
}

// Fetch total reactions for comments
func (f *CommentFetcher) FetchUsersReactions(ctx context.Context, userId int, entities []domain.Entity) (map[int]domain.Reaction, error) {
	url := fmt.Sprintf("/%v/user-reactions", userId)

	body, err := json.Marshal(Request[[]domain.Entity]{Data: entities})
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, "reaction body", string(body))
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	var urr Response[[]domain.Reaction]
	if err := f.reactionsHttpClient.GetJSON(req, &urr); err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, urr)
	if urr.Message != "" {
		return nil, fmt.Errorf(urr.Message)
	}
	return sliceToMap(urr.Data, func(r domain.Reaction) int {
		return r.EntityId
	}), nil
}

func sliceToMap[T any](data []T, getID func(T) int) map[int]T {
	dataMap := make(map[int]T)
	for _, tr := range data {
		dataMap[getID(tr)] = tr
	}
	return dataMap
}
