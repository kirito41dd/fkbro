package btcinfo

import "fmt"

type Address struct {
	// TODO: 只用了关心的字段
	Address 		string 		`json:"address"`	// 地址
	Received		int64 		`json:"received"`	// 总接收
	Sent 			int64 		`json:"sent"`		// 总支出
	Balance 		int64 		`json:"balance"`	// 余额
	Tx_count		int64		`json:"tx_count"`   // 交易次数
	Last_tx			string 		`json:"last_tx"`	// 最后一次交易
	First_tx		string		`json:"first_tx"`	// 第一次交易
}

func (addr *Address) GetRecv() string {
	re := float64(float64(addr.Received) / float64(1e8))
	return fmt.Sprintf("%.8f", re)
}

func (addr *Address) GetSent() string {
	re := float64(float64(addr.Sent) / float64(1e8))
	return fmt.Sprintf("%.8f", re)
}

func (addr *Address) GetBalance() string {
	re := float64(float64(addr.Balance) / float64(1e8))
	return fmt.Sprintf("%.8f", re)
}