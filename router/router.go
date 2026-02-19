package router

import (
	"homework-system/handler"
	"homework-system/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// 生产模式设置
	// gin.SetMode(gin.ReleaseMode)

	r := gin.New() // 不使用默认的Logger和Recovery，我们自己加

	// 添加中间件
	r.Use(gin.Logger())           // 日志
	r.Use(gin.Recovery())         // 崩溃恢复
	r.Use(middleware.Cors())      // 跨域
	r.Use(middleware.RateLimit()) // 限流

	// 初始化handler
	userHandler := handler.NewUserHandler()
	homeworkHandler := handler.NewHomeworkHandler()
	submissionHandler := handler.NewSubmissionHandler()

	// 公开接口
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.RefreshToken)
	}

	// 优秀作业 - 公开接口
	r.GET("/api/submissions/excellent", submissionHandler.GetExcellent)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// 需要认证的接口
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// 用户模块
		user := api.Group("/user")
		{
			user.GET("/profile", userHandler.GetProfile)
			user.DELETE("/account", userHandler.DeleteAccount)
		}

		// 作业模块
		homework := api.Group("/homeworks")
		{
			homework.POST("/", middleware.RequireRole("admin"), homeworkHandler.Create)
			homework.GET("/", homeworkHandler.List)
			homework.GET("/:id", homeworkHandler.GetByID)
			homework.PUT("/:id", middleware.RequireRole("admin"), homeworkHandler.Update)
			homework.DELETE("/:id", middleware.RequireRole("admin"), homeworkHandler.Delete)
		}

		// 提交模块
		submission := api.Group("/submissions")
		{
			// 小登接口
			submission.POST("/", middleware.RequireRole("student"), submissionHandler.Submit)
			submission.GET("/my", middleware.RequireRole("student"), submissionHandler.GetMySubmissions)

			// 老登接口
			submission.GET("/department", middleware.RequireRole("admin"), submissionHandler.GetDepartmentSubmissions)
			submission.POST("/review", middleware.RequireRole("admin"), submissionHandler.Review)
		}
	}

	return r
}
