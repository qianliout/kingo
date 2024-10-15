package cmd

import (
	"fmt"
	"outback/kingo/config"
	"outback/kingo/service/flag"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Option struct {
	ConfigFile string
}

func GetOptionByViper() Option {
	return Option{
		ConfigFile: viper.GetString(flag.ConfigFile),
	}
}

func ValidateConfig(cfg *config.Config) error {
	// 测试数据库是否能链接
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
