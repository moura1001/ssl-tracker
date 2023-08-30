package data

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
