package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/evenq/evenq-cli/src/shared/api"
	"gopkg.in/guregu/null.v3"
)

type Import struct {
	ID          string            `json:"id"`
	OrgID       string            `json:"orgId"`
	EventID     string            `json:"eventId"`
	Status      string            `json:"status"`
	Mappings    map[string]string `json:"mappings"`
	CreatedAt   time.Time         `json:"createdAt"`
	StartedAt   null.Time         `json:"startedAt"`
	CompletedAt null.Time         `json:"completedAt"`
	EventCount  null.Int          `json:"eventCount"`
	FileName    null.String       `json:"fileName"`
	FileHash    null.String       `json:"fileHash"`
	FileSize    null.Int          `json:"fileSize"`
	FileFormat  null.String       `json:"fileFormat"`
	UploadURL   null.String       `json:"uploadUrl"`
}

func CreateImport(ctx context.Context, data Import) (Import, bool) {
	out := Import{}

	if data.EventID == "" {
		log.Println("missing event ID on import creation")
		return out, false
	}

	path := fmt.Sprintf("/events/%v/imports", data.EventID)

	_, err := api.Post(ctx, path, data, &out)

	return out, out.ID != "" && err == nil
}

func StartImport(ctx context.Context, eventID string, importID string, mapping map[string]string) bool {
	out := map[string]interface{}{}

	data := map[string]interface{}{
		"mappings": mapping,
	}

	path := fmt.Sprintf("/events/%v/imports/%v/start", eventID, importID)

	_, err := api.Post(ctx, path, data, &out)
	if err != nil {
		fmt.Println(err)
	}

	return api.IsSuccess(out)
}
