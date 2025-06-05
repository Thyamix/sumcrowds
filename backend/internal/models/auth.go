package models

type AccessToken struct {
	Id        int    `sql:"id"`
	Token     string `sql:"token"`
	ExpiresAt int64  `sql:"expires_at"`
	UserId    int    `sql:"user_id"`
}

type RefreshToken struct {
	Id         int    `sql:"id"`
	Token      string `sql:"token"`
	LastUsedAt int64  `sql:"last_used_at"`
	ExpiresAt  int64  `sql:"expires_at"`
	UserId     int    `sql:"user_id"`
	Revoked    int    `sql:"revoked"`
}

type FestivalAccess struct {
	Id         int   `sql:"id"`
	FestivalId int   `sql:"festival_id"`
	UserId     int   `sql:"user_id"`
	LastUsedAt int64 `sql:"last_used_at"`
}
