package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func Test_WhaleAlert(t *testing.T) {


	client := http.Client{}
	req, _ := http.NewRequest("GET", "https://api.whale-alert.io/v1/transactions", nil)
	q := req.URL.Query()
	q.Add("api_key", "Cq4aEiBAVjYFsqIHjg8GBLq86ct0TZue")
	q.Add("currency", "btc")
	d, _ := time.ParseDuration("-1m")
	_ = time.Now().Add(d * 59)
	str := fmt.Sprintf("%d",0 )
	q.Add("min_value", "500000")
	q.Add("start", str)
	q.Add("cursor", "221a847f-221a847f-5e6cb85a")
	req.URL.RawQuery = q.Encode()
	t.Log(req.URL)
	resp, _ := client.Do(req)
	data , _ := ioutil.ReadAll(resp.Body)
	var w WhaleAlert
	json.Unmarshal(data, &w)
	t.Log(w)
}
