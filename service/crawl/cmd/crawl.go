package cmd

import (
	"context"
	"fmt"
	"os"

	"outback/kingo/config"
	"outback/kingo/dao"
	"outback/kingo/service/crawl/spiders"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

			ss := NewSpider(c)
			for i := range ss {
				ss[i].Start(context.Background())
			}
		},
	}

	return cmd
}

func NewSpider(cfg *config.Config) []spiders.Spider {
	// 测试数据库是否能链接
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/dbname?charset=utf8mb4&parseTime=True&loc=Local", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port)
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	dal := dao.NewCreateDao(db)

	nameS := spiders.NewNameCode(dal)
	res := make([]spiders.Spider, 0)
	res = append(res, nameS)
	return res
}
