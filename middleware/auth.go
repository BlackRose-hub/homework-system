package middleware

import (
	"homework-system/pkg/errcode"
	"homework-system/pkg/jwt"
	"homework-system/pkg/response"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, errcode.Unauthorized)
			c.Abort()
			return
		}

		// 检查格式：Bearer xxxxx
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, http.StatusUnauthorized, errcode.TokenInvalid)
			c.Abort()
			return
		}

		// 解析token
		claims, err := jwt.ParseAccessToken(parts[1])
		if err != nil {
			// 判断是否是过期错误
			if err == jwt.ErrTokenExpired {
				response.Error(c, http.StatusUnauthorized, errcode.TokenExpired)
			} else {
				response.Error(c, http.StatusUnauthorized, errcode.TokenInvalid)
			}
			c.Abort()
			return
		}

		// 检查token是否快过期（剩余时间少于5分钟）
		if claims.ExpiresAt != nil {
			expireTime := claims.ExpiresAt.Time
			if time.Until(expireTime) < 5*time.Minute {
				// 在Header中提示前端需要刷新token
				c.Header("X-Token-Expiring", "true")
			}
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

// RequireRole 角色权限中间件
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, errcode.Unauthorized)
			c.Abort()
			return
		}

		// 检查角色是否在允许列表中
		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, errcode.PermissionDenied)
		c.Abort()
	}
}

// RateLimit 简单的限流中间件（防止暴力请求）
func RateLimit() gin.HandlerFunc {
	// 使用map存储每个IP的请求次数
	limiter := make(map[string]int)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter[ip]++

		// 每秒超过30次请求就限流
		if limiter[ip] > 30 {
			response.Error(c, http.StatusTooManyRequests, 42901)
			c.Abort()
			return
		}

		// 每秒重置
		time.AfterFunc(time.Second, func() {
			limiter[ip]--
		})

		c.Next()
	}
}

// Cors 跨域中间件（给前端用）
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
