package xingoudoc

import "strings"

// title=测试接口,name=api/test2,method=ANY/GET/POST/DEL/PUT/OPTIONS,auth=AUTH/NONE
type DocHead struct {
	Name   string `json:"name"`
	Title  string `json:"title"`
	Method string `json:"method"`
	Auth   string `json:"auth"`
	Func   string `json:"func"`
}

func (dh *DocHead) Analyze(str string) {
	arr := strings.Split(str, ",")
	for _, val := range arr {
		arr1 := strings.Split(val, "=")
		if arr1[0] == "name" {
			dh.Name = arr1[1]
		}
		if arr1[0] == "title" {
			dh.Title = arr1[1]
		}
		if arr1[0] == "method" {
			dh.Method = arr1[1]
		}
		if arr1[0] == "auth" {
			dh.Auth = arr1[1]
		}
	}
	dh.check()
}

func (dh *DocHead) check() {
	if dh.Title == "" {
		panic("接口名称必填")
	}
	if dh.Name == "" {
		panic("接口路由必填")
	}
	if dh.Auth == "" {
		dh.Auth = "NONE"
	}
	if dh.Method == "" {
		dh.Method = "ANY"
	} else {
		var methodArr = []string{"GET", "POST", "PUT", "PATCH", "HEAD", "OPTIONS", "DELETE", "CONNECT", "ANY", "TRACE"}
		// method := data["method"].(string)
		ok := 0
		for _, v := range methodArr {
			if dh.Method == v {
				ok = 1
				break
			}
		}
		if ok == 0 {
			panic("方法书写错误")
		}
	}
	if dh.Func == "" {
		dh.Func = "none"
	}
	// return data, true
}
