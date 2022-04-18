package events

import (
	"context"
	"time"

	"github.com/evenq/evenq-cli/src/shared/api"
)

type Event struct {
	ID         string     `json:"id"`
	OwnerID    string     `json:"ownerId"`
	IsFav      bool       `json:"isFavorite"`
	CreatedAt  time.Time  `json:"createdAt"`
	EventStats EventStats `json:"stats"`
}

type EventStats struct {
	TotalCount               int       `json:"totalCount"`
	TotalSize                int       `json:"totalSize"`
	Oldest                   time.Time `json:"oldest"`
	Newest                   time.Time `json:"newest"`
	PartitionCount           int       `json:"partitionCount"`
	PartitionCountIncomplete bool      `json:"partitionCountIncomplete"`
	PartitionList            []string  `json:"partitionList"` // encounters the first 1000 partition keys we find
}

func Get(ctx context.Context, name string) (Event, bool) {
	evt := Event{}

	err := api.Get(ctx, "/events/"+name, &evt)

	return evt, evt.ID != "" && err == nil
}

func Create(ctx context.Context, name string) (Event, bool) {
	evt := Event{
		ID: name,
	}

	out := map[string]interface{}{}

	_, err := api.Post(ctx, "/events", evt, &out)

	return evt, api.IsSuccess(out) && err == nil
}
