package dongfang

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"outback/kingo/dao"
	"outback/kingo/model"

	"github.com/gocolly/colly"
	"github.com/rs/zerolog/log"
)

// 爬取股票代码
type NameCode struct {
	ErrorURL []string
	PageUrl  string
	create   dao.CreateDal
	Rex      *regexp.Regexp
}

func NewNameCode(cre dao.CreateDal) *NameCode {
	// http://quote.eastmoney.com/center/gridlist.html#hs_a_board
	return &NameCode{
		create:   cre,
		ErrorURL: make([]string, 0),
		PageUrl:  "http://51.push2.eastmoney.com/api/qt/clist/get?cb=jQuery1124044325478855346645_1730194169959&pn=%d&pz=20&po=1&np=1&ut=bd1d9ddb04089700cf9c27f6f7426281&fltt=2&invt=2&dect=1&wbp2u=|0|0|0|web&fid=f3&fs=m:0+t:6,m:0+t:80,m:1+t:2,m:1+t:23,m:0+t:81+s:2048&fields=f12,f14",
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
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Referer", "http://quote.eastmoney.com/center/gridlist.html")
		r.Headers.Set("Cookie", "qgqp_b_id=d2453831d633cf249d9cc79560cc326b; st_si=54291757672398; st_sn=6; st_psi=20241029173647629-113200301321-3262121969; st_asi=delete; st_pvi=55225654094531; st_sp=2024-10-22%2022%3A20%3A24; st_inirUrl=https%3A%2F%2Fwww.google.com%2F")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")
		r.Headers.Set("Referer", "http://www.sse.com.cn/")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:131.0) Gecko/20100101 Firefox/131.0")
		r.Headers.Set("Host", "51.push2.eastmoney.com")
	})

	c.OnResponse(func(resp *colly.Response) {
		log.Info().Str("url", resp.Request.URL.String()).Msg("response received")
		if resp.StatusCode != 200 {
			log.Info().Str("url", resp.Request.URL.String()).Int("status", resp.StatusCode).Msg("response received but not get data")
			return
		}

		body := string(resp.Body)
		start := strings.Index(body, "(")
		end := strings.Index(body, ")")
		body = body[start+1 : end]

		res := new(Response)
		err := json.Unmarshal([]byte(body), res)
		if err != nil {
			log.Error().Err(err).Msg("error")
			return
		}
		for _, ch := range res.Data.Diff {
			data := &model.NameCode{Name: ch.Name, Code: ch.Code}
			if err := s.create.CreateNameCode(context.Background(), data); err != nil {
				log.Error().Err(err).Str("data", data.LogStr()).Msg("create name")
			}
			log.Info().Interface("data", data).Msg("create name")
		}
	})

	// 对visit的线程数做限制，visit可以同时运行多个
	_ = c.Limit(&colly.LimitRule{
		Parallelism: 1,
		RandomDelay: 1 * time.Second,
		DomainGlob:  "*",
		Delay:       1 * time.Second,
	})

	c.OnError(func(response *colly.Response, err error) {
		log.Error().Err(err).Msg(response.Ctx.Get("url"))
	})
	// 上证
	for page := 1; page <= 284; page++ {
		url := fmt.Sprintf(s.PageUrl, page)
		err := c.Visit(url)
		if err != nil {
			log.Err(err).Str("url", url).Msg("start visit")
			continue
		}
		log.Info().Str("url", url).Msg("start visit")
	}
}

type Response struct {
	Data struct {
		Total int `json:"total"`
		Diff  []struct {
			Code string `json:"f12"`
			Name string `json:"f14"`
		} `json:"diff"`
	} `json:"data"`
}
