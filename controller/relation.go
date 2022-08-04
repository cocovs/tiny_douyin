package controller

//关系
import (
	"github.com/cocovs/tiny-douyin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserListResponse struct {
	models.Response
	UserList []models.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, models.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

//关注列表
// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	var userList []models.User
	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}

//粉丝列表
// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	var userList []models.User
	c.JSON(http.StatusOK, UserListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
}
