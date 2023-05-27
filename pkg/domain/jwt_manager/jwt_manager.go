package jwt_manager

import (
	"blazem/pkg/domain/user"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTManager interface {
	CreateJWT(user *user.User, expirationDate time.Time) (string, error)
	ParseJWT(jwtStr string) (*jwt.MapClaims, bool)
	JwtInDate(claims jwt.MapClaims) bool
	GetCurrentUserId(jwtStr string) (string, error)
}

type JWTManage struct {
	SecretKey []byte
}

func NewJWTManager(secretKey []byte) *JWTManage {
	return &JWTManage{
		SecretKey: secretKey,
	}
}

// createJWT creates a JWT and returns the token and an error if there is one
func (jwtMgr *JWTManage) CreateJWT(user *user.User, expirationDate time.Time) (string, error) {
	t := jwt.New(jwt.SigningMethodHS256)

	claims := t.Claims.(jwt.MapClaims)
	claims["exp"] = expirationDate.Unix()
	claims["user"] = user.Id

	jwt, err := t.SignedString(jwtMgr.SecretKey)
	return jwt, err
}

// parseJWT checks if the JWT is valid and in date
func (jwtMgr *JWTManage) ParseJWT(jwtStr string) (*jwt.MapClaims, bool) {
	tkn, err := jwt.Parse(jwtStr, func(t *jwt.Token) (interface{}, error) {
		return jwtMgr.SecretKey, nil
	})
	if err != nil {
		return nil, false
	}

	if tkn.Method != jwt.SigningMethodHS256 {
		return nil, false
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	inDate := jwtMgr.JwtInDate(claims)
	if !ok || !inDate {
		return nil, false
	}

	return &claims, true
}

// GetCurrentUserId returns the Id of the user who is logged in to
// this session
func (jwtMgr *JWTManage) GetCurrentUserId(jwtStr string) (string, error) {
	parseClaims, valid := jwtMgr.ParseJWT(jwtStr)
	if !valid {
		return "", errors.New("JWT not valid")
	}

	claims := *parseClaims
	userId := claims["user"].(string)
	return userId, nil
}

// jwtInDate checks whether the JWT (claims) section is in date
func (jwtMgr *JWTManage) JwtInDate(claims jwt.MapClaims) bool {
	if int64(claims["exp"].(float64)) < time.Now().Local().Unix() {
		return false
	}
	return true
}
