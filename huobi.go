package main

import (
	"encoding/json"
	"github.com/zshorz/ezlog"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

/// https://huobiapi.github.io/docs/spot/v1/cn/
type Response struct {
	Status		string		`json:"status"`		// 接口返回状态
	Ch 			string 		`json:"ch"`			// 接口数据对应的数据流。部分接口没有对应数据流因此不返回此字段
	Ts 			int64 		`json:"ts"`			// 返回的时间戳
	Data 		interface{}	`json:"data"`		// 数据主体
}

type KLine struct {
	// TODO: 只使用了关注的字段
	Id 			int64 		`json:"id"`			// 调整为新加坡时间的时间戳，单位秒，并以此作为此K线柱的id
	Open 		float32 	`json:"open"`		// 开盘
	Close 		float32 	`json:"close"`		// 收盘价
	High 		float32 	`json:"high"`		// 最高
	Low 		float32 	`json:"low"`		// 最低
}

var client *http.Client

type Huobi struct {
	Client 		*http.Client
	Proxy		string
	Logger 		*ezlog.EzLogger
}

func NewHuobi(proxy string) *Huobi {
	Proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxy)
	}
	Transport := &http.Transport{Proxy:Proxy}

	huobi := &Huobi{
		Client: &http.Client{},
		Proxy:  proxy,
		Logger: ezlog.New(os.Stdout, "", ezlog.BitDefault, ezlog.LogAll),
	}

	if proxy != "" {
		huobi.Client.Transport = Transport
	}
	return huobi
}

// btcusdt 60min 1
func (h *Huobi) GetKLine(symbol, period string, size int) ([]*KLine, error) {
	req, _ := http.NewRequest("GET", "https://api.huobi.pro/market/history/kline", nil)
	q := req.URL.Query()
	q.Add("symbol", symbol)
	q.Add("period", period)
	q.Add("size", strconv.Itoa(size))
	req.URL.RawQuery = q.Encode()

	resp, err := h.query(req)
	if err != nil {
		return nil, err
	}
	kline := make([]*KLine,0)
	js, err := json.Marshal(resp.Data)
	if err != nil {
		h.Logger.Debug(err)
		return nil, err
	}
	err = json.Unmarshal(js, &kline)
	if err != nil {
		h.Logger.Debug(err)
		return nil, err
	}
	return kline, nil
}

func (h *Huobi) query(req *http.Request) (*Response, error) {
	resp, err := h.Client.Do(req)
	if err != nil {
		h.Logger.Debug(err, resp)
		return nil, err
	}
	var ret Response
	data, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, &ret)
	if err != nil {
		h.Logger.Debug(err)
		return nil, err
	}
	return &ret, nil
}
