package xingoudoc

import (
	"strconv"
	"strings"
)

// // name=data,type=BigScreen,info=code,must=1
type DocParam struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Info  string      `json:"code"`
	Must  int         `json:"must"`
	Rule  string      `json:"rule"`
	Limit string      `json:"limit"`
	Map   string      `json:"map"`
	Value interface{} `json:"value"`
}

func (dp *DocParam) Analyze(str string) {
	arr := strings.Split(str, ",")
	for _, val := range arr {
		arr1 := strings.Split(val, "=")
		if arr1[0] == "name" {
			dp.Name = arr1[1]
		}
		if arr1[0] == "type" {
			dp.Type = arr1[1]
		}
		if arr1[0] == "info" {
			dp.Info = arr1[1]
		}
		if arr1[0] == "must" {
			dp.Must, _ = strconv.Atoi(arr1[1])
		}
		if arr1[0] == "value" {
			dp.Value = arr1[1]
		}
		if arr1[0] == "rule" {
			dp.Rule = arr1[1]
		}
		if arr1[0] == "limit" {
			dp.Limit = arr1[1]
		}
		if arr1[0] == "map" {
			dp.Map = arr1[1]
		}
	}
	dp.check()
}

func (dp *DocParam) check() {
	if dp.Name == "" {
		panic("请求参数字段必填")
	}
	if dp.Type == "" {
		panic("请求参数类型必填")
	}
	if dp.Info == "" {
		panic("请求参数名称必填")
	}
	if dp.Must == 0 {
		dp.Must = 0
	}
	if dp.Value == nil {
		dp.Value = nil
	}
	// 检查规则
	if dp.Rule == "" {
		dp.Rule = "none"
	}
	// 字符串校验 strict 表示开启特殊字符检查包含则不通过
	if dp.Limit == "" {
		dp.Limit = "strict"
	}
	// 参数分组
	if dp.Map == "" {
		dp.Map = "none"
	}
}
