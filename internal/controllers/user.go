package controllers

import (
	"server/internal/models"
	"server/internal/services"
	"server/pkg/response"
	"strconv"

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

// GetUserById 根据ID获取用户信息
func (c *UserController) GetUserById(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, "Invalid user ID"))
		return
	}

	user, err := c.userService.GetUserById(userId)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(404, response.Error(response.StatusNotFound, err.Error()))
			return
		}
		ctx.JSON(500, response.Error(response.StatusInternalError, err.Error()))
		return
	}

	ctx.JSON(200, response.SuccessWithMessage("Get user successfully", user))
}

// GetUserDetail 获取用户详情（包含统计信息）
func (c *UserController) GetUserDetail(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, response.Error(response.StatusBadRequest, "Invalid user ID"))
		return
	}

	userDetail, err := c.userService.GetUserDetail(userId)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(404, response.Error(response.StatusNotFound, err.Error()))
			return
		}
		ctx.JSON(500, response.Error(response.StatusInternalError, err.Error()))
		return
	}

	ctx.JSON(200, response.SuccessWithMessage("Get user detail successfully", userDetail))
}
