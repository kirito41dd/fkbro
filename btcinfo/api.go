package btcinfo

import (
	"encoding/json"
	"errors"
	"github.com/zshorz/ezlog"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type IAPI interface {
	GetAddressInfo(addr string) (*Address, error)
	GetTransactionInfo(id string) (*Transaction, error)
	GetBlockInfo(want string) (*Block, error)
}

// 实现 IAPI 接口
type BTC_com_api struct {
	Client 	http.Client
	Host	string 		// https://chain.api.btc.com/v3
	Logger  *ezlog.EzLogger
	QueryCnt int64		// 记录查询次数，调试用

	cacheSize int

	addressCache map[string]*Address
	addressLock sync.RWMutex
	addressKeyChan chan string

	transactionCache map[string]*Transaction
	transactionLock sync.RWMutex
	TransactionKeyChan chan string

	blockCache	map[string]*Block
	blockLock	sync.RWMutex
	blockKeyChan chan string
}

func NewBTC_com_api(host string, cacheSize int) (*BTC_com_api) {
	api := &BTC_com_api{}

	api.Client = http.Client{}
	api.Host = host
	api.Logger = ezlog.New(os.Stdout, "", ezlog.BitDefault, ezlog.LogAll)

	api.cacheSize = cacheSize
	api.addressCache = make(map[string]*Address)
	api.addressKeyChan = make(chan string, cacheSize)
	api.transactionCache = make(map[string]*Transaction)
	api.TransactionKeyChan = make(chan string, cacheSize)
	api.blockCache = make(map[string]*Block)
	api.blockKeyChan = make(chan string, cacheSize)


	return api
}


func (api *BTC_com_api) GetAddressInfo(addr string) (*Address, error) {
	api.addressLock.RLock()
	ad, ok := api.addressCache[addr]
	api.addressLock.RUnlock()
	if ok {
		return ad, nil
	}
	// TODO: 不存在,远程调用获取
	data, err := api.doquery(api.Host + "/address/" + addr)
	if err != nil {
		return nil, err
	}
	ret, err := getAddress(data)
	if err != nil {
		return nil, err
	}
	// 加入到缓冲
	api.addressLock.Lock()
	if (len(api.addressKeyChan) >= api.cacheSize) { // 缓冲满了 删除一个
		k := <- api.addressKeyChan
		delete(api.addressCache, k)
	}
	api.addressCache[addr] = ret
	api.addressLock.Unlock()

	return ret, err;

	return nil, errors.New("GetAddressInfo unknow error")
}
func (api *BTC_com_api) GetTransactionInfo(id string) (*Transaction, error) {
	api.transactionLock.RLock()
	trans, ok := api.transactionCache[id]
	api.transactionLock.RUnlock()
	if ok {
		return trans, nil
	}

	// TODO： 不存在
	data, err := api.doquery(api.Host + "/tx/" + id)
	if err != nil {
		return nil, err
	}
	ret, err := getTransaction(data)
	if err != nil {
		return nil, err
	}
	// 加入到缓冲
	api.transactionLock.Lock()
	if (len(api.TransactionKeyChan) >= api.cacheSize) { // 缓冲满了 删除一个
		k := <- api.TransactionKeyChan
		delete(api.transactionCache, k)
	}
	api.transactionCache[id] = ret
	api.transactionLock.Unlock()

	return ret, err;
	return nil, errors.New("GetTransactionInfo unknow error")
}
func (api *BTC_com_api) GetBlockInfo(want string) (*Block, error) {

	api.blockLock.RLock()
	b, ok := api.blockCache[want]
	api.blockLock.RUnlock()
	if ok {
		return b, nil;
	}
	// todo:
	data, err := api.doquery(api.Host + "/block/" + want)
	if err != nil {
		api.Logger.Debug(data,err)
		return nil, err
	}
	ret, err := getBlock(data)
	if err != nil {
		api.Logger.Error(ret,err)
		return nil, err
	}

	if want == "latest" { // 不缓存
		return ret, err
	} else {
		// 加入到缓冲
		api.blockLock.Lock()
		if (len(api.blockKeyChan) >= api.cacheSize) { // 缓冲满了 删除一个
			k := <- api.blockKeyChan
			delete(api.blockCache, k)
		}
		api.blockCache[want] = ret
		api.blockLock.Unlock()
		return ret, nil
	}

	return nil, errors.New("GetBlockInfo unknow error")
}

// 返回 response结果的data字段
//
func (api *BTC_com_api) doquery(url string) (interface{}, error) {
	api.Logger.Debug("do query", url)
	api.QueryCnt++
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		api.Logger.Debug(err)
		return nil, err;
	}
	re, err := api.Client.Do(req)
	if err != nil {
		api.Logger.Debug(err)
		return nil, err;
	}
	body, err := ioutil.ReadAll(re.Body)
	if err != nil {
		api.Logger.Debug(err)
		return nil, err
	}
	var resp Response

	err = json.Unmarshal(body, &resp)
	if err != nil {
		api.Logger.Debug(err,"\n", string(body),"\n", resp)
		return nil, err
	}
	if resp.Err_no != 0 || resp.Error_msg != "" {
		api.Logger.Debug("bad api call:", resp.Err_no,resp.Error_msg)
		return nil, errors.New("api errno != 0")
	}
	return resp.Data, nil
}

func getBlock(i interface{}) (*Block, error) {
	js, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	var b Block
	err = json.Unmarshal(js, &b)
	if err != nil {
		return nil, err
	}
	return &b, err
}
func getAddress(i interface{}) (*Address, error) {
	js, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	var a Address
	err = json.Unmarshal(js, &a)
	if err != nil {
		return nil, err
	}
	return &a, err
}
func getTransaction(i interface{}) (*Transaction, error) {
	js, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	var t Transaction
	err = json.Unmarshal(js, &t)
	if err != nil {
		return nil, err
	}
	return &t, err
}
