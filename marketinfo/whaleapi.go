package marketinfo

import (
	"encoding/json"
	"fmt"
	"github.com/zshorz/ezlog"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type WhaleAPI struct {
	Client 		*http.Client
	Logger 		*ezlog.EzLogger
}

func NewWhaleAPI() *WhaleAPI {
	return &WhaleAPI{
		Client: &http.Client{},
		Logger: ezlog.New(os.Stdout, "", ezlog.BitDefault, ezlog.LogAll),
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
	TransactionType string 		`json:"transaction_type"`
}

type WhaleAddr struct {
	Address 	string 		`json:"address"`
	Owner 		string 		`json:"owner"`
	Owner_type 	string 		`json:"owner_type"`
}

// arg like Cq4aEiBAVjYFsqIHjg8GBLq86ct0TZue btc xxx-xxx-xxx
func (whale *WhaleAPI)WhaleQuery(apiKey, currency, cursor string, start , min_value int64) *WhaleAlert {
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
	whale.Logger.Debug(req.URL)
	resp, err := whale.Client.Do(req)
	if err != nil {
		whale.Logger.Debug(err)
		return nil
	}
	data , _ := ioutil.ReadAll(resp.Body)
	var w WhaleAlert
	err = json.Unmarshal(data, &w)
	if err != nil {
		whale.Logger.Debug(err)
		return nil
	}
	return &w
}

