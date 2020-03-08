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
