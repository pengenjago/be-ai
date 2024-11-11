package dto

import "time"

type UserQuery struct {
	PageNo   int    `query:"pageNo"`
	PageSize int    `query:"pageSize"`
	Search   string `query:"search"`
}

type UserReq struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	NoTelephone string `json:"noTelephone"`
	Password    string `json:"password"`
	Auth        string `json:"auth"`
}

type UserRes struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	NoTelephone string `json:"noTelephone,omitempty"`
	Role        string `json:"role,omitempty"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRes struct {
	AccessToken                 string    `json:"accessToken"`
	AccessTokenExpiredAt        time.Time `json:"accessTokenExpiredAt"`
	RefreshAccessToken          string    `json:"refreshAccessToken"`
	RefreshAccessTokenExpiredAt time.Time `json:"refreshAccessTokenExpiredAt"`
	UserRes                     UserRes   `json:"user"`
}
