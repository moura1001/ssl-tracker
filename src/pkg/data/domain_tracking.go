package data

import "time"

const StatusHealthy = "healthy"
const StatusExpires = "expired"
const StatusInvalid = "invalid"
const StatusOffline = "offline"
const StatusUnresponsive = "unresponsive"

type DomainTrackingInfo struct {
	Issuer        string
	SignatureAlgo string
	PublicKeyAlgo string
	EncodedPEM    string
	PublicKey     string
	Signature     string
	DNSNames      string
	KeyUsage      string
	ExtKeyUsages  []string
	Expires       time.Time
	Status        string
	LastPollAt    time.Time
	Latency       int
	Error         string
}

type DomainTracking struct {
	Id         int64
	UserId     string
	DomainName string

	DomainTrackingInfo
}

type TrackingAndAccount struct {
	NotifyUpfront int

	*DomainTracking
}

// TODO: implementation
func GetAllTrackingsWithAccount() ([]TrackingAndAccount, error) {
	return []TrackingAndAccount{
		{DomainTracking: &DomainTracking{DomainName: "google.com"}},
		{DomainTracking: &DomainTracking{DomainName: "facebook.com"}},
		{DomainTracking: &DomainTracking{DomainName: "youtube.com"}},
		{DomainTracking: &DomainTracking{DomainName: "twitter.com"}},
		{DomainTracking: &DomainTracking{DomainName: "amazon.com"}},
	}, nil
}

// TODO: implementation
func UpdateAllTrackings(trackings []DomainTracking) error {
	return nil
}
