package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type depositRequest struct {
	Project      uint   `json:"project"`
	AmountString string `json:"amount"`
	amount       float64
	trackingID   string
}

func newDepositRequest(r *http.Request) (*depositRequest, error) {
	var deposit depositRequest
	if err := json.NewDecoder(r.Body).Decode(&deposit); err != nil {
		return nil, err
	}
	r.Body.Close()
	if deposit.AmountString == "" || !strings.HasSuffix(deposit.AmountString, ".00") || len(deposit.AmountString) > 8 {
		return nil, fmt.Errorf("amount wrong format: %s", deposit.AmountString)
	}
	amount, err := strconv.ParseFloat(deposit.AmountString, 64)
	if err != nil {
		return nil, err
	}
	deposit.amount = amount
	deposit.trackingID = "generate a unque tracking ID of 127 chars"
	return &deposit, nil
}

type executeRequest struct {
	Project uint `json:"project"`
}

func newExecuteRequest(r *http.Request) (*executeRequest, error) {
	var execute executeRequest
	if err := json.NewDecoder(r.Body).Decode(&execute); err != nil {
		return nil, err
	}
	r.Body.Close()
	return &execute, nil
}
