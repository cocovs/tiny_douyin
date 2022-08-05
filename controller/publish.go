package controller

//发布
import (
	"github.com/cocovs/tiny-douyin/dao/mysql"
	"github.com/cocovs/tiny-douyin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	models.Response
	VideoList []models.Video `json:"video_list"`
}

// Publish 用户发布视频
func Publish(c *gin.Context) {
	//token := c.PostForm("token")

	user := new(models.User)
	username, _ := c.Get("username")

	result := mysql.DB.First(&user, "name = ?", username)

	//拿到用户信息
	//UserOK, exist := QueryUsersToken(token)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	//拿到文件
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	myfile, _ := file.Open()

	FileName := file.Filename
	//上传视频
	err = UpOSS(myfile, FileName)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		//上传视频信息至mysql
		NewVideo(file.Filename, user.Id)

		c.JSON(http.StatusOK, models.Response{
			StatusCode: 0,
			StatusMsg:  file.Filename + " uploaded successfully  wuhu！！！！",
		})
	}

}

// PublishList 根据每个用户的信息来返回视频列表 请求参数token user_id
func PublishList(c *gin.Context) {
	//c.Query()
	userId := c.Query("user_id")
	User := new(models.User)
	mysql.DB.First(User, "id = ?", userId)
	//fmt.Println("publish User: ", User)

	var VideoList []models.Video
	VideoList = UserVideoList(User)
	//fmt.Println("publish video: ", VideoList)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: VideoList,
	})
}
