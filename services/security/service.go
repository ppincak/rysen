package security

type Service struct {
	accounts map[string][]*Account
}

func NewService() *Service {
	return &Service{
		accounts: make(map[string][]*Account),
	}
}

func (service *Service) Initialize(accounts []*Account) {
	for _, account := range accounts {
		service.AddAccount(account)
	}
}

func (service *Service) AddAccount(account *Account) {
	var accounts []*Account
	if list, ok := service.accounts[account.Exchange]; !ok {
		accounts = make([]*Account, 0)
	} else {
		accounts = list
	}
	accounts = append(accounts, account)
	service.accounts[account.Exchange] = accounts
}

func (service *Service) GetAccounts(exchange string) []*Account {
	if accounts, ok := service.accounts[exchange]; ok {
		return accounts
	}
	return make([]*Account, 0)
}
