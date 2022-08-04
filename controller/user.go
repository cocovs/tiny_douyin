package controller

//用户
import (
	"github.com/cocovs/tiny-douyin/dao/mysql"
	"github.com/cocovs/tiny-douyin/models"
	"github.com/cocovs/tiny-douyin/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//// 用户信息   key 为用户名++密码
//// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]models.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

//注册后检验数据库是否有相同账号，如果有回复已经有相同账号，没有便为其注册
//注册后将用户名生成token并返回给客户端
//注册接口更新2022-8-1
func Register(c *gin.Context) {
	//Query返回参数值，如果不存在的话则返回空字符串""
	username := c.Query("username")
	password := c.Query("password")

	user := new(models.User)
	result := mysql.DB.First(&user, "name = ? AND password = ?", username, password)
	if result.RowsAffected == 0 {
		//未找到 就为他注册
		//获得token
		token, err := jwt.GenRegisteredClaims(username)
		if err != nil {
			c.JSON(http.StatusOK, UserResponse{
				Response: models.Response{StatusCode: 1, StatusMsg: "token.SignedString err"},
			})
		}
		//不存在 添加新用户 返回token
		NowId := NewUser(username, password).Id
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 0},
			UserId:   NowId,
			Token:    token,
		})
	} else {
		//数据库中检验到已有用户名
		c.JSON(http.StatusOK, UserLoginResponse{
			//用户已存在
			Response: models.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	}

}

// Login 登录
func Login(c *gin.Context) {
	//登录接口 拿到用户的用户名和密码生成 token 在 users数据库中搜索
	username := c.Query("username")
	password := c.Query("password")

	user := new(models.User)
	result := mysql.DB.First(&user, "name = ? AND password = ?", username, password)

	if result.RowsAffected == 0 {
		//存在
		//登录重新给一个token
		token, err := jwt.GenRegisteredClaims(username)
		if err != nil {
			c.JSON(http.StatusOK, UserResponse{
				Response: models.Response{StatusCode: 1, StatusMsg: "token.SignedString err"},
			})
		}
		//返回响应
		//fmt.Println("exist = true")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}

}

// UserInfo
// douyin/user/
//登陆后用户鉴权一次再进入主页
//用户信息
func UserInfo(c *gin.Context) {
	//接收到客户端发送过来的token，将token进行解码，无误则进行
	user := new(models.User)
	username, _ := c.Get("username")
	mysql.DB.First(user, "name = ?", username)

	//用户存在
	c.JSON(http.StatusOK, UserResponse{
		Response: models.Response{StatusCode: 0},
		User:     *user,
	})

}
