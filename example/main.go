package main

import (
	"log"

	"github.com/stupboy/xingoudoc"
)

func main() {
	var Doc xingoudoc.GinDoc
	var err error
	Doc.FileName = "doc.go"
	Doc.FileDir = "example/utli/"
	Doc.KeyName = "DocJsonString"
	Doc.PackageName = "utli"
	err = Doc.GetApiDoc("example/controller/")
	log.Println(Doc.DocSlice)
	Doc.ToFile()
	if err != nil {
		panic(err)
	}
}
