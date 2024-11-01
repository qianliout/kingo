package model

import (
	"fmt"
	"outback/kingo/utils"
)

// 股票列表

type NameCode struct {
	ID       int64  `gorm:"column:id"`
	Name     string `gorm:"column:name"`     // 名字
	Code     string `gorm:"column:code"`     // 代码
	Price    int64  `gorm:"column:price"`    // 价格的100倍
	CrawlAt  int64  `gorm:"column:crawl_at"` // 爬取时间
	Plate    string `gorm:"column:plate"`    // 板块：主版，三版，创业版
	Industry string `gorm:"column:industry"` // 行业
	Exchange string `gorm:"column:exchange"` // 交易所

	CreatedAt int64 `gorm:"autoCreateTime:milli;column:created_at"` // milliseconds
	UpdatedAt int64 `gorm:"autoUpdateTime:milli;column:updated_at"` // milliseconds
}

func (vi *NameCode) TableName() string {
	return "names"
}

func (vi *NameCode) Serialize() {
}

func (vi *NameCode) LogStr() string {
	str := fmt.Sprintf("code:%s,name:%s", vi.Code, vi.Name)
	return str
}

func (vi *NameCode) Check() error {
	if vi.Name == "" {
		return fmt.Errorf("name is empty")
	}
	if vi.Code == "" {
		return fmt.Errorf("code is empty")
	}
	return nil
}

type Industry struct {
	Name string
	Code string
}

type NubSh struct {
	Result []DataSH `json:"result"`
}

type DataSH struct {
	Code string `json:"SECURITY_CODE_A"`
	Name string `json:"COMPANY_ABBR"`
}

type Crawl struct {
	ID        int64  `gorm:"column:id"`
	UniqueID  int64  `gorm:"column:unique_id"`
	Code      string `gorm:"column:code"` // 代码
	Year      string `gorm:"column:year"` // 报告年份
	CrawlType string `gorm:"column:crawl_type"`
	CrawlAt   int64  `gorm:"column:crawl_at"`
	CreatedAt int64  `gorm:"autoCreateTime:milli;column:created_at"` // milliseconds
	UpdatedAt int64  `gorm:"autoUpdateTime:milli;column:updated_at"` // milliseconds
}

func (vi *Crawl) TableName() string {
	return "crawl"
}

func (vi *Crawl) Serialize() {
	vi.UniqueID = utils.GenerateUUID64(fmt.Sprintf("%s-%s-%s", vi.Code, vi.Year, vi.CrawlType))
}
