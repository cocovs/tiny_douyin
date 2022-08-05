package middlewares

//用户登录以后
//auth 身份验证
import (
	"fmt"
	"github.com/cocovs/tiny-douyin/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {

		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":  2005,
				"msg":   "无效的Token",
				"token": token,
				"err":   err,
			})
			c.Abort()
			return
		}
		fmt.Println("token  ->> mc.Username: ", mc.Username)
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		//c.Set("password", mc.Password)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
