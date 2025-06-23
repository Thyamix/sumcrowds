package models

type Response struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type FestivalData struct {
	Id           int    `sql:"id"`
	Code         string `sql:"code"`
	Pin          int    `sql:"pin" json:"adminPin"`
	Password     string `json:"password"`
	PasswordHash string `sql:"password"`
	CreatedAt    int    `sql:"created_at"`
	ExpiresAt    int    `sql:"expires_at"`
}

type Event struct {
	Id         int `sql:"id"`
	FestivalId int `sql:"festival_id"`
	CreatedAt  int `sql:"created_at"`
	LastUsedAt int `sql:"last_used_at"`
	Active     int `sql:"active"`
}

type ValueChange struct {
	Amount int `json:"amount"`
}
