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
func (f *PostsFetcher) FetchUsersReactions(ctx context.Context, userId int, entities []domain.Entity) (map[int]domain.Reaction, error) {
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

func sliceToMap[T any](data []T, getID func(T) int) map[int]T {
	dataMap := make(map[int]T)
	for _, tr := range data {
		dataMap[getID(tr)] = tr
	}
	return dataMap
}
