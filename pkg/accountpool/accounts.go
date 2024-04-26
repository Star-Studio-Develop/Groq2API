package accountpool

import (
	groq "github.com/learnLi/groq_client"
	"sync"
)

type IAccounts struct {
	Accounts []*groq.Account `json:"accounts"`
	mx       sync.Mutex
}

func (a *IAccounts) Get() *groq.Account {
	a.mx.Lock()
	defer a.mx.Unlock()
	if len(a.Accounts) == 0 {
		return nil
	}
	account := a.Accounts[0]
	a.Accounts = append(a.Accounts[1:], account)
	return account
}

func NewAccounts(accounts []*groq.Account) *IAccounts {
	return &IAccounts{Accounts: accounts}
}
