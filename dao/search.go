package dao

import (
	"context"
	"time"

	"outback/kingo/model"

	"gorm.io/gorm"
)

type SearchDal interface {
	SearchNameCode(ctx context.Context, param model.SearchNameCodeParam) ([]model.NameCode, error)
	SearchCrawl(ctx context.Context, param model.SearchCrawlParam) ([]model.Crawl, error)
}

type SearchDao struct {
	db *gorm.DB
}

func (dal *SearchDao) SearchNameCode(ctx context.Context, param model.SearchNameCodeParam) ([]model.NameCode, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()

	res := make([]model.NameCode, 0)
	db := dal.db.WithContext(ctx).Table(new(model.NameCode).TableName())
	if param.Name != "" {
		db = db.Where("name like ?", "%"+param.Name+"%")
	}
	if param.Code != "" {
		db = db.Where("code = ?", param.Code)
	}

	err := db.Find(&res).Error
	return res, err
}

func (dal *SearchDao) SearchCrawl(ctx context.Context, param model.SearchCrawlParam) ([]model.Crawl, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()

	res := make([]model.Crawl, 0)
	db := dal.db.WithContext(ctx).Table(new(model.Crawl).TableName())
	if param.Code != "" {
		db = db.Where("code = ?", param.Code)
	}
	if param.Year != "" {
		db = db.Where("report_period = ?", param.Code)
	}
	err := db.Find(&res).Error
	return res, err
}

func NewSearchDao(db *gorm.DB) *SearchDao {
	return &SearchDao{db: db}
}
