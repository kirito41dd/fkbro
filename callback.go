package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zshorz/fkbro/btcinfo"
	"github.com/zshorz/fkbro/util"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func help(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.Text = ParseToString("help", nil)
	msg.ParseMode = "Markdown"
	send(&msg, 5)
}

func newest(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	block, err := getNewestBlock()
	if err != nil {
		Log.Error("GetBlockInfo:", err)
		return
	}
	msg.Text = ParseToString("block", block)
	msg.ParseMode = "Markdown"
	send(&msg, 5)
}

func recent(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	msg.ParseMode = "Markdown"

	arr := make([]*btcinfo.Block, 0)
	latestb, err := getNewestBlock()
	if err != nil {
		Log.Error(err)
		return
	}
	arr = append(arr, latestb)
	// TODO: 使用单次查询获得所有区块
	for i := 1; i < 9; i++ {
		b, err := API.GetBlockInfo( strconv.Itoa(int(latestb.Height) - i) )
		if err != nil {
			Log.Error(err)
			return
		}
		arr = append(arr, b)
	}

	msg.Text = ParseToString("recent", arr)
	send(&msg, 5)
}

func q(update *tgbotapi.Update) {
	Log.Debug("recv msg:", update.Message.Text, "chatID:", update.Message.Chat.ID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
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
		msg.Text = ParseToString("block", iface)
	} else if len(want) > 20 && len(want) < 40 { // address 地址
		iface, err = API.GetAddressInfo(want)
		msg.Text = ParseToString("address", iface)
	} else if len(want) == 64 { // transaction 交易
		iface, err = API.GetTransactionInfo(want)
		msg.Text = ParseToString("transaction", iface)
	} else {
		return
	}
	if err != nil {
		Log.Debug(err)
		return
	}
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
	cnt := 20
	for ; cnt > 0; cnt-- {
		_, err := Bot.TgBot.Send(msg)
		if err != nil {
			Log.Debug("send msg:", err, " will retry...", cnt)
			<-time.After(time.Duration(duration)*time.Second)
		} else {
			break
		}
	}
	/// TODO: 调试用，不要每次加载，发行时删除这个
	LoadTemplate(util.Config.StaticPath)
}
