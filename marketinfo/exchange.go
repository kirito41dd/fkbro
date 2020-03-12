package marketinfo

/// TODO: 不要写死
/// https://api.exchangerate-api.com/v4/latest/USD

type Exchange struct {
	Base 		string 				`json:"base"`
	Rates		map[string]float32	`json:"rates"`
	Date 		string				`json:"date"`
	
	Time_last_updated int64 `json:"time_last_updated"`

	btc2usdt		float32
}


func (ex *Exchange) Get1USD2CNY() float32 {
	return ex.Rates["CNY"]
}

func (ex *Exchange) Get1BTC2CNY() float32 {
	return ex.Rates["CNY"] * ex.btc2usdt
}

func (ex *Exchange) SetBTC2USDT(f float32){
	ex.btc2usdt = f
}
