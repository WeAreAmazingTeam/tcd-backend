package auth

import (
	"errors"

	"github.com/WeAreAmazingTeam/tcd-backend/constant"
	"github.com/dgrijalva/jwt-go"
)

func (svc *authService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	claim["the_cloud_donation_user_id"] = userID

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, err := token.SignedString(constant.SecretKey)

	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (svc *authService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Token invalid")
		}
		return []byte(constant.SecretKey), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}
