package util

import (
	"encoding/json"
	"github.com/zshorz/ezlog"
	"io/ioutil"
	"os"
)

type GlobalObj struct {
	BotToken	string	`json:"bot_token"`
	Proxy		string	`json:"proxy"`
	CacheSize	int 	`json:"cache_size"`
	ApiHost		string 	`json:"api_host"`
	StaticPath	string	`json:"static_path"`
	Debug 		bool 	`json:"debug"`
	BotOwner 	string 	`json:"bot_owner"` // username 如 yesiare 不要@
	WhaleApikey string  `json:"whale_apikey"` // https://docs.whale-alert.io/#introduction
	AlertMinV	int64 	`json:"alert_min_v"`
	LookMinV	int64 	`json:"look_min_v"`
	DbAddr 		string  `json:"db_addr"` // 192.168.0.105:3306
	DbUser 		string 	`json:"db_user"`
	DbPasswd 	string  `json:"db_passwd"`
	DbName		string 	`json:"db_name"`
	Rss 		[]string `json:"rss"` // 提前订阅的用户列表，服务启动自动订阅，不加@
}



var Config *GlobalObj
var Logger *ezlog.EzLogger

func (g *GlobalObj) Reload(file string) {

	if !Exists(file) {
		return
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		Logger.Warn(err)
	}
	//data = []byte(fmt.Sprintf("json:%s\n", data))
	err = json.Unmarshal(data, &g)
	if err != nil {
		Logger.Panic(err)
	}
}

func init() {
	Config =  &GlobalObj{
		BotToken: "your_token",
		Proxy:    "socks5://127.0.0.1:1080",
		CacheSize: 5000,
		ApiHost:   "https://chain.api.btc.com/v3" ,
		StaticPath: "static",
		Debug:		true,
		BotOwner: 	"yesiare",
		WhaleApikey: "",
		AlertMinV: 2000000,
		LookMinV:  1000000,
	}
	Logger = ezlog.New(os.Stdout, "", ezlog.BitDefault, ezlog.LogAll)
}

func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}