package bot

import (
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zshorz/ezlog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Bot struct {
	TgBot	*tgbotapi.BotAPI
	Proxy	string
	Token 	string
	Router	*BotRouter
	CallbackQueryRouter *BotRouter
	ExitChan chan interface{}
	Logger *ezlog.EzLogger
}

func NewBot(token string, proxy string) (*Bot, error){
	Proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxy)
	}
	transport := &http.Transport{Proxy:Proxy}
	client := &http.Client{}
	if proxy != "" {
		client.Transport = transport
	}
	bot := &Bot{}
	var err error
	bot.TgBot, err = tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		return nil, err
	}
	bot.Proxy = proxy
	bot.Token = token
	bot.Router = NewBotRouter()
	bot.CallbackQueryRouter = NewBotRouter()
	bot.ExitChan = make(chan interface{})
	bot.Logger = ezlog.New(os.Stdout, "", ezlog.BitDefault, ezlog.LogAll)
	return bot, nil
}

func (bot *Bot) Dispatch(update *tgbotapi.Update) {
	if update.Message != nil {
		if update.Message.Entities != nil {
			for _, entity := range *update.Message.Entities {
				if entity.Type == "bot_command" {
					bot.Logger.Debug("offse:len", entity.Offset, entity.Length, "msg" , update.Message.Text)
					url := update.Message.Text[entity.Offset : entity.Offset+entity.Length]
					i := strings.Index(url, "@") // 提取命令 /help@kirito_testonly_bot
					if i != -1 {
						url = url[:i]
					}
					bot.Router.DoHandle(url, update)
					return
				}
			}
		}
		bot.Router.DoHandle("unknow", update)

	} else if update.CallbackQuery != nil {
		if update.CallbackQuery.Message != nil {
			mp  := make(map[string]string)
			err := json.Unmarshal([]byte(update.CallbackQuery.Data), &mp)
			if err != nil {
				bot.Logger.Error(err)
				return
			}
			if mp["cmd"] == "" {
				return
			}
			bot.Logger.Debug(mp)
			bot.CallbackQueryRouter.DoHandle(mp["cmd"], update)
		}
	} else {
		bot.Logger.Warn("其他消息类型", update)
		//bot.Router.DoHandle("unknow", update)
	}

}

func (bot *Bot) Loop() {
	bot.Logger.Info("Looping...")
	bot.TgBot.RemoveWebhook()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updateChan, err := bot.TgBot.GetUpdatesChan(updateConfig)
	if err != nil {
		bot.Logger.Error("Loop", err)
		return
	}
	for {
		select {
		case <- bot.ExitChan:
			bot.Logger.Info("exit loop")
			bot.TgBot.StopReceivingUpdates()
			return
		case update := <- updateChan:
			u := update
			js, _ := json.Marshal(u)
			bot.Logger.Debug("new update:", string(js))
			bot.Dispatch(&update)
		}
	}
}


