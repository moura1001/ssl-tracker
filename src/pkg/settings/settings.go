package settings

import "github.com/moura1001/ssl-tracker/src/pkg/data"

type accountSettings struct {
	MaxTrackings int
}

var Account = map[string]accountSettings{
	data.PlanFree: {
		MaxTrackings: 2,
	},
	data.PlanStarter: {
		MaxTrackings: 20,
	},
	data.PlanBusiness: {
		MaxTrackings: 200,
	},
}
