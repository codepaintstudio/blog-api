package utils

import (
	"errors"
	
	"golang.org/x/crypto/bcrypt"
)

const (
	// MinPasswordLength 最小密码长度
	MinPasswordLength = 6
	// MaxPasswordLength 最大密码长度
	MaxPasswordLength = 128
)

var (
	ErrPasswordTooShort = errors.New("密码长度不能少于6位")
	ErrPasswordTooLong  = errors.New("密码长度不能超过128位")
)

// HashPassword 对密码进行哈希加密
func HashPassword(password string) (string, error) {
	// 验证密码长度
	if len(password) < MinPasswordLength {
		return "", ErrPasswordTooShort
	}
	if len(password) > MaxPasswordLength {
		return "", ErrPasswordTooLong
	}
	
	// 使用bcrypt加密密码，成本因子为12
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	
	return string(hashedBytes), nil
}

// CheckPassword 验证密码是否匹配
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePassword 验证密码格式
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}
	if len(password) > MaxPasswordLength {
		return ErrPasswordTooLong
	}
	
	// 这里可以添加更多密码复杂度验证
	// 例如：必须包含数字、字母、特殊字符等
	
	return nil
}