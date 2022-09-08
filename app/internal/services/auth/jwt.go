package auth

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt"
)

// JWTClaims is a struct of JWT payload
type JWTClaims struct {
	ID string `json:"user_id"`
	jwt.StandardClaims
}

type JWTUser struct {
	JWT string `json:"jwt"`
	ID  int64  `json:"id"`
}

// Parse JWT into JWTUser
func (ju *JWTUser) Parse(jwtToken, jwtSecret string) (*JWTUser, error) {
	if len(jwtToken) <= 0 {
		return nil, errors.New("auth parsing error: JWT is empty")
	}

	token, err := jwt.ParseWithClaims(jwtToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf(`auth error: %w`, err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		id, err := strconv.ParseInt(claims.ID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(`auth error: %w`, err)
		}

		ju = &JWTUser{JWT: jwtToken, ID: id}
		return ju, nil
	}
	return nil, fmt.Errorf(`auth error: %w`, err)
}
