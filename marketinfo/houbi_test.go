package marketinfo

import (
	"github.com/zshorz/fkbro/util"
	"testing"
)

func Test_KLine(t *testing.T) {
	huobi := NewHuobi(util.Config.Proxy)
	t.Log(util.Config.Proxy)
	klines, err := huobi.GetKLine("btcusdt", "60min", 1)
	t.Log(klines[0], err)
	ex := huobi.GetExchange()
	t.Log(ex)
}
