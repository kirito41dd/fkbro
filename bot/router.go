package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"sync"
)

type Router interface {
	DoHandle(url string, update *tgbotapi.Update)
	AddHandle(url string, handleFunc HandleFunc)
}

type HandleFunc func (*tgbotapi.Update)

type BotRouter struct {
	lock sync.Mutex
	routMap map[string]HandleFunc
}

func NewBotRouter() *BotRouter {
	return &BotRouter{
		lock:    sync.Mutex{},
		routMap: make(map[string]HandleFunc),
	}
}


func (r *BotRouter) DoHandle(url string, update *tgbotapi.Update){
	r.lock.Lock()
	defer r.lock.Unlock()
	f, ok := r.routMap[url]
	if ok {
		go f(update)
	} else {
		for k, v := range r.routMap {
			if strings.HasPrefix(url, k) {
				go v(update)
				break
			}
		}
	}
}
func (r *BotRouter) AddHandle(url string, handleFunc HandleFunc){
	r.lock.Lock()
	defer r.lock.Unlock()
	_, ok := r.routMap[url]
	if !ok {
		r.routMap[url] = handleFunc
	}
}



