package models

type AccessToken struct {
	Id        int64  `sql:"id"`
	Token     string `sql:"token"`
	ExpiresAt int64  `sql:"expires_at"`
	UserId    int64  `sql:"user_id"`
	Revoked   bool   `sql:"revoked"`
}

type RefreshToken struct {
	Id        int64  `sql:"id"`
	Token     string `sql:"token"`
	ExpiresAt int64  `sql:"expires_at"`
	UserId    int64  `sql:"user_id"`
	Revoked   bool   `sql:"revoked"`
}

type FestivalAccess struct {
	Id         int64 `sql:"id"`
	FestivalId int64 `sql:"festival_id"`
	UserId     int64 `sql:"user_id"`
	LastUsedAt int64 `sql:"last_used_at"`
	Revoked    bool  `sql:"revoked"`
}
