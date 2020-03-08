package btcinfo

import (
	"github.com/zshorz/fkbro/util"
	"testing"
)

func Test_getinfo(t *testing.T) {
	api := NewBTC_com_api(util.Config.ApiHost, util.Config.CacheSize)

	data, err := api.doquery(api.Host + "/block/50")
	if err != nil {
		api.Logger.Error(err)
	}
	b, err := getBlock(data)
	if err != nil {
		api.Logger.Error(err)
	}
	api.Logger.Info(b, "\n\n")

	data, err = api.doquery(api.Host + "/address/1HeJdFWF5751beh9edfDhVzjetk5AdiUFM")
	if err != nil {
		api.Logger.Error(err)
	}
	a, err := getAddress(data)
	if err != nil {
		api.Logger.Error(err)
	}
	api.Logger.Info(a, "\n\n")

	data, err = api.doquery(api.Host + "/tx/b85e710221b6c14c47c9503e7ff4b10109573b582c1e4c82644084c7adf17f93")
	if err != nil {
		api.Logger.Error(err)
	}
	trans, err := getTransaction(data)
	if err != nil {
		api.Logger.Error(err)
	}
	api.Logger.Info(trans, "\n\n")


}

func Test_API(t *testing.T) {

	api := NewBTC_com_api(util.Config.ApiHost, util.Config.CacheSize)

	r, err := api.GetBlockInfo("latest")
	if err != nil {
		api.Logger.Error(err)
	}
	api.Logger.Info(r)
	api.Logger.Info(r.GetReword())
	for i:=0; i<10; i++ {

		r, err := api.GetBlockInfo("50")
		if err != nil {
			api.Logger.Error(err)
		}
		api.Logger.Info(r)
	}

	//

	r2, err := api.GetAddressInfo("166WoAxVFTc4H9WLLec76HERvn3LKytjJ9")
	if err != nil {
		api.Logger.Error(err)
	}
	api.Logger.Info(r2)
	for i:=0; i<10; i++ {

		r2, err := api.GetAddressInfo("166WoAxVFTc4H9WLLec76HERvn3LKytjJ9")
		if err != nil {
			api.Logger.Error(err)
		}
		api.Logger.Info(r2)
	}
	//

	r3, err := api.GetTransactionInfo("9b5877539f6b422ae5c72e1a105defa4de5ecbc1d9175a0c30207d887b3864df")
	if err != nil {
		api.Logger.Error(err)
	}
	api.Logger.Info(r3)
	for i:=0; i<10; i++ {

		r3, err := api.GetTransactionInfo("9b5877539f6b422ae5c72e1a105defa4de5ecbc1d9175a0c30207d887b3864df")
		if err != nil {
			api.Logger.Error(err)
		}
		api.Logger.Info(r3)
	}

}
