package token

import (
	"be-ai/internal/constants"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// Payload contains the payload data of the token
type Payload struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssueAt   time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
}

func NewJWT() *JWTMaker {
	return &JWTMaker{secretKey: "QWE!@#ASD$%^ZXC*()"}
}

func (maker *JWTMaker) Create(id, username, role string, duration time.Duration) (string, *Payload, error) {

	payload := &Payload{
		ID:        id,
		Username:  username,
		Role:      role,
		IssueAt:   time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString([]byte(maker.secretKey))

	return tokenStr, payload, err
}

func (maker *JWTMaker) Verify(token string) (*Payload, error) {
	keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, constants.ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyfunc)
	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, constants.ErrInvalidToken
	}

	return payload, nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return []string{p.Username}, nil
}

func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{p.ExpiredAt}, nil
}

func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{p.IssueAt}, nil
}

func (p *Payload) GetIssuer() (string, error) {
	return p.Username, nil
}

func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{p.IssueAt}, nil
}

func (p *Payload) GetSubject() (string, error) {
	return "JWT", nil
}

func (p *Payload) Valid() bool {
	return time.Now().Before(p.ExpiredAt)
}
