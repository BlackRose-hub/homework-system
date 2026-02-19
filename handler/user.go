package handler

import (
	"homework-system/models"
	"homework-system/pkg/errcode"
	"homework-system/pkg/response"
	"homework-system/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct { //用户处理器
	userService *service.UserService //依赖service层
}

func NewUserHandler() *UserHandler { //创建实例
	return &UserHandler{
		userService: service.NewUserService(), //初始化service
	}
}

func (h *UserHandler) Register(c *gin.Context) { //用户注册接口//POST /api/auth/register
	var req service.RegisterRequest                //绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil { //从请求体解析JSON到结构体
		response.Error(c, http.StatusBadRequest, errcode.InvalidParam) //参数绑定失败就返回参数错误
		return
	}

	user, code, err := h.userService.Register(&req) //调用service层处理业务逻辑
	if err != nil {                                 //服务器内部错误
		response.Error(c, http.StatusInternalServerError, code)
		return
	}
	if code != errcode.Success { //业务错误
		response.Error(c, http.StatusBadRequest, code)
		return
	}

	response.Success(c, gin.H{ //返回成功响应
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"department": user.Department,
	})
}

func (h *UserHandler) Login(c *gin.Context) { //用户登录接口//POST /api/auth/login
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.InvalidParam)
		return
	}

	resp, code, err := h.userService.Login(&req) //调用service登录
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code)
		return
	}
	if code != errcode.Success { //登录失败
		response.Error(c, http.StatusUnauthorized, code)
		return
	}

	response.Success(c, resp) //返回token和用户信息
}

func (h *UserHandler) RefreshToken(c *gin.Context) { //刷新Token接口，//POST /api/auth/refresh
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	} //需要refresh_token参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.InvalidParam)
		return
	}
	//调用service刷新token
	accessToken, refreshToken, code, err := h.userService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, code)
		return
	}
	//返回新的token
	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) { //获取用户信息接口，//GET /api/user/profile
	userID, _ := c.Get("user_id") //从上下文获取用户ID
	//调用service获取用户信息
	user, code, err := h.userService.GetProfile(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code)
		return
	}
	if code != errcode.Success {
		response.Error(c, http.StatusNotFound, code)
		return
	}
	//返回用户信息，带部门标签
	response.Success(c, gin.H{
		"id":               user.ID,
		"username":         user.Username,
		"nickname":         user.Nickname,
		"role":             user.Role,
		"department":       user.Department,
		"department_label": models.DepartmentLabel[user.Department],
		"email":            user.Email,
		"created_at":       user.CreatedAt,
	})
}

func (h *UserHandler) DeleteAccount(c *gin.Context) { //注销账号接口，//DELETE /api/user/account
	userID, _ := c.Get("user_id")

	code, err := h.userService.DeleteAccount(userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code)
		return
	}

	response.Success(c, nil) //成功删除，不返回数据
}
