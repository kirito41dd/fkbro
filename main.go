package main

import (
	"flag"
	"github.com/zshorz/ezlog"
	"github.com/zshorz/fkbro/bot"
	"github.com/zshorz/fkbro/btcinfo"
	"github.com/zshorz/fkbro/marketinfo"
	"github.com/zshorz/fkbro/util"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Log = ezlog.New(os.Stdout, "", ezlog.BitDefault, ezlog.LogAll)
var Bot *bot.Bot
var API *btcinfo.BTC_com_api
var HuobiAPI *marketinfo.Huobi

var configfile string

func main() {

	flag.StringVar(&configfile, "c", "config.json", "specify config file")
	flag.Parse()

	util.Config.Reload(configfile)

	API = btcinfo.NewBTC_com_api(util.Config.ApiHost, util.Config.CacheSize)
	HuobiAPI = marketinfo.NewHuobi(util.Config.Proxy)

	for  {
		bot, error := bot.NewBot(util.Config.BotToken, util.Config.Proxy)
		if error != nil {
			Log.Warn(error, "will retry...")
			<- time.After(5*time.Second)
		} else {
			Bot = bot
			break
		}
	}

	if !util.Config.Debug {
		Log.SetLogLevel(ezlog.LogInfo)
		API.Logger.SetLogLevel(ezlog.LogInfo)
		Bot.Logger.SetLogLevel(ezlog.LogInfo)
	}

	Log.Info("create bot success")


	Bot.Router.AddHandle("/help", help)
	Bot.Router.AddHandle("/newest", newest)
	Bot.Router.AddHandle("/q", q)
	Bot.Router.AddHandle("unknow", unknow)
	Bot.Router.AddHandle("/recent", recent)
	Bot.Router.AddHandle("/quotes", quotes)
	Bot.Router.AddHandle("/exchange", exchange)
	Bot.Router.AddHandle("/market", market)
	Bot.Router.AddHandle("/rss", rss)

	Bot.CallbackQueryRouter.AddHandle("showQuotes", showQuotes)

	LoadTemplate(util.Config.StaticPath)
	go doSignal()
	go live()
	Bot.Loop()
	Log.Info("good bye: api query cnt", API.QueryCnt)
}

func doSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGHUP)

	for sig := range sigs{
		if sig == syscall.SIGINT {
			liveExitChan <- 1
			Bot.ExitChan <- 1
		} else if sig == syscall.SIGHUP {
			;
		}
	}

}