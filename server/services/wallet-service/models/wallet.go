package models

type Wallet struct {
	ID      string  `json:"id"`
	UserId  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

type Transaction struct {
	ID        string  `json:"id"`
	UserId    string  `json:"user_id"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
	Status    string  `json:"status"`
}