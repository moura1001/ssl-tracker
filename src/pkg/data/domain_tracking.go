package data

import (
	"time"
)

type DomainTrackingInfo struct {
	Issuer        string
	SignatureAlgo string `bun:"-"`
	PublicKeyAlgo string `bun:"-"`
	EncodedPEM    string `bun:"-"`
	PublicKey     string `bun:"-"`
	Signature     string
	DNSNames      string
	KeyUsage      string
	ExtKeyUsages  []string `bun:",array"`
	Expires       time.Time
	Status        string
	LastPollAt    time.Time
	Latency       int
	Error         string `bun:"-"`
}

type DomainTracking struct {
	Id         int64  `bun:"id,pk,autoincrement"`
	UserId     string `bun:",pk"`
	DomainName string `bun:",pk"`

	DomainTrackingInfo
}
