package data

import "time"

type User struct {
	Id          string
	Email       string
	Password    string
	AccessToken string
	ExpiresAt   time.Time
}
