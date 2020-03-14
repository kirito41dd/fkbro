package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zshorz/fkbro/util"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var liveMap = make(map[string]int64) // username : ChatID
var liveMapLock sync.Mutex
var liveClient http.Client

var lastTime time.Time
var lastCursor string = "0-0-0"
var liveExitChan = make(chan int, 1)

func live() {
	if util.Config.WhaleApikey == "" {
		return
	}
	lastTime = time.Now()
	for {
		if len(liveExitChan) > 0 {
			<-liveExitChan
			break
		}
		w := WhaleQuery(util.Config.WhaleApikey, "btc", lastCursor, lastTime.Unix(), 2000000) // 2000000

		if w == nil || w.Result != "success" || w.Cursor == "0-0-0" { // 没查到
			Log.Debug(w)
		} else {
			lastCursor = w.Cursor
			Log.Debug(w)

			for _,trans := range w.Transactions {
				trans.Symbol = strings.ToUpper(trans.Symbol)
				trans.Amount = math.Floor(trans.Amount)
				trans.Amount_usd = math.Floor(trans.Amount_usd)
				text := ParseToString("alert", trans)
				liveMapLock.Lock()
				for _, id :=range liveMap {
					msg := tgbotapi.NewMessage(id,text)
					msg.ChatID = id
					msg.ParseMode="Markdown"
					go send(&msg, 5)
				}
				liveMapLock.Unlock()
			}

		}
		<-time.After(60*time.Second)
	}
}


// https://docs.whale-alert.io/#transactions
type WhaleAlert struct {
	Result 		string 		`json:"result"` // success
	Cursor 		string		`json:"cursor"`
	Count 		int 		`json:"count"`
	Transactions []*WhaleTrans `json:"transactions"`
}

type WhaleTrans struct {
	Blockchain 		string 		`json:"blockchain"`
	Symbol 			string 		`json:"symbol"`
	Hash 			string 		`json:"hash"`
	Timestamp 		int64 		`json:"timestamp"`
	Amount 			float64 	`json:"amount"`
	Amount_usd 		float64  	`json:"amount_usd"`
	From 			WhaleAddr 	`json:"from"`
	To 				WhaleAddr 	`json:"to"`
}

type WhaleAddr struct {
	Address 	string 		`json:"address"`
	Owner_type 	string 		`json:"owner_type"`
}

// arg like Cq4aEiBAVjYFsqIHjg8GBLq86ct0TZue btc xxx-xxx-xxx
func WhaleQuery(apiKey, currency, cursor string, start , min_value int64) *WhaleAlert {
	req, _ := http.NewRequest("GET", "https://api.whale-alert.io/v1/transactions", nil)
	q := req.URL.Query()
	q.Add("api_key", apiKey)
	q.Add("currency", currency)
	str := fmt.Sprintf("%d", start)
	q.Add("min_value", strconv.Itoa(int(min_value)))
	q.Add("start", str)
	if cursor != "0-0-0" {
		q.Add("cursor", cursor)
	}
	req.URL.RawQuery = q.Encode()
	Log.Debug(req.URL)
	resp, err := liveClient.Do(req)
	if err != nil {
		Log.Debug(err)
		return nil
	}
	data , _ := ioutil.ReadAll(resp.Body)
	var w WhaleAlert
	err = json.Unmarshal(data, &w)
	if err != nil {
		Log.Debug(err)
		return nil
	}
	return &w
}
