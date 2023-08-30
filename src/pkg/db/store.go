package db

import (
	db_account "github.com/moura1001/ssl-tracker/src/pkg/db/account"
	db_domain "github.com/moura1001/ssl-tracker/src/pkg/db/domain"
	db_service "github.com/moura1001/ssl-tracker/src/pkg/db/service"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
)

type store struct {
	Account db_account.AccountStore
	Domain  db_domain.DomainStore
}

var Store *store

func Init() {

	if Store != nil {
		return
	}

	Store = new(store)

	if util.GetEnv("BUN_MODE", "off") == "on" {
		db_service.Init()

		Store.Account = db_account.NewAccountBunStore()
		Store.Domain = db_domain.NewDomainBunStore()
	} else {
		Store.Account = db_account.NewAccountInMemoryStore()
		Store.Domain = db_domain.NewDomainInMemoryStore()
	}
}
