package btcinfo

import (
	"fmt"
	"time"
)

type Block struct {
	Height 				int64 	`json:"height"` 			// 块高度
	Version 			int64 	`json:"version"`			// 块版本
	Mrkl_root			string	`json:"mrkl_root"`			// Merkle Root 默克尔根
	Curr_max_timestamp 	int64 	`json:"curr_max_timestamp"` // 块最大时间戳
	Timestamp			int64 	`json:"timestamp"`			// 块时间戳
	Bits 				int64 	`json:"bits"`
	Nonce				int64		`json:"nonce"`				// 随机数
	Hash 				string	`json:"hash"`				// 块哈希
	Prev_block_hash		string	`json:"prev_block_hash"`	// 前块hash
	Next_block_hash		string	`json:"next_block_hash"`	// 后块hash
	Size 				int64 	`json:"size"`
	Pool_difficulty		int64 	`json:"pool_difficulty"`	// 矿池难度
	Difficulty			float32 `json:"difficulty"`			// 块难度
	Tx_count			int64 	`json:"tx_count"`			// 块奖励
	Reward_block		int64 	`json:"reward_block"`		// 块奖励
	Reward_fees			int64 	`json:"reward_fees"` 		// 块手续费
	Created_at			int64 	`json:"Created_at"`			// 系统处理时间，无业务含义
	Confirmations		int64 	`json:"confiramtions"`		// 确认数
	Extras				BlockExtras	`json:"extras"`
}

type BlockExtras struct {
	Relayed_by	string	`json:"relayed_by"` // 播报方
	Pool_name 	string 	`json:"pool_name"`
	Pool_link	string	`json:"pool_link"`

}

func (b *Block) GetReword() string {
	re := float64(float64(b.Reward_block+b.Reward_fees) / float64(1e8))
	return fmt.Sprintf("%.8f", re)
}

func (b *Block) GetPastTimeByMinute() int64 {
	t := time.Unix(time.Now().Unix()-b.Timestamp, 0)
	return t.Unix()/60
}