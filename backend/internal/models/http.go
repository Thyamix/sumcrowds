package models

type Response struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type FestivalData struct {
	Id           int64  `sql:"id"`
	Code         string `sql:"code"`
	Pin          string `sql:"pin" json:"pin"`
	Password     string `json:"password"`
	PasswordHash string `sql:"password"`
	CreatedAt    int64  `sql:"created_at"`
	LastUsedAt   int64  `sql:"last_used_at"`
}

type Event struct {
	Id         int64 `sql:"id"`
	FestivalId int64 `sql:"festival_id"`
	CreatedAt  int64 `sql:"created_at"`
	LastUsedAt int64 `sql:"last_used_at"`
	Active     bool  `sql:"active"`
}

type ValueChange struct {
	Amount int `json:"amount"`
}
