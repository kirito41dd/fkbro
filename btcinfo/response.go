package btcinfo


type Response struct {
	Data 		interface{} `json:"data"` 	// 返回数据
	Err_no		int `json:"err_no"`			// 错误码 0 正常 1 找不到该资源 2 参数错误
	Error_msg 	string `json:"error_msg"`	// 错误信息，供调试使用。如果没有错误，则此字段不出现。
}
