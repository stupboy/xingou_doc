package xingoudoc

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type GinDoc struct {
	Doc         map[string]interface{}
	DocSlice    []ApiDoc
	FileDir     string
	FileName    string
	JsonName    string
	KeyName     string
	PackageName string
	DocJson     string
}

func (c *GinDoc) MapToHtml() string {
	var html string
	for key, value := range c.Doc {
		valueMap := value.(map[string]interface{})
		html = html + key + " " + valueMap["title"].(string)
		html = html + "<table>"
		html = html + "<tr><td>参数</td></tr>"
		for kp, vp := range valueMap["param"].(map[string]interface{}) {
			html = html + "<tr>"
			vpMap := vp.(map[string]interface{})
			html = html + "<td>" + kp + "</td>"
			html = html + "<td>" + vpMap["info"].(string) + "</td>"
			html = html + "<td>" + vpMap["type"].(string) + "</td>"
			html = html + "<td>" + vpMap["must"].(string) + "</td>"
			html = html + "<td>" + vpMap["value"].(string) + "</td>"
			html = html + "<td>" + vpMap["rule"].(string) + "</td>"
			//html = html + kp + vpMap["info"].(string) + "<br>"
			html = html + "</tr>"
			//print(vp)
		}
		//html = html + "<br>返回<br>"
		//html = html + "<table>"
		html = html + "<tr><td>返回</td></tr>"
		for kr, vr := range valueMap["return"].(map[string]interface{}) {
			html = html + "<tr>"
			vpMap := vr.(map[string]interface{})
			html = html + "<td>" + kr + "</td>"
			html = html + "<td>" + vpMap["info"].(string) + "</td>"
			html = html + "<td>" + vpMap["type"].(string) + "</td>"
			html = html + "<td>" + vpMap["mock"].(string) + "</td>"
			html = html + "<td>" + vpMap["id"].(string) + "</td>"
			html = html + "<td>" + vpMap["pid"].(string) + "</td>"
			//html = html + kp + vpMap["info"].(string) + "<br>"
			html = html + "</tr>"
			//log.Print(vr)
		}
		html = html + "</table>"
		html = html + "<hr>"
		//log.Print(value)
	}
	return html
}

func (c *GinDoc) ToFile() error {
	var err error
	if c.FileName == "" {
		return errors.New("文件名称不存在")
	}
	if c.KeyName == "" {
		return errors.New("变量名称不存在")
	}
	if c.PackageName == "" {
		return errors.New("包名不存在")
	}
	docJson, _ := json.Marshal(c.DocSlice)
	var f *os.File
	saveFileName := c.FileDir + c.FileName
	exist := true
	if _, err := os.Stat(saveFileName); os.IsNotExist(err) {
		exist = false
	}
	if exist { //如果文件存在
		f, _ = os.OpenFile(saveFileName, os.O_RDWR, 0666) //打开文件
	} else {
		f, _ = os.Create(saveFileName) //创建文件
	}
	println(saveFileName, c.FileDir, c.FileName)
	//time.Sleep(2*time.Second)
	// 拼接稳定文件
	doc1 := string(docJson)
	doc1 = strings.Replace(doc1, "\"", "'", -1)
	doc1 = "\"" + doc1 + "\""
	doc1 = "package " + c.PackageName + " \r\nvar " + c.KeyName + " = " + doc1
	_, err = io.WriteString(f, doc1) //写入文件(字符串)
	return err
}

func (c *GinDoc) MapToJson() error {
	var err error
	if c.JsonName == "" {
		return errors.New("文件名称不存在")
	}
	docJson, _ := json.Marshal(c.Doc)
	var f *os.File
	defer f.Close()
	saveFileName := c.JsonName
	exist := true
	if _, err := os.Stat(saveFileName); os.IsNotExist(err) {
		exist = false
	}

	if exist { //如果文件存在
		f, _ = os.OpenFile(saveFileName, os.O_RDWR, 0666) //打开文件
		_ = os.Remove(saveFileName)
	}

	f, _ = os.Create(saveFileName) //创建文件

	doc1 := string(docJson)
	_, err = io.WriteString(f, doc1) //写入文件(字符串)
	return err
}

// 获取每个方法文档
func (c *GinDoc) getMethodDoc() {

}

// 获取每个文件文档
func (c *GinDoc) getFileDoc(apiDir string, f os.FileInfo) error {
	var (
		methodTurn, startNum int
		endNum, otherEndNum  int
		line, lineStr        string
		fs                   *os.File
		rd                   *bufio.Reader
		err                  error
		startTag             bool
		lineNum              int
		temDoc               ApiDoc
		temParam             DocParam
		temReturn            DocReturn
	)
	fs, _ = os.Open(apiDir + f.Name())
	rd = bufio.NewReader(fs)
	for {
		lineNum++
		line, err = rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		// 去除空格 换行等特殊符号
		line = strings.Replace(line, " ", "", -1)
		line = strings.Replace(line, "\r", "", -1)
		line = strings.Replace(line, "\n", "", -1)
		line = strings.Replace(line, `"`, "", -1)
		if len(line) < 4 { //长度低于4个字符跳过
			continue
		}
		// 不是注释内容和方法第一行跳过
		if line[0:2] != "//" && line[0:4] != "func" {
			continue
		}
		switch line {
		case "//@start":
			if methodTurn != methodAnalyzeOver {
				log.Fatal("报错文件:", f.Name(), ",报错行数:", lineNum, ",开始标签必须在结束标签之后")
			}
			methodTurn = methodAnalyzeStart
			temDoc.Params = make(map[string]DocParam)
			temDoc.Returns = make(map[string]DocReturn)
			temParam = DocParam{}
			continue
		case "//@param":
			if methodTurn != methodAnalyzeHead {
				log.Fatal("报错文件:", f.Name(), ",报错行数:", lineNum, ",参数必须在头部标签之后")
			}
			methodTurn = methodAnalyzeParam
			continue
		case "//@return":
			if methodTurn != methodAnalyzeHead && methodTurn != methodAnalyzeParam {
				log.Fatal("报错文件:", f.Name(), ",报错行数:", lineNum, ",返还标签必须在头部标签或参数标签之后")
			}
			methodTurn = methodAnalyzeReturn
			continue
		case "//@end":
			if methodTurn != methodAnalyzeParam && methodTurn != methodAnalyzeReturn && methodTurn != methodAnalyzeHead {
				log.Fatal("报错文件:", f.Name(), ",报错行数:", lineNum, ",结束标签必须在其他标签之后")
			}
			methodTurn = methodAnalyzeEnd
			// temp["param"] = paramMap
			// temp["return"] = returnMap
			// temDoc.Params = DocParms
			// temDoc.Returns = DocReturns
			continue
		}
		switch methodTurn {
		case methodAnalyzeStart: //检查头部
			lineStr = line[2:]
			lineStr = strings.TrimSpace(lineStr)
			temDoc.Head.Analyze(lineStr)
			methodTurn = methodAnalyzeHead
		case methodAnalyzeParam:
			lineStr = line[2:]
			lineStr = strings.TrimSpace(lineStr)
			temParam.Analyze(lineStr)
			temDoc.Params[temParam.Name] = temParam
			temParam = DocParam{}
		case methodAnalyzeReturn:
			lineStr = line[2:]
			lineStr = strings.TrimSpace(lineStr)
			if methodTurn == methodAnalyzeReturn {
				// arr, ok := strToMap(lineStr)
				// if !ok {
				// 	continue
				// }
				// arr, _ = checkReturn(arr)
				temReturn.Analyze(lineStr)
				temDoc.Returns[temReturn.Name] = temReturn
				temReturn = DocReturn{}
				// returnMap[arr["name"].(string)] = arr
			}

		case methodAnalyzeEnd:
			if line[0:4] != "func" {
				continue
			}
			startNum = 4
			endNum = strings.Index(line, "(")
			otherEndNum = strings.LastIndex(line, "(")
			if endNum != otherEndNum {
				startNum = strings.Index(line, ")") + 1
				endNum = otherEndNum
			}
			// temp["func"] = line[startNum:endNum]
			temDoc.Head.Func = line[startNum:endNum]
			// 结束分档分析
			c.DocSlice = append(c.DocSlice, temDoc)
			temDoc = ApiDoc{}
			temDoc.Params = make(map[string]DocParam)
			temDoc.Returns = make(map[string]DocReturn)
			methodTurn = methodAnalyzeOver
		case methodAnalyzeOver: //重新开发分析 next 分析头部
			// 注释开始标记
			startTag, _ = regexp.MatchString("//@start", line)
			if !startTag {
				continue
			}
			methodTurn = methodAnalyzeStart
			// paramMap = make(map[string]interface{})
			// returnMap = make(map[string]interface{})
			temParam = DocParam{}
		}
	}
	if methodTurn == methodAnalyzeHead || methodTurn == methodAnalyzeStart {
		log.Fatal("报错文件:", f.Name(), ",报错行数:", lineNum, ",注释缺失闭合标签")
	}
	return nil
}

func (c *GinDoc) GetApiDoc(apiDir string) error {
	var (
		files []os.FileInfo
		f     os.FileInfo
		err   error
	)
	if c.Doc == nil {
		c.Doc = make(map[string]interface{})
	}
	files, err = ioutil.ReadDir(apiDir)
	if err != nil {
		return errors.New("目录不存在")
	}
	for _, f = range files {
		err = c.getFileDoc(apiDir, f)
	}
	if err != nil {
		return err
	}
	return nil
}
