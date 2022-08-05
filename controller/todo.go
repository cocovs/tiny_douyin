package controller

//数据操作
import (
	"context"
	"fmt"
	"github.com/cocovs/tiny-douyin/dao/mysql"
	"github.com/cocovs/tiny-douyin/models"
	"github.com/cocovs/tiny-douyin/setting"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//鉴权 搜索用户是否在users表中
//true:存在用户 false：无用户
//func QueryUsersToken(token string) (models.User, bool) {
//	var userNow models.User
//	//DB.Where("Token = ?", token).First(&userNew)
//	mysql.DB.Where("token = ?", token).First(&userNow)
//	if userNow.Token != "" {
//		//不为空，已有用户
//		return userNow, true
//	} else {
//		//空 没有用户
//		return userNow, false
//	}
//
//}

// NewUser 上传用户信息到users表
func NewUser(username string, password string) (newUser models.User) {
	//fmt.Println("这里是新增用户")
	newUser = models.User{
		//id 自增长 之后添加雪花算法生成id
		Name:          username,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		//Token:         token,
		Password: password,
	}

	mysql.DB.Create(&newUser) //传入该结构体
	return newUser
}

// NewVideo 上传视频信息到videos表
func NewVideo(dataName string, userID int64) /*(newVideo Video)*/ {
	//需要传入用户id 和 视频访问url
	newVideo := models.Video{
		AuthorId: userID,
		PlayUrl:  setting.MyOss.MyOSSUrl + "/douyin_video/" + dataName,
		//CoverUrl: ,
		ReleaseTime: time.Now().Unix(),
	}

	mysql.DB.Create(&newVideo)
	//return newVideo
}

// NewVideoList feed流视频列表
func NewVideoList(latest_time int64) ([]models.Video, int64) {

	//1.1返回小于最近时间的最近的三个视频
	var videosList []models.Video
	var count int
	//小于最近时间的最近的三个视频放入 videosList
	mysql.DB.Limit(3).Where("release_time < ?", latest_time).Order("id desc").Find(&videosList)

	//2.1根据UserId从数据库中取出用户信息User结构体
	count = len(videosList)

	var timeNow int64
	//遍历视频列表 为每一个视频结构体加上用户信息
	for index, value := range videosList {
		userid := value.AuthorId
		//从数据库中找到用户id
		var usernow models.User
		mysql.DB.Where("id = ?", userid).First(&usernow)
		//分别为每个赋予用户结构体
		value.Author = usernow
		//查询是否点赞 并为该视频赋值 通过搜索
		userLike := new(models.User_favorite_video)
		result := mysql.DB.First(userLike, "user_id = ? AND video_id = ?", userid, value.Id)
		fmt.Println(*userLike)
		if result.RowsAffected == 1 {
			//找到
			value.IsFavorite = true
		} else {
			value.IsFavorite = false
		}
		//如果该视频为视频库中最后一个视频则返回 现在时间
		if index == count-1 {
			timeNow = value.ReleaseTime
			if value.Id == 1 {
				timeNow = time.Now().Unix()
			}
		}
	}
	fmt.Println(videosList)
	return videosList, timeNow
}

// UserVideoList 获取用户发布视频列表
func UserVideoList(User *models.User) []models.Video {
	//搜索数据库中所有  Video.UserId == User.Id  的
	//该用户信息

	//该用户所有视频
	VideoList := make([]models.Video, 0, 0)
	mysql.DB.Find(&VideoList, "author_id = ?", User.Id)

	//遍历加入用户信息
	for _, value := range VideoList {
		value.Author = *User
	}
	return VideoList
}

// UpOSS 视频文件上传至对象存储桶内
func UpOSS(Myfile io.Reader, FileName string) error {
	//
	u, _ := url.Parse(setting.MyOss.MyOSSUrl) //*url.URL
	b := &cos.BaseURL{BucketURL: u}           //*cos.BaseURL 访问桶列表
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  setting.MyOss.SecId,   // 替换为用户的 SecretId，
			SecretKey: setting.MyOss.SectKey, // 替换为用户的 SecretKey，
			//请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
		},
	})
	//对象的存储位置 文件夹名/文件名
	name := "douyin_video/" + FileName
	_, err := c.Object.Put(context.Background(), name, Myfile, nil)
	return err
}
