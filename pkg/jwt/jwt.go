package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte("access_secret_key_2024_change_this_in_production")
	refreshSecret = []byte("refresh_secret_key_2024_change_this_in_production")
	accessExpire  = time.Minute * 15   // 15分钟
	refreshExpire = time.Hour * 24 * 7 // 7天
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// 自定义错误
var (
	ErrTokenExpired = errors.New("token已过期")
	ErrTokenInvalid = errors.New("无效的token")
)

// GenerateTokens 生成双Token
func GenerateTokens(userID uint, role string) (accessToken, refreshToken string, err error) {
	// Access Token
	accessClaims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "homework-system",
		},
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = access.SignedString(accessSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "homework-system",
		},
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refresh.SignedString(refreshSecret)
	return accessToken, refreshToken, err
}

// ParseAccessToken 解析Access Token
func ParseAccessToken(tokenString string) (*Claims, error) {
	return parseToken(tokenString, accessSecret)
}

// ParseRefreshToken 解析Refresh Token
func ParseRefreshToken(tokenString string) (*Claims, error) {
	return parseToken(tokenString, refreshSecret)
}

func parseToken(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		// 判断是否是过期错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

// RefreshTokens 刷新Token
func RefreshTokens(refreshTokenString string) (string, string, error) {
	claims, err := ParseRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}
	return GenerateTokens(claims.UserID, claims.Role)
}
