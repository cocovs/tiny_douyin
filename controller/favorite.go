package controller

//收藏
import (
	"github.com/cocovs/tiny-douyin/dao/mysql"
	"github.com/cocovs/tiny-douyin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction 点赞操作
func FavoriteAction(c *gin.Context) {

	userId, err := strconv.ParseInt((c.PostForm("user_id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "1字符转换错误"})
	}
	videoId, err := strconv.ParseInt(c.PostForm("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "2字符转换错误"})
	}

	actionType := c.PostForm("action_type")

	//创建点赞结构体
	userLike := models.UserLike{
		User_id:  userId,
		Video_id: videoId,
	}
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

// FavoriteList 点赞列表
func FavoriteList(c *gin.Context) {
	userId, err := strconv.ParseInt((c.PostForm("user_id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "点赞列表字符串转换错误"})
	}
	//username, exist := c.Get("username")
	//if !exist {
	//	//未获取到该用户名
	//	c.JSON(http.StatusOK, models.Response{StatusCode: 1, StatusMsg: "未获取到该用户名"})
	//}

	//从`user_favorite_videos`中搜索该用户id 返回所有的视频id切片
	videIdList := make([]models.Video, 0, 0)
	result := mysql.DB.Find(videIdList, "user_id = ?", userId)

	//点赞列表
	favoriteList := make([]models.Video, 0, 0)

	// 填充favoriteList
	for i := int64(0); i < result.RowsAffected; i++ {
		if isVideoLike(video[i], userId) {
			videoList = append(videoList, video[i])
		}
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: models.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

// 查询用户存在
func isUserExists(username string) int64 {
	user := new(models.User)

	mysql.DB.First(user, "name = ?", username)

	return user.Id
}

// 查询video是否被当前user点赞并设置isFavorite
func isVideoLike(video models.Video, userId int64) bool {
	videoId := strconv.FormatInt(video.Id, 10)
	key := "video:liked:" + videoId
	isMember := RDB.SIsMember(ctx, key, userId)
	mysql.DB.Where("videoId = ?", videoId).Update("is_favorite", isMember.Val())
	return isMember.Val()
}
