package marketinfo

/// https://huobiapi.github.io/docs/spot/v1/cn/
type Response struct {
	Status		string		`json:"status"`		// 接口返回状态
	Ch 			string 		`json:"ch"`			// 接口数据对应的数据流。部分接口没有对应数据流因此不返回此字段
	Ts 			int64 		`json:"ts"`			// 返回的时间戳
	Data 		interface{}	`json:"data"`		// 数据主体
}
