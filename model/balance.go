package model

import (
	"fmt"

	"outback/kingo/utils"
)

// 资产负债表

type Balance struct {
	ID             int64  `gorm:"column:id"`
	UniqueID       int64  `gorm:"column:unique_id"`
	Name           string `gorm:"column:name"`          // 名字
	Code           string `gorm:"column:code"`          // 代码
	ReportPeriod   string `gorm:"column:report_period"` // 报告期
	Year           int64  `gorm:"column:year"`
	Month          int64  `gorm:"column:month"`
	MoneyFunds     int64  `gorm:"column:money_funds"`     // 现金资产
	TransFinance   int64  `gorm:"column:trans_finance"`   // 交易性金融资产
	AccountReceive int64  `gorm:"column:account_receive"` // 应收账款
	NoteReceive    int64  `gorm:"column:note_receive"`    // 应收票据
	AccountPay     int64  `gorm:"column:account_pay"`     // 应付账款
	NotePay        int64  `gorm:"column:note_pay"`        // 应付票据
	Assets         int64  `gorm:"column:assets"`          // 固定资产
	Stock          int64  `gorm:"column:stock"`           // 存货
	Construct      int64  `gorm:"column:construct"`       // 在建工程
	ShortLoan      int64  `gorm:"column:short_loan"`      // 短期借款
	LongLoan       int64  `gorm:"column:long_loan"`       // 长期借款
	Capital        int64  `gorm:"column:capital"`         // 实收资本

	CreatedAt int64 `gorm:"autoCreateTime:milli;column:created_at"` // milliseconds
	UpdatedAt int64 `gorm:"autoUpdateTime:milli;column:updated_at"` // milliseconds
}

func (vi *Balance) TableName() string {
	return "balances"
}

func (vi *Balance) Serialize() {
	vi.UniqueID = utils.GenerateUUID64(fmt.Sprintf("%s-%s", vi.Code, vi.ReportPeriod))
}
