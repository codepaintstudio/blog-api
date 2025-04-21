package controllers

import (
	"server/internal/models"
	"server/internal/services"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		userService: services.NewUserService(db),
	}
}

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context) {
	var request models.RegisterRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, "Invalid request format"))
		return
	}

	baseUser, err := c.userService.Register(&request)
	if err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, err.Error()))
		return
	}

	ctx.JSON(200, response.Success(baseUser))
}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var request models.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, "Invalid request format"))
		return
	}

	data, err := c.userService.Login(&request)
	if err != nil {
		ctx.JSON(401, response.Error(response.StatusUnauthorized, err.Error()))
		return
	}

	ctx.JSON(200, response.Success(data))
}
