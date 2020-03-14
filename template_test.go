package main

import (
	"github.com/zshorz/fkbro/btcinfo"
	"github.com/zshorz/fkbro/util"
	"testing"
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
	w := WhaleTrans{
		Blockchain: "",
		Symbol:     "",
		Hash:       "",
		Timestamp:  0,
		Amount:     0,
		Amount_usd: 0,
		From:       WhaleAddr{},
		To:         WhaleAddr{},
	}
	w.Amount_usd = 1234567
	w.Amount = 123
	w.From.Owner_type= "from"
	w.To.Owner_type = "to"
	w.Hash = "asd"
	str := ParseToString("alert", w)
	t.Log(str)
}




