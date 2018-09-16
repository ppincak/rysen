package server

type Routes struct {
	symbols           string
	schema            string
	createSchema      string
	statistics        string
	live              string
	feeds             string
	feedSubscriptions string
	subscribeToFeed   string
}

func newRoutesV1() *Routes {
	return &Routes{
		symbols:           "/v1/symbols",
		schema:            "/v1/schemas",
		createSchema:      "/v1/schemas/:name",
		statistics:        "/v1/statistics",
		live:              "v1/live",
		feeds:             "v1/feeds",
		feedSubscriptions: "v1/feeds/subscriptions",
		subscribeToFeed:   "v1/feeds/subscribe",
	}
}

var RoutesV1 *Routes = newRoutesV1()
