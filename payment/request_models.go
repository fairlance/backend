package payment

import (
	"encoding/json"
	"net/http"

	"fmt"
	"time"
)

type depositRequest struct {
	ProjectID uint
}

type deposit struct {
	projectID uint
	trackID   string
}

func newDepositFromRequest(r *http.Request) (deposit, error) {
	var depositReq depositRequest
	if err := json.NewDecoder(r.Body).Decode(&depositReq); err != nil {
		return deposit{}, err
	}
	r.Body.Close()
	trackID := timeToMillis(time.Now())
	return deposit{
		projectID: depositReq.ProjectID,
		trackID:   fmt.Sprintf("%d", trackID),
	}, nil
}

type executeRequest struct {
	ProjectID uint `json:"projectID"`
}

type execute struct {
	projectID uint
}

func newExecuteFromRequest(r *http.Request) (*execute, error) {
	var executeReq executeRequest
	if err := json.NewDecoder(r.Body).Decode(&executeReq); err != nil {
		return nil, err
	}
	r.Body.Close()
	return &execute{
		projectID: executeReq.ProjectID,
	}, nil
}

func timeToMillis(t time.Time) int64 {
	return t.UnixNano() / 1000000
}
