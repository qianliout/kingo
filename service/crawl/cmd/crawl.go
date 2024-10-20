package cmd

import (
	"context"
	"fmt"
	"os"

	"outback/kingo/dao"
	"outback/kingo/service/crawl/spiders"
	"outback/kingo/service/crawl/spiders/profile"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"outback/kingo/config"

	"github.com/spf13/cobra"
)

func NewCrawlCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crawl",
		Short: "crawl",
		Long:  "\n crawl data",
		Run: func(cmd *cobra.Command, args []string) {
			option := GetOptionByViper()
			c, err := config.ParseConfig(option.ConfigFile)
			if err != nil {
				os.Exit(500)
				return
			}
			if err := ValidateConfig(c); err != nil {
				os.Exit(500)
				return
			}

			ss, err := NewSpider(c)
			if err != nil {
				os.Exit(500)
			}
			for i := range ss {
				ss[i].Start(context.Background())
			}
		},
	}

	return cmd
}

func NewSpider(cfg *config.Config) ([]spiders.Spider, error) {
	// 测试数据库是否能链接
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Database.Username,
		cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	dal := dao.NewCreateDao(db)
	sea := dao.NewSearchDao(db)

	// nameS := names.NewNameCode(dal)
	pro := profile.NewStarkSpider(dal, sea)
	res := make([]spiders.Spider, 0)
	// res = append(res, nameS)
	res = append(res, pro)
	return res, nil
}
