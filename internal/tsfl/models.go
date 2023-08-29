package tsfl

import (
	"fmt"
	"github.com/google/uuid"
)

// InputMoneyRequest save money in a penny
type InputMoneyRequest struct {
	UserId string `json:"userId"`
	Count  uint64 `json:"count"`
}

type OutputMoneyRequest struct {
	UserId string `json:"userId"`
	Count  uint64 `json:"count"`
	Url    string `json:"url"`
}

type TransactionInfo struct {
	TransactionId uuid.UUID `json:"transactionId"`
	UserId        uuid.UUID `json:"userId"`
	Count         uint64    `json:"count"`
	Input         bool      `json:"input	"`
	Confirmed     bool      `json:"confirmed"`
	Info          string    `json:"info"`
}

type RespondError struct {
	Status  int
	Message string
	Err     error
}

func (r *RespondError) Error() string {
	return fmt.Sprintf("Status: %d\nMessage: %s\nError: %s", r.Status, r.Message, r.Err.Error())
}
