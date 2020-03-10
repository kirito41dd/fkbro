package test

import (
	"crypto/tls"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zshorz/fkbro/btcinfo"
	"github.com/zshorz/fkbro/util"
	"github.com/zshorz/ezlog"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test_(t *testing.T) {
	ti := time.Now()
	t.Log(ti.String())
	t.Log(ti.Format("2006-01-02 15:04:05"))
}

func Test_demo(t *testing.T) {
	ezlog.Info("fkbro demo test...")
	// 本地开发 使用http代理
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(util.Config.Proxy)
	}
	transport := &http.Transport{Proxy:proxy}
	client := &http.Client{}
	client.Transport = transport
	bot, err := tgbotapi.NewBotAPIWithClient(util.Config.BotToken, client)
	// bot, err := tgbotapi.NewBotAPI("1088231432:AAGFMLS3GVBh_P1a0u_IYm-NDWe6mfnS5r8")
	if err != nil {
		ezlog.Panic(err)
	}

	bot.RemoveWebhook()
	//bot.Debug = true
	ezlog.Debugf("Authorized on account: %s", bot.Self.UserName)
	//botName := bot.Self.UserName

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60


	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		ezlog.Info(update.UpdateID)
		ezlog.Debugf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		msg.ReplyToMessageID = 0//update.Message.MessageID
		ezlog.Info("replayto:", msg.ReplyToMessageID)
		ezlog.Debug("message type: ", update.Message.Chat.Type)
		if update.Message.Chat.Type != "private" && update.Message.Entities != nil &&len(*update.Message.Entities) > 0 {
			util.Logger.Debug("arr len:", len(*update.Message.Entities))
			entities := (*update.Message.Entities)[0]
			ezlog.Info(entities)
			if entities.Type != "mention" {
				continue
			}
			// @user
			user := update.Message.Text[entities.Offset+1 : entities.Length]
			util.Logger.Debug("user:", user)
			if user != bot.Self.UserName {
				continue
			}
			msg.Text = "这里是群组"
			_, err := bot.Send(msg)
			if err != nil {
				ezlog.Error(err)
			}
		} else if update.Message.Chat.Type == "private" {
			msg.Text = "这里是私聊"
			_, err := bot.Send(msg)
			if err != nil {
				ezlog.Error(err)
			}
		}

	}

}

func Test_btc_com_api(t *testing.T) {
	_ = &http.Transport{
		TLSClientConfig:        &tls.Config{
			InsecureSkipVerify:          true,
		},
	}

	client := http.Client{}

	req, err := http.NewRequest("GET", "https://chain.api.btc.com/v3/block/latest", nil)
	if err != nil {
		t.Error(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var response btcinfo.Response
	var block btcinfo.Block
	err = json.Unmarshal(body, &response)
	js, err := json.Marshal(response.Data)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(js, &block)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(body))
	t.Log(response)
	t.Log(block)

	t.Log("\n")
	req, err = http.NewRequest("GET", "https://chain.api.btc.com/v3/tx/b85e710221b6c14c47c9503e7ff4b10109573b582c1e4c82644084c7adf17f93?verbose=3", nil)
	if err != nil {
		t.Error(err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Error(err)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	js, err = json.Marshal(response.Data)
	var trans btcinfo.Transaction
	err = json.Unmarshal(js, &trans)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(body))
	t.Log(response)
	t.Log(trans)

}

func Test_re2(t *testing.T) {
	str := "/q sadg 9b5877539f6b422ae5c72e1a105defa4de5ecbc1d9175a0c30207d887b3864df  620555 123  0000000000000000000de95c2c8ccce2a3aee3f56d17d70627e4c9984bd378e5 "
	ss2 := strings.Fields(str)

	re , _ := regexp.Compile(`[\w]{1,}`)
	re.Longest()
	t.Log(re.String())
	regexp.MatchString(re.String(), str)
	//ss := re.FindStringSubmatch(str)
	for _, s := range ss2 {
		t.Log(s)
	}
}

func Test_kline(t *testing.T) {

	Proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(util.Config.Proxy)
	}
	transport := &http.Transport{Proxy:Proxy}
	client := http.Client{}
	if util.Config.Proxy != "" {
		client.Transport = transport
	}

	req, _ := http.NewRequest("GET", "https://api.huobi.pro/market/history/kline", nil)
	q := req.URL.Query()
	q.Add("symbol", "btcusdt")
	q.Add("period", "60min")
	q.Add("size", "1")
	req.URL.RawQuery = q.Encode()

	resp, _ := client.Do(req)
	t.Log(req.URL.String())
	data,_ := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{})
	t.Log(data)
	json.Unmarshal(data, &m)
	t.Log(m)
}