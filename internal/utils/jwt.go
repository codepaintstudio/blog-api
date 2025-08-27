package utils

import (
	"errors"
	"time"

	"blog-api/pkg/config"
	
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Type     string `json:"type"` // access 或 refresh
	jwt.RegisteredClaims
}

var (
	ErrTokenExpired     = errors.New("token已过期")
	ErrTokenNotValidYet = errors.New("token尚未生效")
	ErrTokenMalformed   = errors.New("token格式错误")
	ErrTokenInvalid     = errors.New("token无效")
)

// GenerateTokens 生成访问令牌和刷新令牌
func GenerateTokens(userID uint, username, role string) (accessToken, refreshToken string, err error) {
	cfg := config.GetConfig()
	now := time.Now()
	
	// 生成访问令牌
	accessClaims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.JWT.Issuer,
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.JWT.AccessExpire) * time.Second)),
		},
	}
	
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}
	
	// 生成刷新令牌
	refreshClaims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.JWT.Issuer,
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.JWT.RefreshExpire) * time.Second)),
		},
	}
	
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}
	
	return accessToken, refreshToken, nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()
	
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("意外的签名方法")
		}
		return []byte(cfg.JWT.Secret), nil
	})
	
	if err != nil {
		// 根据错误类型返回相应的错误
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		} else {
			return nil, ErrTokenInvalid
		}
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, ErrTokenInvalid
}

// RefreshToken 刷新访问令牌
func RefreshToken(refreshTokenString string) (newAccessToken string, err error) {
	// 解析刷新令牌
	claims, err := ParseToken(refreshTokenString)
	if err != nil {
		return "", err
	}
	
	// 验证令牌类型
	if claims.Type != "refresh" {
		return "", errors.New("无效的刷新令牌")
	}
	
	// 生成新的访问令牌
	cfg := config.GetConfig()
	now := time.Now()
	
	newClaims := Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.JWT.Issuer,
			Subject:   claims.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(cfg.JWT.AccessExpire) * time.Second)),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ValidateTokenType 验证令牌类型
func ValidateTokenType(tokenString, expectedType string) error {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return err
	}
	
	if claims.Type != expectedType {
		return errors.New("令牌类型不匹配")
	}
	
	return nil
}