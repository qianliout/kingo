package utils

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"regexp"
	"strconv"
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

func ReportDate(res []string) []Report {
	ans := make([]Report, 0)
	for i := 0; i < len(res); i++ {
		if res[i] == "报表日期" {
			for j := i + 1; j < len(res); j++ {
				if res[j] == "一、营业总收入" || res[j] == "一、经营活动产生的现金流量" || res[j] == "流动资产" || res[j] == "资产" || res[j] == "一、营业收入" {
					break
				}
				ans = append(ans, Report{ReportPeriod: res[j]})
			}
		}
	}
	// 处理一下，处理成202006 202009 这类这样便于查询
	// 2020-03-31
	for i, ch := range ans {
		split := strings.Split(ch.ReportPeriod, "-")
		if len(split) < 2 {
			continue
		}
		ans[i].ReportPeriod = strings.Join(split[:2], "")
		ans[i].Year = GetInt64(split[0])
		ans[i].Month = GetInt64(split[1])
	}
	return ans
}

// ParsePeriodCnt 查询当前的报表中有几个时间区
func ParsePeriodCnt(res []string) int {
	start := 0

	for i := 0; i < len(res); i++ {
		if res[i] == "报表日期" {
			start = i
		}
		if res[i] == "一、营业总收入" || res[i] == "一、经营活动产生的现金流量" || res[i] == "流动资产" || res[i] == "资产" || res[i] == "一、营业收入" {
			return i - start - 1
		}
	}
	return 0
}

func SetField(obj interface{}, fieldName string, value interface{}) error {
	// 获取对象的反射值
	val := reflect.ValueOf(obj)

	// 检查是否是指针
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 检查是否是结构体
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("obj must be a struct or pointer to struct")
	}

	// 获取字段
	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field %s not found in struct", fieldName)
	}

	// 检查字段是否可设置
	if !field.CanSet() {
		return fmt.Errorf("field %s is not settable", fieldName)
	}

	// 设置字段值
	fieldVal := reflect.ValueOf(value)
	if !fieldVal.Type().ConvertibleTo(field.Type()) {
		return fmt.Errorf("value type %s is not convertible to field type %s", fieldVal.Type(), field.Type())
	}

	field.Set(fieldVal.Convert(field.Type()))
	return nil
}

func GenerateUUID64(str string) int64 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(str))
	ans := h.Sum32()
	return int64(ans)
}

func GetReportYear(str string) string {
	// 2024-03-31
	split := strings.Split(str, "-")
	if len(split) > 0 {
		return split[0]
	}
	return ""
}

func GetInt64(a string) int64 {
	i, _ := strconv.ParseInt(a, 10, 64)
	return i
}

type Report struct {
	ReportPeriod string
	Year         int64
	Month        int64
}
