package marketinfo

import (
	"fmt"
)

// 这个map里的币种会被当作主流币,都与usdt比较
var Symbol_name = map[string]string{
	"btcusdt" : "BTC/USDT",
	"ethusdt" : "ETH/USDT",
	"bchusdt" : "BCH/USDT",
	"bsvusdt" : "BSV/USDT",
	"ltcusdt" : "LTC/USDT",
	"eosusdt" : "EOS/USDT",
	"xrpusdt" : "XRP/USDT",
	"etcusdt" : "ETC/USDT",
}

var period_time = map[string]string{
	"1day" : "1日",
	"60min": "1小时",
}

type KLine struct {
	// TODO: 只使用了关注的字段
	Id 			int64 		`json:"id"`			// 调整为新加坡时间的时间戳，单位秒，并以此作为此K线柱的id
	Open 		float32 	`json:"open"`		// 开盘
	Close 		float32 	`json:"close"`		// 收盘价
	High 		float32 	`json:"high"`		// 最高
	Low 		float32 	`json:"low"`		// 最低

	// 下面都是为了在模板中显示加入的变量
	symbol 		string 		// "btcusdt"
	name 		string 		// 比如 "BTC/USDT"
	time 		string      // 比如 "当天"
	yesterday	float32		// 昨日收盘价
}

func (kl *KLine) GetName() string {
	return kl.name
}

func (kl *KLine) GetTime() string {
	return kl.time
}

func (kl *KLine) GetSymbol() string {
	return kl.symbol
}

func (kl *KLine) SetYesterday(y float32)  {
	kl.yesterday = y
}
func (kl *KLine) GetChange() string {
	var v float32 = (kl.Close - kl.yesterday) / kl.yesterday * 100
	if v > 0 {
		return fmt.Sprintf("+%.2f", v)
	}
	return fmt.Sprintf("%.2f", v)
}
