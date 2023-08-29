package data

import (
	"context"

	"github.com/moura1001/ssl-tracker/src/pkg/db"
	"github.com/moura1001/ssl-tracker/src/pkg/logger"
	"github.com/moura1001/ssl-tracker/src/pkg/util"
	"github.com/uptrace/bun"
)

const (
	PlanFree     = "FREE"
	PlanStarter  = "STARTER"
	PlanBusiness = "BUSINESS"
)

type TrackingAndAccount struct {
	NotifyUpfront int

	DomainTracking
}

type Account struct {
	Id                 int64 `bun:"id,pk,autoincrement"`
	UserId             string
	Email              string
	SubscriptionStatus string
	Plan               string
	NotifyUpfront      int
	DefaultNotifyEmail string
}

func GetUserAccount(userId string) (*Account, error) {
	account := new(Account)
	ctx := context.Background()
	err := db.Bun.NewSelect().
		Model(account).
		Where("user_id = ?", userId).
		Scan(ctx)
	return account, err
}

func GetAccount(query util.Map) (*Account, error) {
	account := new(Account)
	builder := db.Bun.NewSelect().Model(account)
	for k, v := range query {
		if v != "" {
			builder.Where("? = ?", bun.Ident(k), v)
		}
	}
	err := builder.Scan(context.Background())
	return account, err
}

func UpdateAccount(account *Account) error {
	_, err := db.Bun.NewUpdate().
		Model(account).
		WherePK().
		Exec(context.Background())
	return err
}

func CreateAccountForUserIfNotExist(user *User) (*Account, error) {
	if acc, err := GetUserAccount(user.Id); err == nil {
		return acc, nil
	}

	acc := Account{
		UserId:             user.Id,
		NotifyUpfront:      7,
		DefaultNotifyEmail: user.Email,
		Plan:               PlanFree,
	}

	_, err := db.Bun.NewInsert().Model(&acc).Exec(context.Background())
	if err != nil {
		return nil, err
	}

	logger.Log("event", "new account signup", "id", acc.Id)
	return &acc, nil
}
