package models

type Response struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type ValueChange struct {
	Amount int `json:"amount"`
}
