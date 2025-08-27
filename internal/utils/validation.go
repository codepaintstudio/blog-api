package utils

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// IsEmail 检查字符串是否为有效邮箱格式
func IsEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) == 0 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

// IsValidUsername 检查用户名是否有效
func IsValidUsername(username string) bool {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	
	// 用户名只能包含字母、数字、下划线和连字符
	validChars := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return validChars.MatchString(username)
}