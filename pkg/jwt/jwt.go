package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	ID        int64
	Login     string
	SessionID string
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data *JWTData) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        data.ID,
		"login":     data.Login,
		"sessionID": data.SessionID,
		"exp":       jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return s, nil
}

func (j *JWT) Parse(tokenString string) (*JWTData, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return nil, errors.New("invalid 'id' claim")
	}

	id := int64(idFloat)

	login, ok := claims["login"].(string)
	if !ok {
		return nil, errors.New("invalid 'login' claim")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return nil, errors.New("token expired")
		}
	}

	return &JWTData{
		ID:    id,
		Login: login,
	}, nil
}
