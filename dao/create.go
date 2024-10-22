package dao

import (
	"context"
	"time"

	"outback/kingo/model"

	"gorm.io/gorm"
)

type CreateDal interface {
	CreateProfile(ctx context.Context, data *model.Profile) error
	CreateBalance(ctx context.Context, data *model.Balance) error
	CreateCashFlow(ctx context.Context, data *model.CashFlow) error
	CreateNameCode(ctx context.Context, data *model.NameCode) error
	CreateCrawl(ctx context.Context, data *model.Crawl) error
}

type CreateDao struct {
	db *gorm.DB
}

func (dal *CreateDao) CreateBalance(ctx context.Context, data *model.Balance) error {
	data.Serialize()

	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()
	db := dal.db.WithContext(ctx).Table(data.TableName())
	if err := db.Where("unique_id = ?", data.UniqueID).Delete(&model.Balance{}).Error; err != nil {
		return err
	}

	err := db.Create(data).Error
	return err
}

func (dal *CreateDao) CreateCashFlow(ctx context.Context, data *model.CashFlow) error {
	data.Serialize()

	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()

	db := dal.db.WithContext(ctx).Table(data.TableName())

	if err := db.Where("unique_id = ?", data.UniqueID).Delete(&model.CashFlow{}).Error; err != nil {
		return err
	}

	err := db.Create(data).Error
	return err
}

func (dal *CreateDao) CreateProfile(ctx context.Context, data *model.Profile) error {
	data.Serialize()
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()
	db := dal.db.Table(data.TableName()).WithContext(ctx)
	if err := db.Where("unique_id = ?", data.UniqueID).
		Delete(&model.Profile{}).Error; err != nil {
		return err
	}

	return db.Create(data).Error
}

func (dal *CreateDao) CreateNameCode(ctx context.Context, data *model.NameCode) error {
	data.Serialize()
	if err := data.Check(); err != nil {
		return nil
	}
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()

	db := dal.db.Table(data.TableName()).WithContext(ctx)
	res := make([]*model.NameCode, 0)
	if err := db.Where("code = ?", data.Code).Find(&res).Error; err != nil {
		return err
	}
	if len(res) > 0 {
		return nil
	}
	err := db.Create(data).Error
	return err
}

func (dal *CreateDao) CreateCrawl(ctx context.Context, data *model.Crawl) error {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()
	data.Serialize()
	db := dal.db.Table(data.TableName()).WithContext(ctx)
	res := make([]*model.NameCode, 0)
	if err := db.Where("unique_id = ?", data.UniqueID).Find(&res).Error; err != nil {
		return err
	}
	if len(res) > 0 {
		up := map[string]interface{}{"crawl_at": data.CrawlAt}
		return db.Where("id = ?", res[0].ID).Updates(up).Error
	}
	err := db.Create(data).Error
	return err
}

func NewCreateDao(db *gorm.DB) *CreateDao {
	return &CreateDao{db: db}
}
