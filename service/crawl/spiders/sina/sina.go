package sina

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"outback/kingo/consts"
	"strconv"
	"strings"
	"time"

	"outback/kingo/config"
	"outback/kingo/dao"
	"outback/kingo/model"
	"outback/kingo/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"

	"github.com/gocolly/colly"
)

// type Document struct {
// 	*Selection
// 	Url      *url.URL
// 	rootNode *html.Node
// }
// type Selection struct {
// 	Nodes    []*html.Node
// 	document *Document
// 	prevSel  *Selection
// }

type StarkSpider struct {
	create dao.CreateDal
	search dao.SearchDal
}

func NewStarkSpider(cre dao.CreateDal, sea dao.SearchDal) *StarkSpider {

	return &StarkSpider{create: cre, search: sea}
}

func (s *StarkSpider) Start(ctx context.Context) {

	// 声明初始化NewCollector对象时可以指定Agent，连接递归深度，URL过滤以及domain限制等
	c := colly.NewCollector(
		// colly.AllowedDomains("news.baidu.com"),
		colly.UserAgent("Opera/9.80 (Windows NT 6.1; U; zh-cn) Presto/2.9.168 Version/11.50"),
		colly.MaxDepth(-1),
	)

	// 发出请求时附的回调
	c.OnRequest(func(r *colly.Request) {
		// Request头部设定
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", "")
		r.Headers.Set("Referer", "http://vip.stock.finance.sina.com.cn/")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")
		r.Headers.Set("authority", "money.finance.sina.com.cn")
		r.Headers.Set("authority", "money.finance.sina.com.cn")
		r.Headers.Set("accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
		r.Headers.Set("Cookie", "name=sinaAds; post=massage; page=23333; NowDate=Sat Oct 30 2021 11:42:10 GMT+0800 (ä¸­å›½æ ‡å‡†æ—¶é—´); UOR=www.google.com,finance.sina.com.cn,; SINAGLOBAL=101.206.250.69_1606034230.695058; U_TRS1=00000000.5f3ccf7.5fba2341.a0a53296; kke_CnLv1_PPT_v2=know; UM_distinctid=179ea4edb87329-0f1204ed8e9d55-1f3b6254-13c680-179ea4edb88403; __gads=ID=eb2fd0309922d2cf-2262930549c90069:T=1623133708:RT=1623133708:S=ALNI_MZAaX1lyVKcp3US2kTz_5qbQ6cJ_g; SR_SEL=1_511; Apache=182.150.57.253_1635498661.533198; ULV=1635498802120:12:12:4:182.150.57.253_1635498661.533198:1635498661440; _s_upa=3; U_TRS2=000000fd.9066544.617bbf44.a46a46b3; FIN_ALL_VISITED=sh603155%2Csz002932%2Csz300459%2Csz002171%2Csz002756%2Csz002240; rotatecount=2; FINA_V_S_2=sh603155,sz002932,sz300459,sz002171,sz002756,sz002240; display=hidden; sinaH5EtagStatus=y")
	})

	c.OnResponse(func(resp *colly.Response) {
		url := resp.Request.URL.String()
		log.Info().Str("url", url).Msg("get response url")
		if resp.StatusCode != http.StatusOK {
			log.Info().Msgf("Status is not OK:%s", url)
			return
		}

		if strings.Contains(url, "hq.sinajs.cn/list") {
			fmt.Println("解析股价")
			// s.ParseStarkPrice(bytes.NewReader(resp.Body))
		} else {
			// goquery直接读取resp.Body的内容
			htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
			if err != nil {
				log.Error().Err(err).Msg("error")
				return
			}

			if strings.Contains(url, "vFD_ProfitStatement") {
				htmlDoc.Find(`#ProfitStatementNewTable0`).Each(s.ParseProfile)
			}

			if strings.Contains(url, "vFD_CashFlow") {
				htmlDoc.Find(`#ProfitStatementNewTable0`).Each(s.ParseCash)
			}
			//
			if strings.Contains(url, "vFD_BalanceSheet") {
				htmlDoc.Find(`#BalanceSheetNewTable0`).Each(s.ParseBalance)
			}
		}
	})

	// 对visit的线程数做限制，visit可以同时运行多个
	if err := c.Limit(&colly.LimitRule{
		Delay:       15 * time.Second,
		RandomDelay: 15 * time.Second,
		DomainGlob:  "*",
		Parallelism: 1,
	}); err != nil {
		log.Error().Err(err).Msg("Limit")
	}

	c.OnError(func(response *colly.Response, err error) {
		url := response.Request.URL.String()
		status := response.StatusCode
		if status == 456 {
			log.Info().Msg("IP已被封禁了")
			return
		}

		log.Error().Err(err).Msgf("get :%s", url)
	})

	codes, err := s.search.SearchNameCode(context.Background(), model.SearchNameCodeParam{})
	if err != nil {
		log.Error().Err(err).Msg("find name and code")
		return
	}
	years := config.GetConfig().CrawlConfig.Period
	for i := range codes {
		for j := range years {
			crawl, err := s.search.SearchCrawl(ctx, model.SearchCrawlParam{Code: codes[i].Code, Year: years[j]})
			if err != nil {
				log.Error().Err(err).Msg("SearchCrawl")
				continue
			}
			if len(crawl) >= 3 {
				log.Info().Str("Code", codes[i].Code).Str("year", years[i]).Msg("data has crawled")
				continue
			}

			urls := make([]string, 0)

			proUrl := fmt.Sprintf("https://money.finance.sina.com.cn/corp/go.php/vFD_ProfitStatement/stockid/%s/ctrl/%s/displaytype/4.phtml", codes[i].Code, years[j])
			cashUrl := fmt.Sprintf("https://money.finance.sina.com.cn/corp/go.php/vFD_CashFlow/stockid/%s/ctrl/%s/displaytype/4.phtml", codes[i].Code, years[j])
			balanceUrl := fmt.Sprintf("https://money.finance.sina.com.cn/corp/go.php/vFD_BalanceSheet/stockid/%s/ctrl/%s/displaytype/4.phtml", codes[i].Code, years[j])
			urls = append(urls, proUrl, cashUrl, balanceUrl)
			for _, ur := range urls {
				if err := c.Visit(ur); err != nil {
					log.Error().Err(err).Str("url", ur).Msgf("Visit")
					continue
				}
				log.Info().Str("url", ur).Msgf("Visit start")
			}
		}

		// ti := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
		// if codes[i].CrawlDate < ti.Unix() && len(codes[i].Code) > 0 && string([]rune(codes[i].Code)[:1]) != "3" {
		// 	if codes[i].SHSZ != "" {
		// 		balanceUrl := fmt.Sprintf("https://hq.sinajs.cn/list=%s%s", codes[i].SHSZ, codes[i].Code)
		// 		if err := c.Visit(balanceUrl); err != nil {
		// 			log.Error().Err(err).Msgf("Visit:%s", balanceUrl)
		// 		}
		// 	} else {
		// 		balanceUrl := fmt.Sprintf("https://hq.sinajs.cn/list=sh%s", codes[i].Code)
		// 		if err := c.Visit(balanceUrl); err != nil {
		// 			log.Error().Err(err).Msgf("Visit:%s", balanceUrl)
		// 		}
		//
		// 		balanceUrl = fmt.Sprintf("https://hq.sinajs.cn/list=sz%s", codes[i].Code)
		// 		if err := c.Visit(balanceUrl); err != nil {
		// 			log.Error().Err(err).Msgf("Visit:%s", balanceUrl)
		// 		}
		// 	}
		// }
	}

}

func (s *StarkSpider) ParseProfile(i int, selection *goquery.Selection) {

	res := make([]string, 0)
	name, code := "", ""

	selection.Find(" tr td").Each(
		func(i int, selection *goquery.Selection) {
			t := selection.Text()
			res = append(res, t)
		},
	)

	selection.Find("tr th ").Each(
		func(i int, selection *goquery.Selection) {
			na, co := utils.ParseNameCode(selection)
			name = na
			code = co
		})

	incomes, err := parseProfile(name, code, res)
	if err != nil {
		log.Error().Err(err).Msg("parseProfile")
		return
	}
	for i := range incomes {
		if err := s.create.CreateProfile(context.Background(), incomes[i]); err != nil {
			log.Error().Err(err).Msg("CreateProfile")
		}
	}
	crawl := &model.Crawl{
		Code:      code,
		Year:      utils.GetReportYear(incomes[0].ReportPeriod),
		CrawlType: consts.ReportTypeProfile,
		CrawlAt:   time.Now().UnixMilli(),
	}
	if err := s.create.CreateCrawl(context.Background(), crawl); err != nil {
		log.Error().Err(err).Msg("CreateCrawl")
	}

	log.Info().Msgf("写入利润表成功,Name %s:Code：%s", name, code)
}

func (s *StarkSpider) ParseCash(j int, selection *goquery.Selection) {
	res := make([]string, 0)
	name, code := "", ""

	selection.Find(" tr td ").Each(
		func(i int, selection *goquery.Selection) {
			t := selection.Text()
			res = append(res, t)
		},
	)

	selection.Find("tr th ").Each(
		func(i int, selection *goquery.Selection) {
			na, co := utils.ParseNameCode(selection)
			name = na
			code = co
		})

	cashs, err := parseCashFlow(name, code, res)
	if err != nil {
		log.Error().Err(err).Msg("parseProfile")
		return
	}
	for i := range cashs {
		if err := s.create.CreateCashFlow(context.Background(), cashs[i]); err != nil {
			log.Error().Err(err).Msg("CreateProfile")
		}
	}
	crawl := &model.Crawl{
		Code:      code,
		Year:      utils.GetReportYear(cashs[0].ReportPeriod), // 这里需要进一步解析
		CrawlType: consts.ReportTypeCash,
		CrawlAt:   time.Now().UnixMilli(),
	}
	if err := s.create.CreateCrawl(context.Background(), crawl); err != nil {
		log.Error().Err(err).Msg("CreateCrawl")
	}

	log.Info().Msgf("写入现金表成功,Name:%s Code：%s", name, code)
}

func (s *StarkSpider) ParseBalance(i int, selection *goquery.Selection) {

	res := make([]string, 0)
	name, code := "", ""

	selection.Find(" tr td").Each(
		func(i int, selection *goquery.Selection) {
			t := selection.Text()
			res = append(res, t)
		},
	)

	selection.Find("tr th ").Each(
		func(i int, selection *goquery.Selection) {
			na, co := utils.ParseNameCode(selection)
			name = na
			code = co
		})

	balance, err := parseBalance(name, code, res)
	if err != nil {
		log.Error().Err(err).Msg("parseProfile")
		return
	}
	for i := range balance {
		if err := s.create.CreateBalance(context.Background(), balance[i]); err != nil {
			log.Error().Err(err).Msg("CreateProfile")
		}
	}
	crawl := &model.Crawl{
		Code:      code,
		Year:      utils.GetReportYear(balance[0].ReportPeriod),
		CrawlType: consts.ReportTypeProfile,
		CrawlAt:   time.Now().UnixMilli(),
	}
	if err := s.create.CreateCrawl(context.Background(), crawl); err != nil {
		log.Error().Err(err).Msg("CreateCrawl")
	}

	log.Info().Msgf("写入资产负债表,Name: %s Code：%s", name, code)
}

func parseProfile(name, code string, res []string) ([]*model.Profile, error) {
	per := utils.ParsePeriodCnt(res)
	date := utils.ReportDate(res)
	if per != len(date) {
		log.Error().Err(fmt.Errorf("日期列不符:%s,%s", name, code))
		return nil, fmt.Errorf("日期列不符")
	}
	ans := make([]*model.Profile, per)
	for i := range ans {
		ans[i] = &model.Profile{}
	}
	for i := 0; i < len(date); i++ {
		ans[i].ReportPeriod = date[i]
		ans[i].Code = code
		ans[i].Name = name
	}

	item := map[string]string{
		"一、营业总收入":     "OperateIn",
		"二、营业总成本":     "OperateAllCost",
		"营业成本":        "OperateCost",
		"营业税金及附加":     "Tax",
		"销售费用":        "SalesCost",
		"管理费用":        "ManageCost",
		"财务费用":        "FinancialCost",
		"研发费用":        "RDCost",
		"五、净利润":       "NetProfit",
		"稀释每股收益(元/股)": "EarnPerShare",
		"投资收益":        "Invest",
		"公允价值变动收益":    "FairIn",
	}
	i := 0
	for i < len(res) {
		fi, ok := item[res[i]]
		if !ok {
			i++
			continue
		}
		for j := 0; j < per; j++ {
			_ = utils.SetField(ans[j], fi, ParseNum(res[i+j+1]))
		}
		i += per + 1
	}

	return ans, nil
}

func ParseNum(res string) int64 {
	if res == "--" {
		return 0
	}
	if strings.Contains(res, "亿") {
		fmt.Println("has 亿")
	}
	if strings.Contains(res, "万") {
		fmt.Println("has 万")
	}
	// 注意可能会有负数哦
	if parseInt, err := strconv.ParseFloat(strings.ReplaceAll(res, ",", ""), 64); err == nil {
		return int64(parseInt)
	}
	return 0
}

func parseCashFlow(name, code string, res []string) ([]*model.CashFlow, error) {
	per := utils.ParsePeriodCnt(res)
	date := utils.ReportDate(res)
	if per != len(date) {
		log.Error().Err(fmt.Errorf("日期列不符:%s,%s", name, code))
		return nil, fmt.Errorf("日期列不符")
	}
	ans := make([]*model.CashFlow, per)
	for i := range ans {
		ans[i] = &model.CashFlow{}
	}
	for i := 0; i < len(date); i++ {
		ans[i].ReportPeriod = date[i]
		ans[i].Code = code
		ans[i].Name = name
	}

	item := map[string]string{
		"销售商品流入":   "SaleIn",
		"税费返还":     "TaxIn",
		"经营活动流入小计": "SumIn",
		"购买商品的流出":  "SaleOut",
		"支付给员工的流出": "EmpOut",
		"流出小计":     "SumOut",
		"经营活动现金净额": "Netflow",
	}

	i := 0
	for i < len(res) {
		fi, ok := item[res[i]]
		if !ok {
			i++
			continue
		}
		for j := 0; j < per; j++ {
			_ = utils.SetField(ans[j], fi, ParseNum(res[i+j+1]))
		}
		i += per + 1
	}

	return ans, nil
}

func parseBalance(name, code string, res []string) ([]*model.Balance, error) {

	per := utils.ParsePeriodCnt(res)
	date := utils.ReportDate(res)
	if per != len(date) {
		log.Error().Err(fmt.Errorf("日期列不符:%s,%s", name, code))
		return nil, fmt.Errorf("日期列不符")
	}
	ans := make([]*model.Balance, per)
	for i := range ans {
		ans[i] = &model.Balance{}
	}
	for i := 0; i < len(date); i++ {
		ans[i].ReportPeriod = date[i]
		ans[i].Code = code
		ans[i].Name = name
	}

	item := map[string]string{
		"货币资金":    "MoneyFunds",
		"交易性金融资产": "TransFinance",
		"应收账款":    "AccountReceive",
		"应收票据":    "NoteReceive",
		"应付账款":    "AccountPay",
		"应付票据":    "NotePay",
		"固定资产":    "Assets",
		"存货":      "Stock",
		"在建工程":    "Construct",
		"短期借款":    "ShortLoan",
		"长期借款":    "LongLoan",
		"实收资本":    "Capital",
	}

	i := 0
	for i < len(res) {
		fi, ok := item[res[i]]
		if !ok {
			i++
			continue
		}
		for j := 0; j < per; j++ {
			_ = utils.SetField(ans[j], fi, ParseNum(res[i+j+1]))
		}
		i += per + 1
	}

	return ans, nil
}
