package cmd

import (
	"context"
	"fmt"
	"os"
	"outback/kingo/service/crawl/spiders/names"
	"outback/kingo/service/flag"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"outback/kingo/dao"
	"outback/kingo/service/crawl/spiders"

	"outback/kingo/config"

	"github.com/spf13/cobra"
)

func NewCrawlNameCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "name",
		Short: "crawl name",
		Long:  "\n crawl name data",
		Run: func(cmd *cobra.Command, args []string) {
			option := flag.GetOptionByViper()
			c, err := config.ParseConfig(option.ConfigFile)
			if err != nil {
				os.Exit(500)
				return
			}
			if err := flag.ValidateConfig(c); err != nil {
				os.Exit(500)
				return
			}

			ss, err := NewNameSpider(c)
			if err != nil {
				os.Exit(500)
			}

			ss.Start(context.Background())
		},
	}

	return cmd
}

func NewNameSpider(cfg *config.Config) (spiders.Spider, error) {
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

	nameS := names.NewNameCode(dal)
	return nameS, nil
}
