package marketinfo

import (
	"encoding/json"
	"errors"
	"github.com/zshorz/ezlog"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)



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

// 获取汇率，用的不是火币api
func (h *Huobi) GetExchange() *Exchange {
	req, err := http.NewRequest("GET","https://api.exchangerate-api.com/v4/latest/USD", nil)
	resp, err := h.Client.Do(req)
	if err != nil {
		h.Logger.Debug(err)
		return nil
	}
	data, _ := ioutil.ReadAll(resp.Body)
	var ex Exchange
	err = json.Unmarshal(data, &ex)
	if err != nil {
		h.Logger.Debug(err)
		return nil
	}
	return &ex
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
	for _, k := range kline {
		k.name = Symbol_name[symbol]
		if k.name == "" {
			k.name = symbol
		}
		k.symbol = symbol
		k.time = period_time[period]
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
	if ret.Status != "ok" {
		h.Logger.Debug("return status =", ret.Status)
		return nil, errors.New("status != ok")
	}
	return &ret, nil
}
