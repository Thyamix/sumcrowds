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
