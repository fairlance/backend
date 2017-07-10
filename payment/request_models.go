package payment

import (
	"encoding/json"
	"fmt"
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
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
	trackID, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("could not generate uuid: %v", err)
		return deposit{}, nil
	}
	return deposit{
		projectID: depositReq.ProjectID,
		trackID:   trackID.String(),
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
