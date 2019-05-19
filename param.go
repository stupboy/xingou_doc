package xingoudoc

import (
    "errors"
    "fmt"
    "log"
    "strconv"
    "strings"
)

// 获取参数转存
type AllParam struct {
    // 存储参数
    Data map[string]interface{}
    // 存储参数错误
    ErrMap map[string]interface{}
    // 查询条件字符串
    Condition string
    //W2 string
}

type PageInfo struct {
    Page int `json:"page"`
    Size int `json:"size"`
}

type WhereInfo struct {
    Key   string `json:"key"`
    Alias string `json:"alias"`
    Type  int    `json:"type"` //类型  等于  like 大于  小
}

type WhereString string

func (whereString WhereString) Analyze() (key string,way string,alias string) {
    // 键|= >= <= in key@tabe
    arr := strings.Split(string(whereString),"|")
    if len(arr) == 1{
        key = arr[0]
        way = "="
    }
    return key,way,alias
}

func (allParam *AllParam) InitData(data interface{}) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = errors.New("初始化数据失败")
        }
    }()
    allParam.Data = data.(map[string]interface{})
    return err
}

func (allParam *AllParam)InitW1(){
    if allParam.Condition == "" {
        allParam.Condition = " 1=1 "
    }
}

func (allParam *AllParam)AddFormat(format string,argus ...interface{}){
    allParam.InitW1()
    v := fmt.Sprintf(format,argus...)
    if strings.Contains(v,"%!") {
        panic("参数类型填写错误")
    }
    if len(v) > 0 {
        allParam.Condition += v
    }
}

func (allParam *AllParam)AddOther(key string,value string){
    if allParam.Condition == ""{
        allParam.Condition = " 1=1 "
    }
    allParam.Condition += " and " + key + " = " + value
}

func (allParam *AllParam)AddEqual(key string,argus ...string){
    var (
        p   string
        val interface{}
    )
    if _,ok := allParam.Data[key]; !ok{
        return
    }
    if allParam.Condition == ""{
        allParam.Condition = " 1=1 "
    }
    val = allParam.Data[key]
    switch val.(type) {
    case string:
        p = "'" + val.(string) + "'"
    case int:
        p = strconv.Itoa(val.(int))
    default:
        p = "''"
    }
    if len(argus) == 0{
        allParam.Condition += " and " + key + "=" + p
    }else{
        allParam.Condition += " and " + argus[0] + "=" + p
    }
}
func (allParam *AllParam)AddLike(key string,argus ...string){
    var (
        p   string
        val interface{}
    )
    if _,ok := allParam.Data[key]; !ok{
        return
    }
    if allParam.Condition == ""{
        allParam.Condition = " 1=1 "
    }
    val = allParam.Data[key]
    switch val.(type) {
    case string:
        p = "'%" + val.(string) + "%'"
    case int:
        p = "'%" + strconv.Itoa(val.(int)) + "%'"
    default:
        p = "''"
    }
    if len(argus) == 0{
        allParam.Condition += " and " + key + " like " + p
    }else{
        allParam.Condition += " and " + argus[0] + " like " + p
    }
}


// 构造查询条件  username
func (allParam AllParam) CreateWhere(param ...WhereInfo){

}

// 增加查询条件
func (allParam AllParam) AddWhere(param ...WhereString){
    for _,str := range param{
        key,way,alias := str.Analyze()
        log.Println(key,way,alias)
    }
}

// 获取分页数据
func (allParam AllParam) Page() (info PageInfo){
    info.Page = allParam.Int("page",1)
    info.Size = allParam.Int("size",10)
    return info
}

// 请求参数格式化
func (allParam AllParam) String(key string, argus ...string) (param string) {
    defer func() {
        if r := recover(); r != nil {
            allParam.ErrMap[key] = r
            if len(argus) > 0 {
                param = argus[0]
            } else {
                param = ""
            }
        }
    }()
    if val, ok := allParam.Data[key]; ok {
        param = val.(string)
    }
    return param
}

func (allParam AllParam) Int(key string, argus ...int) (param int) {
    defer func() {
        if r := recover(); r != nil {
            allParam.ErrMap[key] = r
            if len(argus) > 0 {
                param = argus[0]
            } else {
                param = 0
            }
        }
    }()
    if val, ok := allParam.Data[key]; ok {
        param = val.(int)
    }
    return param
}