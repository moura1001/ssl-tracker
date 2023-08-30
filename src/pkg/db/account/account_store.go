package db_account

import (
	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

type AccountStore interface {
	GetUserAccount(userId string) (*data.Account, error)
	GetAccount(query util.Map) (*data.Account, error)
	UpdateAccount(account *data.Account) error
	CreateAccountForUserIfNotExist(user *data.User) (*data.Account, error)
}
