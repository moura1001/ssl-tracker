package db_account

import (
	"context"

	db_service "github.com/moura1001/ssl-tracker/src/pkg/db/service"

	"github.com/moura1001/ssl-tracker/src/pkg/data"
	"github.com/moura1001/ssl-tracker/src/pkg/logger"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
	"github.com/uptrace/bun"
)

type AccountBunStore struct{}

func NewAccountBunStore() AccountBunStore {
	return AccountBunStore{}
}

func (abs AccountBunStore) GetUserAccount(userId string) (*data.Account, error) {
	account := new(data.Account)
	ctx := context.Background()
	err := db_service.Bun.NewSelect().
		Model(account).
		Where("user_id = ?", userId).
		Scan(ctx)
	return account, err
}

func (abs AccountBunStore) GetAccount(query util.Map) (*data.Account, error) {
	account := new(data.Account)
	builder := db_service.Bun.NewSelect().Model(account)
	for k, v := range query {
		if v != "" {
			builder.Where("? = ?", bun.Ident(k), v)
		}
	}
	err := builder.Scan(context.Background())
	return account, err
}

func (abs AccountBunStore) UpdateAccount(account *data.Account) error {
	_, err := db_service.Bun.NewUpdate().
		Model(account).
		WherePK().
		Exec(context.Background())
	return err
}

func (abs AccountBunStore) CreateAccountForUserIfNotExist(user *data.User) (*data.Account, error) {
	if acc, err := abs.GetUserAccount(user.Id); err == nil {
		return acc, nil
	}

	acc := data.Account{
		UserId:             user.Id,
		NotifyUpfront:      7,
		DefaultNotifyEmail: user.Email,
		Plan:               data.PlanFree,
	}

	_, err := db_service.Bun.NewInsert().Model(&acc).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	logger.Log("event", "new account signup", "id", acc.Id)
	return &acc, nil
}
