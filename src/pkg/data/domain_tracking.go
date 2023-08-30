package data

import (
	"time"
)

type DomainTrackingInfo struct {
	Issuer        string
	SignatureAlgo string
	PublicKeyAlgo string
	EncodedPEM    string
	PublicKey     string
	Signature     string
	DNSNames      string
	KeyUsage      string
	ExtKeyUsages  []string `bun:",array"`
	Expires       time.Time
	Status        string
	LastPollAt    time.Time
	Latency       int
	Error         string
}

type DomainTracking struct {
	Id         int64 `bun:"id,pk,autoincrement"`
	UserId     string
	DomainName string

	DomainTrackingInfo
}
