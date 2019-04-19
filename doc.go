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

type NoteDoc struct {
    Doc         map[string]interface{}
    FileDir     string
    FileName    string
    JsonName    string
    KeyName     string
    PackageName string
    DocJson     string
}

func (c *NoteDoc) MapToHtml() string {
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
            log.Print(vp)
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
            log.Print(vr)
        }
        html = html + "</table>"
        html = html + "<hr>"
        log.Print(value)
    }
    return html
}

func (c *NoteDoc) MapToFile() error {
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
    docJson, _ := json.Marshal(c.Doc)
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
    log.Println(saveFileName, c.FileDir, c.FileName)
    //time.Sleep(2*time.Second)
    // 拼接稳定文件
    doc1 := string(docJson)
    doc1 = strings.Replace(doc1, "\"", "'", -1)
    doc1 = "\"" + doc1 + "\""
    doc1 = "package " + c.PackageName + " \r\nvar " + c.KeyName + " = " + doc1
    _, err = io.WriteString(f, doc1) //写入文件(字符串)
    return err
}

func (c *NoteDoc) MapToJson() error {
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

func (c *NoteDoc) GetApiDoc(apiDir string) error {
    if c.Doc == nil {
        c.Doc = make(map[string]interface{})
    }
    files, err := ioutil.ReadDir(apiDir)
    if err != nil {
        return errors.New("目录不存在")
    }
    methodStart := 0
    errMsg := ""
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
            // 去除空格 换行等特殊符号
            line = strings.Replace(line, " ", "", -1)
            line = strings.Replace(line, "\r", "", -1)
            line = strings.Replace(line, "\n", "", -1)
            //log.Println(line)
            if len(line) < 4 {
                continue
            }
            if line[0:2] != "//" && line[0:4] != "func" {
                continue
            }

            if line[0:4] == "func" && methodStart == 4 {
                startNum := 4
                endNum := strings.Index(line, "(")
                otherEndNum := strings.LastIndex(line, "(")
                if endNum != otherEndNum {
                    startNum = strings.Index(line, ")") + 1
                    endNum = otherEndNum
                }
                temp["func"] = line[startNum:endNum]
                c.Doc[temp["name"].(string)] = temp
                methodStart = 0
                continue
            }

            startTag, _ := regexp.MatchString("//@start", line)
            if startTag {
                if methodStart == 1 {
                    errMsg = "注释没有闭合标签"
                    break
                }
                methodStart = 1
                paramMap = make(map[string]interface{})
                returnMap = make(map[string]interface{})
                continue
            }
            endTag, _ := regexp.MatchString("//@end", line)
            if endTag {
                if methodStart != 3 {
                    errMsg = "解析解析必须在返回参数解析之后"
                    break
                }
                methodStart = 4
                temp["param"] = paramMap
                temp["return"] = returnMap
                continue
            }
            paramTag, _ := regexp.MatchString("//@param", line)
            if paramTag {
                if methodStart != 1 {
                    errMsg = "参数解析必须在开始标签之后"
                    break
                }
                methodStart = 2
                continue
            }
            returnTag, _ := regexp.MatchString("//@return", line)
            if returnTag {
                if methodStart != 2 {
                    errMsg = "返回参数解析必须在参数解析标签之后"
                    break
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
        if errMsg != "" {
            break
        }
        if methodStart == 1 {
            errMsg = "注释缺失闭合标签"
        }
    }
    if errMsg != "" {
        return errors.New(errMsg)
    }
    return nil
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
    // 检查规则
    if _, ok := data["rule"]; !ok {
        data["rule"] = "none"
    }
    // 字符串校验 strict 表示开启特殊字符检查包含则不通过
    if _, ok := data["check"]; !ok {
        data["check"] = "strict"
    }
    // 参数分组
    if _, ok := data["map"]; !ok{
        data["map"] = "none"
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
    } else {
        var methodArr = []string{"GET", "POST", "PUT", "PATCH", "HEAD", "OPTIONS", "DELETE", "CONNECT", "ANY", "TRACE"}
        method := data["method"].(string)
        ok := 0
        for _, v := range methodArr {
            if method == v {
                ok = 1
                break
            }
        }
        if ok == 0 {
            panic("方法书写错误")
        }
    }
    if _, ok := data["func"]; !ok {
        data["func"] = "none"
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
