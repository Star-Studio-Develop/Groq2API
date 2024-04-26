package global

import (
	"github.com/patrickmn/go-cache"
	"groqai2api/pkg/accountpool"
	"groqai2api/pkg/proxypool"
)

var (
	Cache         *cache.Cache
	Host          string
	Port          string
	ChinaPrompt   string
	ProxyPool     *proxypool.IProxy
	Authorization string
	AccountPool   *accountpool.IAccounts
)
