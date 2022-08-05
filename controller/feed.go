package controller

//视频流
import (
	"github.com/cocovs/tiny-douyin/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	models.Response
	//视频列表
	VideoList []models.Video `json:"video_list,omitempty"`
	//本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	NextTime int64 `json:"next_time,omitempty"`
}

// FavoriteRelation 点赞关系在user_favorite_video 表中查找用户点赞视频的关系
func FavoriteRelation() {

}

//投稿时间倒叙，从最近的开始30个视频
//该视频列中发布最早的视频时间 作为下次视频列表的最近时间
//获取视频 对于是否点赞 应该在返回时搜索  user_favorite_video表 中该用户是否点赞
// 视频流
func Feed(c *gin.Context) {
	latest_time, _ := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	//token := c.Query("token")
	videosList, next_time := NewVideoList(latest_time) //获取视频列表函数
	//获取发布最早时间

	c.JSON(http.StatusOK, FeedResponse{
		Response:  models.Response{StatusCode: 0}, //状态码，0-成功，其他值-失败
		VideoList: videosList,                     //视频列表
		//本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
		NextTime: next_time,
	})
}
