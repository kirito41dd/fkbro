package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zshorz/fkbro/btcinfo"
	"github.com/zshorz/fkbro/data"
	"github.com/zshorz/fkbro/marketinfo"
	"github.com/zshorz/fkbro/util"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 测试专用
func test(update *tgbotapi.Update) {
	if (util.Config.Debug){
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = "Markdown"

		alert := data.Alert{
			ID:        0,
			TimeStamp: 124,
			Symbol:    "btc",
			Hash:      "548448",
			Amount:    1234,
			AmountUsd: 55464468,
			FromAddr:  "",
			FromOwner: "fo",
			TomAddr:   "",
			TomOwner:  "tow",
		}
		alert.SetURL("https://www.blockchain.com/btc/tx/b3b6cda500a9fc9b4ed0df5d7cd07f747edcd7488bf4dc4cc97047e6e2354a95")
		msg.Text = ParseToString("alert", &alert)

		send(&msg, 5)
	}
}

func now(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ParseMode = "Markdown"
	timenow := time.Now()
	timenow = time.Date(timenow.Year(),timenow.Month(),timenow.Day(),0,0,0,0,timenow.Location())
	var want string = "btc"
	fds := strings.Fields(update.Message.Text)
	for _,s := range fds {
		if s[0] != '/' && s[0] != '@' {
			want = s
			break
		}
	}

	rep := calc(want, timenow.Unix(), 5000000)
	rep.Time = "#简报"

	msg.Text = ParseToString("report", rep)

	send(&msg, 5)
}

func rss(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	if update.Message.From.UserName != util.Config.BotOwner {
		return
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	fields := strings.Fields(update.Message.Text)
	arg := ""
	for _, s := range fields {
		if s[0] == '@' {
			arg = s[1:]
		}
	}
	if arg == "" {
		liveMapLock.Lock()
		msg.Text = "当前订阅:"
		for k, _ := range liveMap {
			msg.Text += "\n@" + k
		}
		liveMapLock.Unlock()
		send(&msg, 5)
		return
	}
	chat, err := Bot.TgBot.GetChat(tgbotapi.ChatConfig{
		ChatID:             0,
		SuperGroupUsername: "@"+arg,
	})
	if err != nil {
		Log.Debug(err)
		msg.Text = "获取群组或频道失败!"
		send(&msg, 5)
		return
	}
	id := chat.ID
	Log.Debug("chatId", id)

	liveMapLock.Lock()
	_, ok := liveMap[arg]
	if ok {
		delete(liveMap, arg)
		msg.Text = "@" + arg + " 取消订阅"
	} else {
		liveMap[arg] = id
		msg.Text = "@" + arg + " 加入订阅"
	}
	liveMapLock.Unlock()

	send(&msg, 5)
}

func help(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "@" + update.Message.From.UserName + " - " + update.Message.From.FirstName + "\n"
	msg.Text += ParseToString("help", nil)
	msg.ParseMode = "Markdown"
	send(&msg, 5)
}

func newest(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "@" + update.Message.From.UserName + " - " + update.Message.From.FirstName + "\n"
	block, err := getNewestBlock()
	if err != nil {
		Log.Error("GetBlockInfo:", err)
		return
	}
	msg.Text += ParseToString("block", block)
	msg.ParseMode = "Markdown"
	send(&msg, 5)
}

func recent(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "@" + update.Message.From.UserName + " - " + update.Message.From.FirstName + "\n"
	msg.ParseMode = "Markdown"

	arr := make([]*btcinfo.Block, 0)
	latestb, err := getNewestBlock()
	if err != nil {
		Log.Error(err)
		return
	}
	arr = append(arr, latestb)
	// TODO: 使用单次查询获得所有区块
	for i := 1; i < 5; i++ {
		b, err := API.GetBlockInfo( strconv.Itoa(int(latestb.Height) - i) )
		if err != nil {
			Log.Error(err)
			break
		}
		arr = append(arr, b)
		<-time.After(time.Millisecond * 100) // 间隔100毫秒
	}

	msg.Text += ParseToString("recent", arr)
	send(&msg, 5)
}

func q(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "@" + update.Message.From.UserName + " - " + update.Message.From.FirstName + "\n"
	msg.ParseMode = "Markdown"
	fields := strings.Fields(update.Message.Text)
	var want string
	for _, want = range fields {
		if strings.HasPrefix(want, "/") || strings.HasPrefix(want, "@") {
			want = ""
			continue
		} else {
			break
		}
	}

	re , _ := regexp.Compile(`[\w]{1,}`)
	re.Longest()
	want = re.FindString(want)

	var iface interface{}
	var err error
	if len(want) == 0 {
		Log.Debug("want len = 0 return")
		return
	}
	if len(want) < 20 || strings.HasPrefix(want, "000000") { // block 块
		iface, err = API.GetBlockInfo(want)
		msg.Text += ParseToString("block", iface)
	} else if len(want) > 20 && len(want) < 40 { // address 地址
		iface, err = API.GetAddressInfo(want)
		msg.Text += ParseToString("address", iface)
	} else if len(want) == 64 { // transaction 交易
		iface, err = API.GetTransactionInfo(want)
		msg.Text += ParseToString("transaction", iface)
	} else {
		return
	}
	if err != nil {
		Log.Debug(err)
		return
	}
	send(&msg, 5)
}

// 把逻辑独立处理，给callbackquery用
func _quotes(arg string) string {
	klines, err := HuobiAPI.GetKLine(arg, "1day", 2)
	if err != nil {
		return ""
	}
	klines[0].SetYesterday(klines[1].Close)
	return ParseToString("quotes", klines[0])
}

func quotes(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "@" + update.Message.From.UserName + " - " + update.Message.From.FirstName + "\n"
	msg.ParseMode = "Markdown"

	fields := strings.Fields(update.Message.Text)
	var want string
	for _, s := range fields {
		if strings.HasPrefix(s, "/") || strings.HasPrefix(s, "@") {
			continue
		} else {
			want = s
			break
		}
	}
	if want == "" {
		want = "btcusdt"
	}
	msg.Text += _quotes(want)
	send(&msg, 5)
}

func exchange(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "@" + update.Message.From.UserName + " - " + update.Message.From.FirstName + "\n"
	msg.ParseMode = "Markdown"

	kline, err := HuobiAPI.GetKLine("btcusdt", "1day", 1)
	ex := HuobiAPI.GetExchange()

	if err != nil || ex == nil {
		Log.Debug(err, "or", "ex is nil")
		return
	}
	ex.SetBTC2USDT(kline[0].Close)
	msg.Text += ParseToString("exchange", ex)
	send(&msg, 5)
}

func market(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = "@" + update.Message.From.UserName + " - " + update.Message.From.FirstName + "\n"
	msg.ParseMode = "Markdown"

	var klines = make([]*marketinfo.KLine,0)

	for k,_ := range marketinfo.Symbol_name {
		kline, err := HuobiAPI.GetKLine(k, "1day", 2)
		if err != nil || len(kline) != 2 {
			continue
		}
		kline[0].SetYesterday(kline[1].Close)
		klines = append(klines, kline[0])
	}
	length := len(klines)
	for i := 0; i< length; i++ {
		for j := 0; j < length-1-i; j++ {
			if klines[j].Close < klines[j+1].Close {
				klines[j], klines[j+1] = klines[j+1], klines[j]
			}
		}
	}


	rows := make([][]tgbotapi.InlineKeyboardButton, len(klines))

	for k, v := range klines {
		rows[k] = make([]tgbotapi.InlineKeyboardButton,1)
		text := v.GetName() + " 涨幅: " + v.GetChange()+ "% 价格: " + fmt.Sprintf("%.3f", v.Close) + " .............."
		mp := make(map[string]string)
		mp["cmd"] = "showQuotes"
		mp["arg"] = v.GetSymbol()
		data, err := json.Marshal(mp)
		if err != nil {
			Log.Error(err)
			return
		}
		rows[k][0] = tgbotapi.NewInlineKeyboardButtonData(text, string(data))
	}

	var numKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)
	// 生成键盘

	msg.ReplyMarkup = numKeyboard

	msg.Text += ParseToString("market", klines)
	send(&msg, 5)
}

func unknow(update *tgbotapi.Update) {
	q(update)
}

var newestBlockCache btcinfo.Block
var newestBlockTime time.Time
var newestBlockLock sync.Mutex
func getNewestBlock() (*btcinfo.Block, error) {
	t := time.Now()
	d := t.Sub(newestBlockTime)
	if d > 30*time.Second {
		b, err := API.GetBlockInfo("latest")
		if err != nil {
			return nil, err
		}
		newestBlockLock.Lock()
		newestBlockCache = *b
		newestBlockTime = time.Now()
		newestBlockLock.Unlock()
		return b, err
	} else {
		b := newestBlockCache
		return &b , nil
	}

}

func send(msg *tgbotapi.MessageConfig, duration int) {
	if len(msg.Text) == 0 {
		return
	}
	msg.DisableWebPagePreview = true // 不生成连接预览
	cnt := 20
	for ; cnt > 0; cnt-- {
		_, err := Bot.TgBot.Send(msg)
		if err != nil {
			Log.Debug("send msg:", err, " will retry...", cnt-1)
			<-time.After(time.Duration(duration)*time.Second)
		} else {
			break
		}
	}
	//调试用，不要每次加载
	if util.Config.Debug {
		LoadTemplate(util.Config.StaticPath)
	}
}

func addRss(username string) {
	chat, err := Bot.TgBot.GetChat(tgbotapi.ChatConfig{
		ChatID:             0,
		SuperGroupUsername: "@"+username,
	})
	if err != nil {
		return
	}
	liveMapLock.Lock()
	_, ok := liveMap[username]
	if ok {
		delete(liveMap, username)
	} else {
		liveMap[username] = chat.ID
	}
	liveMapLock.Unlock()
}

