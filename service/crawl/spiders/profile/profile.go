package profile

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"outback/kingo/dao"
	"outback/kingo/items"
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

			// if strings.Contains(url, "vFD_CashFlow") {
			// 	htmlDoc.Find(`#ProfitStatementNewTable0`).Each(s.ParseCash)
			// }
			//
			// if strings.Contains(url, "vFD_BalanceSheet") {
			// 	htmlDoc.Find(`#BalanceSheetNewTable0`).Each(s.ParseBalance)
			// }
		}

	})

	// 对visit的线程数做限制，visit可以同时运行多个
	if err := c.Limit(&colly.LimitRule{
		Parallelism: 1,
		Delay:       15 * time.Second,
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

	codes, err := s.search.SearchNameCode(context.Background(), items.SearchNameCodeParam{})
	if err != nil {
		log.Error().Err(err).Msg("find name and code")
		return
	}
	years := []string{"2023", "2021"}
	for i := range codes {
		for j := range years {
			proUrl := fmt.Sprintf("https://money.finance.sina.com.cn/corp/go.php/vFD_ProfitStatement/stockid/%s/ctrl/%s/displaytype/4.phtml", codes[i].Code, years[j])
			if err := c.Visit(proUrl); err != nil {
				log.Error().Err(err).Msgf("Visit：%s", proUrl)
				continue
			}

			// if codes[i].CashFlow == 0 {
			// 	cashUrl := fmt.Sprintf("https://money.finance.sina.com.cn/corp/go.php/vFD_CashFlow/stockid/%s/ctrl/%s/displaytype/4.phtml", codes[i].Code, years[j])
			// 	if err := c.Visit(cashUrl); err != nil {
			// 		log.Error().Err(err).Msgf("Visit: %s", cashUrl)
			// 	}
			// }
			//
			// if codes[i].Balance == 0 {
			// 	balanceUrl := fmt.Sprintf("https://money.finance.sina.com.cn/corp/go.php/vFD_BalanceSheet/stockid/%s/ctrl/%s/displaytype/4.phtml", codes[i].Code, years[j])
			// 	if err := c.Visit(balanceUrl); err != nil {
			// 		log.Error().Err(err).Msgf("Visit:%s", balanceUrl)
			// 	}
			// }
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

	log.Info().Msgf("写入利润表成功,Name %s:Code：%s", name, code)
}

// func (s *StarkSpider) ParseCash(i int, selection *goquery.Selection) {
// 	res := make([]string, 0)
// 	name, code := "", ""
//
// 	selection.Find(" tr td ").Each(
// 		func(i int, selection *goquery.Selection) {
// 			t := selection.Text()
// 			res = append(res, t)
// 		},
// 	)
//
// 	selection.Find("tr th ").Each(
// 		func(i int, selection *goquery.Selection) {
// 			na, co := spiders.parseNameCode(selection)
// 			name = na
// 			code = co
// 		})
//
// 	cashs, err := parseCashFlow(name, code, res)
// 	if err != nil {
// 		log.Error().Err(err).Msg("parseProfile")
// 		return
// 	}
// 	for i := range cashs {
// 		if err := s.create.CreateCodeCashFlow(context.Background(), cashs[i]); err != nil {
// 			log.Error().Err(err).Msg("CreateProfile")
// 		}
//
// 		updater := map[string]interface{}{"cash_flow": time.Now().Unix()}
// 		if err := s.create.UpdateNameCode(context.Background(), code, updater); err != nil {
// 			log.Error().Err(err).Msg("UpdateNameCode")
// 		}
// 	}
// 	log.Info().Msgf("写入现金表成功,Name:%s Code：%s", name, code)
//
// }
//
// func (s *StarkSpider) ParseBalance(i int, selection *goquery.Selection) {
//
// 	res := make([]string, 0)
// 	name, code := "", ""
//
// 	selection.Find(" tr td").Each(
// 		func(i int, selection *goquery.Selection) {
// 			t := selection.Text()
// 			res = append(res, t)
// 		},
// 	)
//
// 	selection.Find("tr th ").Each(
// 		func(i int, selection *goquery.Selection) {
// 			na, co := spiders.parseNameCode(selection)
// 			name = na
// 			code = co
// 		})
//
// 	cashs, err := parseBalance(name, code, res)
// 	if err != nil {
// 		log.Error().Err(err).Msg("parseProfile")
// 		return
// 	}
// 	for i := range cashs {
// 		if err := s.create.CreateBalance(context.Background(), cashs[i]); err != nil {
// 			log.Error().Err(err).Msg("CreateProfile")
// 		}
// 		updater := map[string]interface{}{"balance": time.Now().Unix()}
// 		if err := s.create.UpdateNameCode(context.Background(), code, updater); err != nil {
// 			log.Error().Err(err).Msg("UpdateNameCode")
// 		}
// 	}
// 	log.Info().Msgf("写入资产负债表,Name: %s Code：%s", name, code)
//
// }

func parseProfile(name, code string, res []string) ([]*items.Profile, error) {
	per := utils.ParsePeriod(res)
	date := utils.ReportDate(res)
	if per != len(date) {
		log.Error().Err(fmt.Errorf("日期列不符:%s,%s", name, code))
		return nil, fmt.Errorf("日期列不符")
	}
	ans := make([]*items.Profile, per)
	for i := range ans {
		ans[i] = &items.Profile{}
	}
	for i := 0; i < len(date); i++ {
		ans[i].ReportPeriod = date[i]
		ans[i].Code = code
		ans[i].Name = name
	}
	i := 0
	for i < len(res) {
		switch res[i] {
		case "一、营业总收入":
			for j := 0; j < per; j++ {
				ans[j].OperateIn = ParseNums(res[i+j+1])
			}
			i = i + per
		case "二、营业总成本":
			for j := 0; j < per; j++ {
				ans[j].OperateAllCost = ParseNums(res[i+j+1])
			}
			i = i + per
		case "营业成本":
			for j := 0; j < per; j++ {
				ans[j].OperateCost = ParseNums(res[i+j+1])
			}
			i = i + per
		case "营业税金及附加":
			for j := 0; j < per; j++ {
				ans[j].Tax = ParseNums(res[i+j+1])
			}
			i = i + per
		case "销售费用":
			for j := 0; j < per; j++ {
				ans[j].SalesCost = ParseNums(res[i+j+1])
			}
			i = i + per
		case "管理费用":
			for j := 0; j < per; j++ {
				ans[j].OperateCost = ParseNums(res[i+j+1])
			}
			i = i + per
		case "财务费用":
			for j := 0; j < per; j++ {
				ans[j].FinancialCost = ParseNums(res[i+j+1])
			}
			i = i + per
		case "研发费用":
			for j := 0; j < per; j++ {
				ans[j].RDCost = ParseNums(res[i+j+1])
			}
			i = i + per
		case "五、净利润":
			for j := 0; j < per; j++ {
				ans[j].NetProfit = ParseNums(res[i+j+1])
			}
			i = i + per
		case "稀释每股收益(元/股)":
			for j := 0; j < per; j++ {
				ans[j].EarnPerShare = ParseNums(res[i+j+1])
			}
			i = i + per
		case "投资收益":
			for j := 0; j < per; j++ {
				ans[j].Invest = ParseNums(res[i+j+1])
			}
		case "公允价值变动收益":
			for j := 0; j < per; j++ {
				ans[j].FairIn = ParseNums(res[i+j+1])
			}
		}
		i = i + per
	}

	return ans, nil
}

func ParseNums(res string) int64 {
	if strings.Contains(res, "亿") {
		fmt.Println("has 亿")
	}
	if strings.Contains(res, "万") {
		fmt.Println("has 万")
	}
	if parseInt, err := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(res, ",", ""), "-", ""), 64); err == nil {
		return int64(parseInt)
	}
	return 0
}

// func Parse() {
// 	for j := 0; j < per; j++ {
// 		// parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64)
// 		ss := strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1)
// 		if ss != "" {
// 			if parseInt, err := strconv.ParseFloat(ss, 64); err == nil {
// 				ans[j].OperateIn = int64(parseInt)
// 			}
// 		}
// 	}
// }

//
// func parseCashFlow(name, code string, res []string) ([]items.CashFlow, error) {
// 	per := utils.ParsePeriod(res)
// 	date := utils.ReportDate(res)
// 	if per != len(date) {
// 		log.Error().Err(fmt.Errorf("日期列不符"))
// 		return nil, fmt.Errorf("日期列不符")
// 	}
// 	ans := make([]items.CashFlow, per)
// 	for i := 0; i < len(date); i++ {
// 		ans[i].ReportPeriod = date[i]
// 		ans[i].Code = code
// 		ans[i].Name = name
// 	}
// 	i := 0
// 	for i < len(res) {
// 		switch res[i] {
// 		case "销售商品、提供劳务收到的现金":
// 			for j := 0; j < per; j++ {
// 				parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64)
// 				if err != nil {
// 					log.Error().Err(err).Msg("strconv.ParseInt")
// 				}
// 				ans[j].SalesCash = parseInt
// 			}
// 			i = i + per
// 		case "经营活动现金流入小计":
// 			for j := 0; j < per; j++ {
// 				parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64)
// 				if err != nil {
// 					log.Error().Err(err).Msg("strconv.ParseInt")
// 				}
// 				ans[j].SumInFow = parseInt
// 			}
// 			i = i + per
// 		case "购买商品、接受劳务支付的现金":
// 			for j := 0; j < per; j++ {
// 				parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64)
// 				if err != nil {
// 					log.Error().Err(err).Msg("strconv.ParseInt")
// 				}
// 				ans[j].BuyCash = parseInt
// 			}
// 			i = i + per
// 		case "经营活动现金流出小计":
// 			for j := 0; j < per; j++ {
// 				parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64)
// 				if err != nil {
// 					log.Error().Err(err).Msg("strconv.ParseInt")
// 				}
// 				ans[j].SumOutFow = parseInt
// 			}
// 			i = i + per
// 		case "经营活动产生的现金流量净额":
// 			for j := 0; j < per; j++ {
// 				parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64)
// 				if err != nil {
// 					log.Error().Err(err).Msg("strconv.ParseInt")
// 				}
// 				ans[j].NetCashFlow = parseInt
// 			}
// 			i = i + per
// 		}
// 		i = i + 1
// 	}
//
// 	return ans, nil
// }
//
// func parseBalance(name, code string, res []string) ([]items.Balance, error) {
// 	per := Period(res)
// 	date := ReportDate(res)
// 	if per != len(date) {
// 		log.Error().Err(fmt.Errorf("日期列不符"))
// 		return nil, fmt.Errorf("日期列不符")
// 	}
// 	ans := make([]items.Balance, per)
// 	for i := 0; i < len(date); i++ {
// 		ans[i].ReportingPeriod = date[i]
// 		ans[i].Code = code
// 		ans[i].Name = name
// 	}
// 	i := 0
// 	for i < len(res) {
// 		switch res[i] {
// 		case "货币资金":
// 			for j := 0; j < per; j++ {
// 				if parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64); err == nil {
// 					ans[j].MoneyFunds = int64(parseInt)
// 				}
// 			}
// 			i = i + per
// 		case "交易性金融资产":
// 			for j := 0; j < per; j++ {
// 				if parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64); err == nil {
// 					ans[j].TransFinance = int64(parseInt)
// 				}
// 			}
// 			i = i + per
// 		case "存货":
// 			for j := 0; j < per; j++ {
// 				if parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64); err == nil {
// 					ans[j].Stock = int64(parseInt)
// 				}
// 			}
// 			i = i + per
// 		case "短期借款":
// 			for j := 0; j < per; j++ {
// 				if parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64); err == nil {
// 					ans[j].ShortLoan = int64(parseInt)
// 				}
// 			}
// 			i = i + per
// 		case "长期借款":
// 			for j := 0; j < per; j++ {
// 				if parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64); err == nil {
// 					ans[j].LongLoan = int64(parseInt)
// 				}
// 			}
// 			i = i + per
//
// 		case "实收资本(或股本)", "股本", "实收资本":
// 			for j := 0; j < per; j++ {
// 				if parseInt, err := strconv.ParseFloat(strings.Replace(strings.Replace(res[i+j+1], ",", "", -1), "-", "", -1), 64); err == nil {
// 					ans[j].Capital = int64(parseInt)
// 				}
// 			}
// 		}
// 		i = i + 1
// 	}
//
// 	return ans, nil
// }

var (
	re = regexp.MustCompile("[0-9]+")
)
