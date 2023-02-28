package security

import (
	"strings"
	"time"
	"websocket-service/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT ...
func GenerateJWT(m map[string]interface{}, tokenExpireTime time.Duration, tokenSecretKey string) (tokenString string, err error) {
	defer errors.WrapCheck(&err, "GenerateJWT")

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	for key, value := range m {
		claims[key] = value
	}

	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(tokenExpireTime).Unix()

	tokenString, err = token.SignedString([]byte(tokenSecretKey))
	if err != nil {
		return "", errors.Wrap(err, "token.SignedString")
	}

	return tokenString, nil
}

// ExtractClaims extracts claims from given token
func ExtractClaims(tokenString string, tokenSecretKey string) (claims jwt.MapClaims, err error) {
	defer errors.WrapCheck(&err, "ExtractClaims")
	var (
		token *jwt.Token
	)

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return []byte(tokenSecretKey), nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "jwt.Parse")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractToken checks and returns token part of input string
func ExtractToken(bearer string) (token string, err error) {
	defer errors.WrapCheck(&err, "ExtractToken")
	strArr := strings.Split(bearer, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}
	return token, errors.New("wrong token format")
}

type TokenInfo struct {
	Id string
}

func ParseClaims(token string, secretKey string) (resp TokenInfo, err error) {
	defer errors.WrapCheck(&err, "ParseClaims")
	var (
		ok     bool
		claims jwt.MapClaims
	)
	resp = TokenInfo{}

	claims, err = ExtractClaims(token, secretKey)
	if err != nil {
		return resp, err
	}

	resp.Id, ok = claims["user_id"].(string)
	if !ok {
		err = errors.New("cannot parse 'user_id' field")
		return resp, err
	}
	resp.Id, ok = claims["email"].(string)
	if !ok {
		err = errors.New("cannot parse 'email' field")
		return resp, err
	}

	return
}
