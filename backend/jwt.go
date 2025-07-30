package main

/**
import (
    "time"
	"os"
	"github.com/golang-jwt/jwt/v5"
)


func generateJWT(userID, email, guestID string) (string, error) {
    experationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		Email: email,
		GuestID: guestID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(experationTime),
			issuedAt: jwt.NewNumbericDate(time.Now()),
			subject: userId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func parseJWT(jwtToken string) (*Claims, error) {
	claims := &Claims
	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, err
}

**/
