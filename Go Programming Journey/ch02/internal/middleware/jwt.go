package middleware

import (
	"demo/ch02/pkg/app"
	"demo/ch02/pkg/errcode"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token string
			ecode = errcode.Success
		)
		// 获取 token
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}
		if token == "" {
			ecode = errcode.InvalidParams
		} else {
			// ParseToken() 解析 token
			_, err := app.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					ecode = errcode.UnauthorizedTokenTimeout
				default:
					ecode = errcode.UnauthorizedTokenError
				}
			}
		}
		if ecode != errcode.Success {
			response := app.NewResponse(c)
			response.ToErrorResponse(ecode)
			/*
				Abort() 可防止调用挂起的处理程序
				请注意，这不会停止当前处理程序
				假设您有一个授权中间件来验证当前请求是否已获得授权
				如果授权失败（例如：密码不匹配）
				请调用 Abort 以确保不调用此请求的其余处理程序
			*/
			c.Abort()
			return
		}
		c.Next()
	}
}
