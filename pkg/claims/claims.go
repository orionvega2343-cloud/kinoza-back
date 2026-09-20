package claims

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserId string
	Email  string
	Role   string
	jwt.RegisteredClaims
}
