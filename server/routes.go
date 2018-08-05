package server

type Routes struct {
	symbols           string
	config            string
	statistics        string
	live              string
	feedList          string
	feedSubscriptions string
	subscribeToFeed   string
}

func newRoutesV1() *Routes {
	return &Routes{
		symbols:           "/v1/symbols",
		config:            "/v1/config",
		statistics:        "/v1/statistics",
		live:              "v1/live",
		feedList:          "v1/feeds",
		feedSubscriptions: "v1/feeds/subscriptions",
		subscribeToFeed:   "v1/feeds/subscribe",
	}
}

var RoutesV1 *Routes = newRoutesV1()
