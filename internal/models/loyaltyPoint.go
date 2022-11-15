package models

// LoyaltyPoint
//
//	{
//	   "order": "<number>",
//	   "status": "PROCESSED",
//	   "accrual": 500
//	}
type LoyaltyPoint struct {
	Status      string  `json:"status"`
	Accrual     float64 `json:"accrual,omitempty"`
	NumberOrder string  `json:"order"`
}
