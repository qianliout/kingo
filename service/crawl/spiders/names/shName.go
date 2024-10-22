package names

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"outback/kingo/dao"
	"outback/kingo/model"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

// 爬取上证的股票代码

type NameCode struct {
	ErrorURL []string
	PageUrl  string
	create   dao.CreateDal
}

func NewNameCode(cre dao.CreateDal) *NameCode {
	return &NameCode{
		create:   cre,
		ErrorURL: make([]string, 0),
		PageUrl: "http://query.sse.com.cn/security/stock/getStockListData2.do?&isPagination=true&" +
			"stockCode=&csrcCode=&areaName=&stockType=1&pageHelp.cacheSize=1&pageHelp.beginPage=%d&" +
			"pageHelp.pageSize=25&pageHelp.pageNo=%d&pageHelp.endPage=651&_=%d",
	}
}

func (s *NameCode) Start(ctx context.Context) {
	// NewCollector(options ...func(*Collector)) *Collector
	// 声明初始化NewCollector对象时可以指定Agent，连接递归深度，URL过滤以及domain限制等
	c := colly.NewCollector(
		colly.UserAgent("Opera/9.80 (Windows NT 6.1; U; zh-cn) Presto/2.9.168 Version/11.50"),
		// colly.AllowedDomains("sina.com.cn"),
		colly.MaxDepth(-1),
	)

	// 发出请求时附的回调
	c.OnRequest(func(r *colly.Request) {
		// Request头部设定
		// r.Headers.Set("Host", "baidu.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", "")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")
		r.Headers.Set("Referer", "http://www.sse.com.cn/")
		r.Headers.Set("accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
		log.Error().Str("url", r.URL.String()).Msg("start crawl")
	})

	c.OnResponse(func(resp *colly.Response) {
		log.Info().Str("url", resp.Request.URL.String()).Msg("response received")
		if resp.StatusCode != 200 {
			log.Info().Str("url", resp.Request.URL.String()).Int("status", resp.StatusCode).Msg("response received but not get data")
			return
		}

		res := new(model.NubSh)
		err := json.Unmarshal(resp.Body, res)
		if err != nil {
			log.Error().Err(err).Msg("error")
			return
		}
		for i := range res.Result {
			data := &model.NameCode{Name: res.Result[i].Name, Code: res.Result[i].Code}
			if err := s.create.CreateNameCode(context.Background(), data); err != nil {
				log.Error().Err(err).Str("data", data.LogStr()).Msg("create name")
			}
			log.Info().Interface("data", data).Msg("create name")
		}
	})

	// 对visit的线程数做限制，visit可以同时运行多个
	_ = c.Limit(&colly.LimitRule{
		Parallelism: 1,
		RandomDelay: 15 * time.Second,
		DomainGlob:  "*",
		Delay:       15 * time.Second,
	})

	c.OnError(func(response *colly.Response, err error) {
		log.Error().Err(err).Msg(response.Ctx.Get("url"))
	})
	// 上证
	for page := 1; page <= 5; page++ {
		url := fmt.Sprintf(s.PageUrl, page, page, time.Now().UnixMilli())
		err := c.Visit(url)
		if err != nil {
			log.Err(err).Str("url", url).Msg("start visit")
			continue
		}
		log.Info().Str("url", url).Msg("start visit")
	}
}
