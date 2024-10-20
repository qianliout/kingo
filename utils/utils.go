package utils

import (
	"fmt"
	"reflect"
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
