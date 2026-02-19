package handler

import (
	"homework-system/models"
	"homework-system/pkg/errcode"
	"homework-system/pkg/response"
	"homework-system/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HomeworkHandler struct { //作业处理器
	homeworkService *service.HomeworkService
}

func NewHomeworkHandler() *HomeworkHandler { //创建实例
	return &HomeworkHandler{
		homeworkService: service.NewHomeworkService(),
	}
}

func (h *HomeworkHandler) Create(c *gin.Context) { //发布作业接口// POST /api/homeworks，需要管理员权限
	var req service.CreateHomeworkRequest // 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, errcode.InvalidParam)
		return
	}

	creatorID, _ := c.Get("user_id") // 获取操作者ID
	// 调用service创建作业
	homework, code, err := h.homeworkService.Create(creatorID.(uint), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code)
		return
	}
	if code != errcode.Success {
		response.Error(c, http.StatusBadRequest, code)
		return
	}
	// 返回创建的作业信息
	response.Success(c, gin.H{
		"id":               homework.ID,
		"title":            homework.Title,
		"description":      homework.Description,
		"department":       homework.Department,
		"department_label": models.DepartmentLabel[homework.Department],
		"creator": gin.H{
			"id":       homework.Creator.ID,
			"nickname": homework.Creator.Nickname,
		},
		"deadline":   homework.Deadline,
		"allow_late": homework.AllowLate,
		"version":    homework.Version,
		"created_at": homework.CreatedAt,
	})
}

func (h *HomeworkHandler) GetByID(c *gin.Context) { //获取作业详情接口，// GET /api/homeworks/:id
	id, err := strconv.ParseUint(c.Param("id"), 10, 32) // 解析URL参数，10进制，32位
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.InvalidParam)
		return
	}
	// 调用service获取作业
	homework, code, err := h.homeworkService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code)
		return
	}
	if code != errcode.Success {
		response.Error(c, http.StatusNotFound, code)
		return
	}
	// 返回作业详情（带部门标签）
	response.Success(c, gin.H{
		"id":               homework.ID,
		"title":            homework.Title,
		"description":      homework.Description,
		"department":       homework.Department,
		"department_label": models.DepartmentLabel[homework.Department],
		"creator": gin.H{
			"id":               homework.Creator.ID,
			"nickname":         homework.Creator.Nickname,
			"department":       homework.Creator.Department,
			"department_label": models.DepartmentLabel[homework.Creator.Department],
		},
		"deadline":   homework.Deadline,
		"allow_late": homework.AllowLate,
		"version":    homework.Version,
		"created_at": homework.CreatedAt,
		"updated_at": homework.UpdatedAt,
	})
}

func (h *HomeworkHandler) List(c *gin.Context) { //作业列表接口// GET /api/homeworks?department=backend&page=1&page_size=10
	// 获取查询参数
	department := c.Query("department")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1")) //默认第一页
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 调用service - 这里只有3个返回值
	result, code, err := h.homeworkService.List(department, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code)
		return
	}
	if code != errcode.Success {
		response.Error(c, http.StatusBadRequest, code)
		return
	}

	response.Success(c, result)
}

// Update 修改作业
func (h *HomeworkHandler) Update(c *gin.Context) {
	// 获取路径参数
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.InvalidParam)
		return
	}

	// 绑定请求体
	var req service.UpdateHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ 参数绑定错误: %v", err)
		response.Error(c, http.StatusBadRequest, errcode.InvalidParam)
		return
	}
	req.ID = uint(id)

	// 调用service
	operatorID, _ := c.Get("user_id")
	homework, code, err := h.homeworkService.Update(operatorID.(uint), &req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code)
		return
	}

	// 处理业务错误
	if code != errcode.Success {
		httpCode := http.StatusBadRequest
		switch code {
		case errcode.HomeworkVersionExpired:
			httpCode = http.StatusConflict
		case errcode.HomeworkNotPermission:
			httpCode = http.StatusForbidden
		case errcode.HomeworkNotFound:
			httpCode = http.StatusNotFound
		}
		response.Error(c, httpCode, code)
		return
	}

	response.Success(c, gin.H{
		"id":          homework.ID,
		"title":       homework.Title,
		"description": homework.Description,
		"deadline":    homework.Deadline,
		"allow_late":  homework.AllowLate,
		"version":     homework.Version,
		"updated_at":  homework.UpdatedAt,
	})
}

// Delete 删除作业
func (h *HomeworkHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, errcode.InvalidParam)
		return
	}

	operatorID, _ := c.Get("user_id")
	code, err := h.homeworkService.Delete(operatorID.(uint), uint(id))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, code)
		return
	}
	if code != errcode.Success {
		httpCode := http.StatusBadRequest
		switch code {
		case errcode.HomeworkNotPermission:
			httpCode = http.StatusForbidden
		case errcode.HomeworkNotFound:
			httpCode = http.StatusNotFound
		}
		response.Error(c, httpCode, code)
		return
	}

	response.Success(c, gin.H{
		"message": "作业删除成功",
	})
}
