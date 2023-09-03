package data

import "time"

type User struct {
	Id          string
	Email       string
	AccessToken string
	ExpiresAt   time.Time
}
