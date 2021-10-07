package common

import (
	"fmt"

	"github.com/micro/go-micro/v2/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

type MysqlConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Pwd      string `json:"pwd"`
	Database string `json:"database"`
	Port     int64  `json:"port"`
}

// 获取mysql配置
func GetMysqlFromConfig(config config.Config, path ...string) *MysqlConfig {
	mc := &MysqlConfig{}

	config.Get(path...).Scan(mc)

	return mc
}

func InitMysql(mc *MysqlConfig) error {

	dsn := mc.User + ":" + mc.Pwd + "@tcp(" + mc.Host + ":" + fmt.Sprintf("%d", mc.Port) + ")/" + mc.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	var err error

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
		// 设置gorm日志级别
		Logger: logger.Default.LogMode(logger.Info),
	})

	return err
}

func GetMysql() *gorm.DB {
	return db
}
