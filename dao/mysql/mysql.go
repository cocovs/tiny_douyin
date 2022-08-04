package mysql

import (
	"fmt"
	"github.com/cocovs/tiny-douyin/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全局变量 打开数据库DB
var (
	DB *gorm.DB
)

func InitMySQL() {
	//连接数据库(加载文件config中的配置)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.MyCnf.User, setting.MyCnf.Password, setting.MyCnf.Host, setting.MyCnf.Port, setting.MyCnf.DB)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//DB.AutoMigrate(&models.User{})
	//DB.AutoMigrate(&models.Video{})
	if err != nil {
		fmt.Println("数据库连接错误： ", err)
	}
}
