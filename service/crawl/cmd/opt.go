package cmd

import (
	"fmt"
	"strings"

	"outback/kingo/config"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	ConfigFile = "config-file"
)

type Option struct {
	ConfigFile string
}

func GetOptionByViper() Option {
	return Option{
		ConfigFile: viper.GetString(ConfigFile),
	}
}

func AddOption(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringP(ConfigFile, "c", "", "config file")

	rootCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		err := viper.BindPFlag(flag.Name, flag)
		if err != nil {
			panic(err)
		}
	})

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func ValidateConfig(cfg *config.Config) error {
	// 测试数据库是否能链接
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/dbname?charset=utf8mb4&parseTime=True&loc=Local", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port)
	_, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}
