package main

import (
	"github.com/zshorz/fkbro/btcinfo"
	"github.com/zshorz/fkbro/data"
	"github.com/zshorz/fkbro/util"
	"log"
	"testing"
	"time"
)

func Test_temp(t *testing.T) {
	LoadTemplate(util.Config.StaticPath)
	api := btcinfo.NewBTC_com_api(util.Config.ApiHost, util.Config.CacheSize)
	b,_ := api.GetBlockInfo("latest")
	b1, _ := api.GetBlockInfo("100")

	arr := make([]*btcinfo.Block,0)
	arr = append(arr, b)
	arr = append(arr, b1)


	t.Log(b.Extras)
	t.Log(ParseToString("recent", arr))

	tran , _ := api.GetTransactionInfo("9b5877539f6b422ae5c72e1a105defa4de5ecbc1d9175a0c30207d887b3864df")
	t.Log(ParseToString("transaction", tran))
}

func Test_alert(t *testing.T) {
	LoadTemplate(util.Config.StaticPath)
	alert := data.Alert{
		ID:        0,
		TimeStamp: 124,
		Symbol:    "btc",
		Hash:      "548448",
		Amount:    1234,
		AmountUsd: 0,
		FromAddr:  "",
		FromOwner: "fo",
		TomAddr:   "",
		TomOwner:  "tow",
	}
	t.Error(ParseToString("alert", &alert))
}

func Test_daily(t *testing.T) {
	data.Setup( "192.168.0.105:3306", "fkbro", "root", "123456")
	LoadTemplate(util.Config.StaticPath)
	now := time.Now()
	futher := now.AddDate(0,0,1)
	_ = time.Date(futher.Year(), futher.Month(), futher.Day(), 0,0,0,0, futher.Location())
	//<- time.After(time.Second * time.Duration(t0.Unix()-now.Unix()))
	last0 := time.Date(now.Year(), now.Month(), now.Day(), 0,0,0,0, now.Location())
	report := calc("btc", last0.Unix())
	for k,v := range report.Data {
		t.Log(k)
		t.Log(v)
	}
	report.Time = "#日报"
	text := ParseToString("report", report)
	//t.Log(text)
	log.Print(text)
	report = calc("usdt", 0)
	report.Time = "#日报"
	text = ParseToString("report", report)
	log.Print(text)
}








