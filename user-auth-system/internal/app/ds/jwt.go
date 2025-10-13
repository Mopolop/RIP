package ds

import (
	"github.com/golang-jwt/jwt"
	"time"
	"user-auth-system/internal/app/role"
)

type JWTClaims struct {
	jwt.StandardClaims
	UserID int       `json:"user_id"`
	Role   role.Role `json:"role"`
}

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResp struct {
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
}
