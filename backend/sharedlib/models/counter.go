package models

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

type FestivalAccess struct {
	Id         int64 `sql:"id"`
	FestivalId int64 `sql:"festival_id"`
	UserId     int64 `sql:"user_id"`
	LastUsedAt int64 `sql:"last_used_at"`
	Revoked    bool  `sql:"revoked"`
}
