package models

type User struct {
	Id       int    `gorm:"primarykey;column:id" json:"id"`
	Username string `gorm:"column:username" json:"username"`
	Email    string `gorm:"column:email" json:"email"`
	Password string `gorm:"column:password" json:"password"`
}

// 基础用户信息返回
type BaseUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// 用户注册Request
type RegisterRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"re_password"`
}

// 用户登录Request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// 用户登录response
type LoginResponse struct {
	Token string   `json:"token"`
	User  BaseUser `json:"user"`
}
