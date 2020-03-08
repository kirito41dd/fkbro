package btcinfo

// https://btc.com/api-doc
// 交易
type Transaction struct {
	// TODO: 只用了关心的字段
	Confirmations		int64 	`json:"confirmations"`	// 确认数
	Block_height		int64 	`json:"block_height"` 	// 所在块高度
	Block_time			int64 	`json:"block_time"`		// 所在块时间
	Created_at			int64 	`json:"created_at"`		// 系统处理时间，没有业务含义
	Fee					int64 	`json:"fee"`			// 该交易手续费
	Hash 				string	`json:"hash"`			// 交易hash
	Inputs				[]Input	`json:"inputs"`			// 输入
	Inputs_count		int64 	`json:"inputs_count"`	// 输入数量
	Inputs_value		int64 	`json:"inputs_value"`	// 输入金额
	Outputs				[]Output `json:"outputs"`		// 输出
	Outputs_count		int64 	`json:"outputs_count"` 	// 输出数量
	Outputs_value		int64 	`json:"outputs_value"`	// 输出数量
	Size 				int64 	`json:"size"`			// 交易体积
}


type Input struct {
	// TODO: 只用了关心的字段
	Prev_addresses		[]string	`json:"prev_addresses"` // 输入地址
	Prev_position		int64 		`json:"prev_position"`	// 前向交易的输出位置
	Prev_tx_hash		string		`json:"prev_tx_hash"`	// 前向交易hash
	Prev_value			int64 		`json:"prev_value"`		// 输入金额
}

type Output struct {
	// TODO: 只用了关心的字段
	Addresses	[]string 	`json:"addresses"`	// 输出地址
	Value 		int64 		`json:"value"`			// 输出金额
}
