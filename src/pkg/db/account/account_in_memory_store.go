package db_account

import (
	"fmt"

	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/logger"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

type AccountInMemoryStore struct {
	accounts []data.Account
	userId   int64
}

func NewAccountInMemoryStore() *AccountInMemoryStore {
	return &AccountInMemoryStore{
		accounts: []data.Account{},
	}
}

func (abs AccountInMemoryStore) GetUserAccount(userId string) (*data.Account, error) {
	for _, account := range abs.accounts {
		if account.UserId == userId {
			return &account, nil
		}
	}
	return nil, fmt.Errorf("user %s does not exist", userId)
}

func (abs AccountInMemoryStore) GetAccount(query util.Map) (*data.Account, error) {
	for _, account := range abs.accounts {
		isEqualsCount := 0
		queryEquals := util.Map{}
		for k, v := range query {
			if v != "" {
				switch k {
				case "email":
					isEquals := account.Email == v.(string)
					queryEquals[k] = isEquals
					if isEquals {
						isEqualsCount++
					}
				case "subscription_status":
					isEquals := account.SubscriptionStatus == v.(string)
					queryEquals[k] = isEquals
					if isEquals {
						isEqualsCount++
					}
				case "plan":
					isEquals := account.Plan == v.(string)
					queryEquals[k] = isEquals
					if isEquals {
						isEqualsCount++
					}
				case "notify_upfront":
					isEquals := account.NotifyUpfront == v.(int)
					queryEquals[k] = isEquals
					if isEquals {
						isEqualsCount++
					}
				case "default_notify_email":
					isEquals := account.DefaultNotifyEmail == v.(string)
					queryEquals[k] = isEquals
					if isEquals {
						isEqualsCount++
					}
				default:
					break
				}
			}
		}
		if len(queryEquals) == isEqualsCount {
			return &account, nil
		}
	}

	return nil, fmt.Errorf("no account was found for the query %v", query)
}

func (abs *AccountInMemoryStore) UpdateAccount(account *data.Account) error {
	for i, acc := range abs.accounts {
		if acc.Id == account.Id {
			abs.accounts[i] = *account
		}
	}
	return nil
}

func (abs *AccountInMemoryStore) CreateAccountForUserIfNotExist(user *data.User) (*data.Account, error) {
	var (
		userId = user.Id
	)

	if len(userId) > 20 {
		if acc, err := abs.GetUserAccount(userId); err == nil {
			return acc, nil
		}
	} else {
		userId = util.NewID()
	}

	acc := data.Account{
		Id:                 abs.userId,
		UserId:             userId,
		NotifyUpfront:      7,
		DefaultNotifyEmail: user.Email,
		Plan:               data.PlanBusiness,
	}

	abs.accounts = append(abs.accounts, acc)

	abs.userId = abs.userId + 1

	logger.Log("event", "new account signup", "id", acc.Id)
	return &acc, nil
}
