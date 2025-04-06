package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"posts/domain"
	"posts/internal"
)

type PostsFetcher struct {
	reactionsClient internal.ClientInterface
	commentsClient  internal.ClientInterface
}

func NewPostsFetcher(reactionsClient internal.ClientInterface, commentsClient internal.ClientInterface) PostsFetcher {
	return PostsFetcher{
		reactionsClient: reactionsClient,
		commentsClient:  commentsClient,
	}
}

// Fetch total reactions for comments
func (f *PostsFetcher) FetchUserReactions(ctx context.Context, userId int, entities []domain.Entity) (map[int]domain.Reaction, error) {
	url := fmt.Sprintf("/%v/user-reactions", userId)

	body, err := json.Marshal(domain.Request[[]domain.Entity]{Data: entities})
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, "reaction body", string(body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var urr domain.Response[[]domain.Reaction]
	f.reactionsClient.GetJSON(req, &urr)

	fmt.Fprintln(os.Stdout, urr)
	if urr.Message != "" {
		return nil, fmt.Errorf(urr.Message)
	}
	return sliceToMap(urr.Data, func(r domain.Reaction) int {
		return r.EntityId
	}), nil
}

// Fetch total reactions for comments
func (f *PostsFetcher) FetchReactionStats(ctx context.Context, entities []domain.Entity) (map[int]domain.ReactionStat, error) {
	url := "/stats"

	body, err := json.Marshal(domain.Request[[]domain.Entity]{Data: entities})
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(os.Stdout, "reaction body", string(body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var reactionStats domain.Response[[]domain.ReactionStat]

	err = f.reactionsClient.GetJSON(req, &reactionStats)
	if err != nil {
		return nil, err
	}
	if reactionStats.Message != "" {
		return nil, fmt.Errorf(reactionStats.Message)
	}
	return sliceToMap(reactionStats.Data, func(stat domain.ReactionStat) int {
		return stat.EntityId
	}), nil
}

func (f *PostsFetcher) FetchCommentsCount(ctx context.Context, entities []domain.Entity) (map[int]domain.Comment, error) {
	url := "/count"
	body, err := json.Marshal(domain.Request[[]domain.Entity]{Data: entities})
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(os.Stdout, "reaction body", string(body))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var commentsCounts domain.Response[[]domain.Comment]
	err = f.commentsClient.GetJSON(req, &commentsCounts)
	fmt.Fprintln(os.Stdout, err.Error())
	if err != nil {
		return nil, err
	}

	return sliceToMap(commentsCounts.Data, func(c domain.Comment) int {
		return c.EntityID
	}), nil
}

func sliceToMap[T any](data []T, getID func(T) int) map[int]T {
	dataMap := make(map[int]T)
	for _, tr := range data {
		dataMap[getID(tr)] = tr
	}
	return dataMap
}
