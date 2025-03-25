package auth

import (
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
)

const UserClaim = "user"

type JwtManager struct {
	secretKey string
}

func NewJwtManager(secretKey string) JwtManager {
	return JwtManager{secretKey: secretKey}
}

func (jm *JwtManager) GenerateJwt(userId uint64) (string, error) {
	claims := jwt.MapClaims{"sub": strconv.FormatUint(userId, 10)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jm.secretKey))
	if err != nil {
		return "", TokenSigningFailed
	}
	return signedToken, nil
}

func (jm *JwtManager) ParseJwt(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, InvalidSignature
		}
		return []byte(jm.secretKey), nil
	})
}
