package utils

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ParseNameCode(selection *goquery.Selection) (name, code string) {
	t := selection.Text()
	rcode := regexp.MustCompile(`\(.*?\)`)
	co := rcode.FindString(t)
	code = strings.Replace(strings.Replace(co, "(", "", -1), ")", "", -1)
	rname := regexp.MustCompile(`(.*?)\(`)
	na := rname.FindString(t)
	name = strings.Replace(strings.Replace(na, "(", "", -1), ")", "", -1)
	return name, code
}

func ReportDate(res []string) []string {
	ans := make([]string, 0)
	for i := 0; i < len(res); i++ {
		if res[i] == "报表日期" {
			for j := i + 1; j < len(res); j++ {
				if res[j] == "一、营业总收入" || res[j] == "一、经营活动产生的现金流量" || res[j] == "流动资产" || res[j] == "资产" || res[j] == "一、营业收入" {
					return ans
				}
				ans = append(ans, res[j])
			}
		}
	}
	return ans
}

// 查询当前的报表中有几个时间区
func ParsePeriod(res []string) int {
	start := 0

	for i := 0; i < len(res); i++ {
		if res[i] == "报表日期" {
			start = i
		}
		if res[i] == "一、营业总收入" || res[i] == "一、经营活动产生的现金流量" || res[i] == "流动资产" || res[i] == "资产" || res[i] == "一、营业收入" {
			return i - start - 1
		}
	}
	return len(res)
}
