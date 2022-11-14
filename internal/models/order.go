package models

import "time"

type Order struct {
	UserLogin   string
	Status      string    `json:"status"`
	Accrual     float64   `json:"accrual,omitempty"`
	Created     time.Time `json:"uploaded_at"`
	NumberOrder int       `json:"number"`
}
