package controller

//收藏
import (
	"fmt"
	"github.com/cocovs/tiny-douyin/dao/mysql"
	"github.com/cocovs/tiny-douyin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//required int64 user_id = 1; // 用户id
//required string token = 2; // 用户鉴权token
//required int64 video_id = 3; // 视频id
//required int32 action_type = 4; // 1-点赞，2-取消点赞

// FavoriteAction 点赞操作
func FavoriteAction(c *gin.Context) {

	//userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	userName, _ := c.Get("username")
	userNow := new(models.User)
	mysql.DB.First(userNow, "name = ?", userName)
	userId := userNow.Id
	fmt.Println("FavoriteAction userId: ", userId)

	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType := c.Query("action_type")
	fmt.Println("videoId: ", videoId)
	fmt.Println("actionType: ", actionType)
	//创建点赞结构体
	userLike := models.User_favorite_video{
		User_id:  userId,
		Video_id: videoId,
	}
	fmt.Println("userLike: ", userLike)
	// 1 点赞
	if actionType == "1" {
		//则在该表中添加一行信息 为该用户id 和 视频id
		result := mysql.DB.Create(userLike)
		if result.Error != nil {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "点赞失败"})
		} else {
			c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "成功"})
		}
	} else {
		//2 取消赞 删除该列信息
		result := mysql.DB.Delete(userLike)
		if result.Error != nil {
			c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "删除失败"})
		} else {
			c.JSON(http.StatusOK, models.Response{StatusCode: 0, StatusMsg: "成功"})
		}
	}

}

// FavoriteList 点赞列表（返回点过赞的所有视频
func FavoriteList(c *gin.Context) {
	var err error
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	fmt.Println("FavoriteList userid: ", userId)
	//从`user_favorite_videos`中搜索该用户id 返回所有的视频id切片
	vIdList := make([]models.User_favorite_video, 0, 0)
	mysql.DB.Find(&vIdList, "user_id = ?", userId)

	fmt.Println("FavoriteList vidlist: ", vIdList)

	//通过视频id切片 在 videos表中获取到所有的视频做成一个切片
	videoList := make([]models.Video, 0, 0)
	videoNow := new(models.Video)
	for _, v := range vIdList {
		//搜索video 并加入是否喜欢
		videoNow, err = isVideoLike(userId, v.Video_id)
		if err != nil {
			fmt.Println("未找到该视频", v.Video_id)
		}
		videoList = append(videoList, *videoNow)
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

// 查询用户存在
func isUserExists(username string) uint64 {
	user := new(models.User)
	mysql.DB.First(user, "name = ?", username)
	return user.Id
}

// 查询video是否被当前user点赞并设置isFavorite
func isVideoLike(userId int64, videoId int64) (*models.Video, error) {
	//首先搜索videos表 获取到需要的视频信息
	videoNow := new(models.Video)
	result := mysql.DB.First(videoNow, "id = ?", videoId)
	fmt.Println("isVideoLike videoNow: ", videoNow)
	if result.RowsAffected == 0 {
		//在videos表中通过视频id未找到视频
		return nil, result.Error
	}
	//查找点赞关系 并设置是否点赞
	result = mysql.DB.First(&models.User_favorite_video{}, "user_id = ? AND video_id = ?", userId, videoId)
	if result.RowsAffected != 0 {
		//找到
		videoNow.IsFavorite = true
	} else {
		//未找到
		videoNow.IsFavorite = false
	}
	return videoNow, nil
}
