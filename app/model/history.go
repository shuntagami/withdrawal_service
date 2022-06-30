package model

const AmountLimit = 100000

type History struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	Amount int `json:"amount"`
}
