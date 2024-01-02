package jwtx

import "github.com/golang-jwt/jwt"

func GenJwtToken(secret string, nowDate int64, accessExpire int64, uid interface{}) (jwtToken string, err error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = nowDate + accessExpire
	claims["iat"] = nowDate
	claims["id"] = uid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secret))
}
