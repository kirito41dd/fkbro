package main

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

var Temp map[string]*template.Template = make(map[string]*template.Template)
var FuncMap = make(template.FuncMap)

func GetBtcValueString(v... int64) string {
	var total int64
	for _, j := range v {
		total = total + j
	}
	return fmt.Sprintf("%.8f", float64(total)/1e8)
}

func GetPastTime(b int64) string {
	t := time.Unix(time.Now().Unix()-b, 0)
	total_min := t.Unix()/60
	day := total_min/(60*24);
	total_min -= day*(60*24)
	hour := total_min/60
	total_min -= hour*60
	var str string
	year := day / 365
	if year > 0 {
		return fmt.Sprintf("%d年前", year)
	}
	if day >0 {
		str += fmt.Sprintf("%d天", day)
	}
	if day > 0 || hour > 0 {
		str += fmt.Sprintf("%d小时", hour)
	}
	str += fmt.Sprintf("%d分钟前", total_min)
	return str
}

func LoadTemplate(root string){
	FuncMap["GetBtcValueString"] = GetBtcValueString
	FuncMap["GetPastTime"] = GetPastTime
	if root[len(root)-1:] != "/" {
		root = root + "/"
	}
	loadtmp("block", root + "block.tmp")
	loadtmp("transaction", root + "transaction.tmp")
	loadtmp("address", root + "address.tmp")
	loadtmp("help", root + "help.tmp")
	loadtmp("recent", root + "recent.tmp")
	loadtmp("quotes", root + "quotes.tmp")
	loadtmp("exchange", root + "exchange.tmp")
	loadtmp("market", root + "market.tmp")
}

func ParseToString(key string, i interface{}) (str string) {
	t, ok := Temp[key]
	if !ok {
		return ""
	}
	buf := bytes.Buffer{}
	err := t.Execute(&buf, i)
	if err != nil {
		Log.Debug(err)
		return ""
	}
	return string(buf.Bytes())
}

func loadtmp(key string, path string) {
	// new  里面要传文件的 basename
	// https://stackoverflow.com/questions/49043292/error-template-is-an-incomplete-or-empty-template
	swa := template.New(key+".tmp").Funcs(FuncMap)

	t, err := swa.ParseFiles(path)
	//t, err := template.ParseFiles(path)
	if err != nil {
		Log.Error(err)
		return
	}
	t.Funcs(FuncMap)
	Temp[key] = t
}
