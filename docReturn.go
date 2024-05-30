package xingoudoc

import "strings"

type DocReturn struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Info string `json:"info"`
}

// name=code,type=int,info=返回状态
func (dr *DocReturn) Analyze(str string) {
	arr := strings.Split(str, ",")
	for _, val := range arr {
		arr1 := strings.Split(val, "=")
		if arr1[0] == "name" {
			dr.Name = arr1[1]
		}
		if arr1[0] == "type" {
			dr.Type = arr1[1]
		}
		if arr1[0] == "info" {
			dr.Info = arr1[1]
		}
	}
	dr.check()
}

func (dr *DocReturn) check() {

}
