package xingoudoc

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type doc struct {
	msg string
}

func (c *doc) Error()string{
	return c.msg
}

func GetApiDoc(dir ...string) (map[string]interface{},error){
	var apiDir string
	if len(dir) == 0{
		apiDir = "controller/"
	}
	data := make(map[string]interface{})
	files, err := ioutil.ReadDir(apiDir)
	if err != nil {
		return data,errors.New("目录不存在")
	}
	methodStart := 0
	for _, f := range files {
		fs, _ := os.Open(apiDir + f.Name())
		rd := bufio.NewReader(fs)
		temp := make(map[string]interface{})
		returnMap := make(map[string]interface{})
		paramMap := make(map[string]interface{})
		for {
			line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
			if err != nil || io.EOF == err {
				break
			}
			if len(line) < 2 {
				continue
			}
			if line[0:2] != "//" {
				continue
			}
			startTag, _ := regexp.MatchString("//@start", line)
			if startTag {
				if methodStart == 1 {
					panic("注释没有闭合标签")
				}
				methodStart = 1
				paramMap = make(map[string]interface{})
				returnMap = make(map[string]interface{})
				continue
			}
			endTag, _ := regexp.MatchString("//@end", line)
			if endTag {
				if methodStart != 3 {
					panic("解析解析必须在返回参数解析之后")
				}
				methodStart = 0
				temp["param"] = paramMap
				temp["return"] = returnMap
				data[temp["name"].(string)] = temp
				//log.Println("结束解析注释")
				continue
			}
			paramTag, _ := regexp.MatchString("//@param", line)
			if paramTag {
				if methodStart != 1 {
					panic("参数解析必须在开始标签之后")
				}
				methodStart = 2
				continue
			}
			returnTag, _ := regexp.MatchString("//@return", line)
			if returnTag {
				if methodStart != 2 {
					panic("返回参数解析必须在参数解析标签之后")
				}
				methodStart = 3
				continue
			}
			lineStr := line[2:]
			lineStr = strings.TrimSpace(lineStr)
			if methodStart == 1 {
				temp, _ = strToMap(lineStr)
				temp, _ = checkMethodHead(temp)
			}
			if methodStart == 2 {
				arr, _ := strToMap(lineStr)
				arr, _ = checkParam(arr)
				paramMap[arr["name"].(string)] = arr
			}
			if methodStart == 3 {
				arr, _ := strToMap(lineStr)
				arr, _ = checkReturn(arr)
				returnMap[arr["name"].(string)] = arr
			}
		}
		if methodStart == 1 {
			panic("注释缺失闭合标签")
		}
	}
	return data,nil
}

func checkParam(data map[string]interface{}) (map[string]interface{}, bool) {
	// name=username,type=string,info=书籍编号,must=1
	if _, ok := data["name"]; !ok {
		panic("请求参数字段必填")
	}
	if _, ok := data["type"]; !ok {
		panic("请求参数类型必填")
	}
	if _, ok := data["info"]; !ok {
		panic("请求参数名称必填")
	}
	if _, ok := data["must"]; !ok {
		data["must"] = "0"
	}
	if _, ok := data["value"]; !ok {
		data["value"] = "none"
	}
	if _, ok := data["rule"]; !ok {
		data["rule"] = "none"
	}
	return data, true
}

func checkReturn(data map[string]interface{}) (map[string]interface{}, bool) {
	// name=list,type=array,mock=1-5,info=记录集合,id=1
	if _, ok := data["name"]; !ok {
		panic("返回参数字段必填")
	}
	if _, ok := data["type"]; !ok {
		panic("返回参数类型必填")
	}
	if _, ok := data["info"]; !ok {
		panic("返回参数名称必填")
	}
	if _, ok := data["id"]; !ok {
		data["id"] = "0"
	}
	if _, ok := data["pid"]; !ok {
		data["pid"] = "0"
	}
	if _, ok := data["mock"]; !ok {
		data["mock"] = "none"
	}
	return data, true
}

func checkMethodHead(data map[string]interface{}) (map[string]interface{}, bool) {
	//title=测试接口,name=api/test2,method=ANY/GET/POST/DEL/PUT/OPTIONS,auth=AUTH/NONE
	if _, ok := data["title"]; !ok {
		panic("接口名称必填")
	}
	if _, ok := data["name"]; !ok {
		panic("接口路由必填")
	}
	if _, ok := data["auth"]; !ok {
		data["auth"] = "NONE"
	}
	if _, ok := data["method"]; !ok {
		data["method"] = "ANY"
	}
	return data, true
}

func strToMap(s string) (map[string]interface{}, bool) {
	data := make(map[string]interface{})
	arr := strings.Split(s, ",")
	for _, val := range arr {
		arr1 := strings.Split(val, "=")
		if len(arr1) != 2 {
			panic("注释书写错误")
		}
		data[arr1[0]] = arr1[1]
	}
	return data, true
}
