package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
)

type depositRequest struct {
	Project uint
	Amount  string
}

type deposit struct {
	project uint
	amount  float64
	trackID string
}

func newDepositFromRequest(r *http.Request) (deposit, error) {
	var depositReq depositRequest
	if err := json.NewDecoder(r.Body).Decode(&depositReq); err != nil {
		return deposit{}, err
	}
	r.Body.Close()
	if depositReq.Amount == "" || !strings.HasSuffix(depositReq.Amount, ".00") || len(depositReq.Amount) > 8 {
		return deposit{}, fmt.Errorf("amount wrong format: %s", depositReq.Amount)
	}
	amount, err := strconv.ParseFloat(depositReq.Amount, 64)
	if err != nil {
		return deposit{}, err
	}
	trackID, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("could not generate uuid: %v", err)
		return deposit{}, nil
	}
	return deposit{
		amount:  amount,
		project: depositReq.Project,
		trackID: trackID.String(),
	}, nil
}

// type execute struct {
// 	TrackID string `json:"trackID"`
// }

type execute struct {
	ProjectID uint `json:"projectID"`
}

func newExecuteFromRequest(r *http.Request) (execute, error) {
	var executeReq execute
	if err := json.NewDecoder(r.Body).Decode(&executeReq); err != nil {
		return execute{}, err
	}
	r.Body.Close()
	return executeReq, nil
}
