package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zshorz/fkbro/data"
	"github.com/zshorz/fkbro/util"
	"math"
	"strings"
	"sync"
	"time"
)

var liveMap = make(map[string]int64) // username : ChatID
var liveMapLock sync.Mutex

var liveExitChan = make(chan int, 1)
var liveStateMap = make(map[string]*liveState)
type liveState struct {
	LastTime 	time.Time
	LastCursor	string
}


// 每日报告
type report struct {
	Name string // 币种
	Time string // 时间间隔
	Data map[string]*reportData
	IN   int64 // 流入交易所
	OUT  int64 // 流出交易所
	IN_usd   int64 // 流入交易所
	OUT_usd  int64 // 流出交易所
}
func (re *report) GetINOUT() int64 {
	ret := re.IN - re.OUT
	return ret
}
func (re *report) AbsUsd() int64 {
	return int64 (math.Abs(float64(re.IN_usd - re.OUT_usd)))
}
func (re *report) Abs() int64 {
	return int64 (math.Abs(float64(re.IN - re.OUT)))
}

type reportData struct {
	TotalIn 	int64
	TotalOut	int64
	UnknownIN	int64
	UnknownOut 	int64
	TotalIn_usd 	int64
	TotalOut_usd	int64
	UnknownIN_usd	int64
	UnknownOut_usd 	int64
}

func live() {
	liveStateMap["btc"] = &liveState{LastTime:time.Now(), LastCursor:"0-0-0"}
	liveStateMap["usdt"] = &liveState{LastTime:time.Now(), LastCursor:"0-0-0"}

	if util.Config.WhaleApikey == "" {
		return
	}
	data.Setup(util.Config.DbAddr, util.Config.DbName, util.Config.DbUser, util.Config.DbPasswd)
	go daily()
	for {
		if len(liveExitChan) > 0 {
			<-liveExitChan
			break
		}

		look("btc")
		look("usdt")

		<-time.After(60*time.Second)
	}
}

func look(currency string) { // btc usdc


	w := WhaleAPI.WhaleQuery(util.Config.WhaleApikey, currency, liveStateMap[currency].LastCursor, liveStateMap[currency].LastTime.Unix(), util.Config.LookMinV)

	if w == nil || w.Result != "success" || w.Cursor == "0-0-0" { // 没查到
		Log.Debug(w)
		if time.Now().Unix() - liveStateMap[currency].LastTime.Unix() > 60*60 && liveStateMap[currency].LastCursor == "0-0-0"{ // 60 分钟没查好
			liveStateMap[currency].LastTime = time.Now()
		}
	} else {
		liveStateMap[currency].LastCursor = w.Cursor
		Log.Debug(w)

		for _,trans := range w.Transactions {
			if trans.TransactionType != "transfer" {
				continue
			}
			if trans.From.Owner_type == "unknown" {
				trans.From.Owner = "unknown"
			}
			if trans.To.Owner_type == "unknown" {
				trans.To.Owner = "unknown"
			}

			al := data.Alert{
				TimeStamp: trans.Timestamp,
				Symbol:    trans.Symbol,
				Hash:      trans.Hash,
				Amount:    int64(trans.Amount+0.5),
				AmountUsd: int64(trans.Amount_usd+0.5),
				FromAddr:  trans.From.Address,
				FromOwner: trans.From.Owner,
				TomAddr:   trans.To.Address,
				TomOwner:  trans.To.Owner,
			}
			switch trans.Blockchain {
			case "bitcoin":
				al.SetURL("https://blockchair.com/zh/bitcoin/transaction/" + al.Hash)
			case "ethereum":
				al.SetURL("https://blockchair.com/zh/ethereum/transaction/0x" + al.Hash)
			case "tron":
				al.SetURL("https://tronscan.org/#/transaction/" + al.Hash)
			}
			Log.Debug(trans)
			if trans.From.Owner != trans.To.Owner{ // 双方账户不一样才处理
				err := al.Insert() // 存入数据库
				if err != nil && strings.HasPrefix(err.Error(), "Error 1062") { // 有重复,就不通知了
					continue
				}
				if int64(trans.Amount_usd) >= util.Config.AlertMinV {
					alert(&al) // 机器人推送
				}
			}


		}

	}
}

func alert(al *data.Alert) {
	liveMapLock.Lock()
	for _, id :=range liveMap {
		text := ParseToString("alert", al)
		msg := tgbotapi.NewMessage(id,text)
		msg.ChatID = id
		msg.ParseMode="Markdown"
		send(&msg, 1)
	}
	liveMapLock.Unlock()
}

func daily() {
	Log.Info("daily 任务开启")
	for {
		// 计算下一个0点
		now := time.Now()
		futher := now.AddDate(0,0,1)
		t0 := time.Date(futher.Year(), futher.Month(), futher.Day(), 0,0,0,0, futher.Location())
		dur := time.Second * time.Duration(t0.Unix()-now.Unix())
		Log.Info("等待执行", int64(dur/time.Second/60/60), "小时后")
		<- time.After(dur)
		last0 := time.Date(now.Year(), now.Month(), now.Day(), 0,0,0,0, now.Location())
		// 开始执行 0 点任务
		btcrep := calc("btc", last0.Unix(), 5000000) // 500w$ 才会打印
		usdtrep := calc("usdt", last0.Unix(), 5000000)
		btcrep.Time = "#日报"
		usdtrep.Time = "#日报"
		doReport(btcrep)
		doReport(usdtrep)

		t:= time.Now()
		if  t.Weekday() == 0 { // 一周任务
			nt := t.AddDate(0,0,-7)
			btcrep := calc("btc", nt.Unix(), 25000000)
			usdtrep := calc("usdt", nt.Unix(), 25000000)
			btcrep.Time = "#周报"
			usdtrep.Time = "#周报"
			doReport(btcrep)
			doReport(usdtrep)
			data.DeleteAlert(nt.Unix())
			Log.Info("周任务 本次完成")
		}

		Log.Info("daily 本次完成")
	}
}

func doReport(rep *report) {
	liveMapLock.Lock()
	for _, id :=range liveMap {
		text := ParseToString("report", rep)
		msg := tgbotapi.NewMessage(id,text)
		msg.ChatID = id
		msg.ParseMode="Markdown"
		send(&msg, 1)
	}
	liveMapLock.Unlock()
}

// 统计币种情况, 不会返回 nil
func calc(currency string, timestamp int64, minVal int64) *report{
	re := &report{
		Name: currency,
		Data: make(map[string]*reportData),
	}
	alerts := data.GetAlertsByTimeStamp(currency, timestamp)
	if alerts == nil {
		return re
	}
	// 初始化
	for _, alert := range alerts {
		if alert.FromOwner != "unknown" {
			re.Data[alert.FromOwner] = &reportData{}
		}
		if alert.TomOwner != "unknown" {
			re.Data[alert.TomOwner] = &reportData{}
		}
	}
	// 计算数据
	for _, alert := range alerts {
		// 计算支出
		if alert.FromOwner != "unknown" {
			data := re.Data[alert.FromOwner]
			data.TotalOut += alert.Amount
			data.TotalOut_usd += alert.AmountUsd
			if alert.TomOwner == "unknown" && alert.FromOwner != "tether treasury" { // 泰达币金库不应视为交易所
				data.UnknownOut += alert.Amount
				data.UnknownOut_usd += alert.AmountUsd
				re.OUT += alert.Amount // 流出交易所
				re.OUT_usd += alert.AmountUsd // 流出交易所
			}
		}
		// 计算收入
		if alert.TomOwner != "unknown" {
			data := re.Data[alert.TomOwner]
			data.TotalIn += alert.Amount
			data.TotalIn_usd += alert.AmountUsd
			if alert.FromOwner == "unknown" && alert.TomOwner != "tether treasury" {
				data.UnknownIN += alert.Amount
				data.UnknownIN_usd += alert.AmountUsd
				re.IN += alert.Amount // 流入交易所
				re.IN_usd += alert.AmountUsd // 流入交易所
			}
		}
	}
	// 交易额小于 minVal 美元的不显示
	for k,v := range re.Data {
		if v.TotalIn_usd + v.TotalOut_usd < minVal {
			delete(re.Data, k)
		}
	}
	return re
}

