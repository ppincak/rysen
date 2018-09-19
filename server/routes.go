package server

type Routes struct {
	symbols         string
	schema          string
	statistics      string
	live            string
	feeds           string
	subscribeToFeed string
}

func newRoutesV1() *Routes {
	return &Routes{
		symbols:         "/v1/symbols",
		schema:          "/v1/schemas",
		statistics:      "/v1/statistics",
		live:            "v1/live",
		feeds:           "v1/feeds",
		subscribeToFeed: "v1/feeds/subscribe",
	}
}

var RoutesV1 *Routes = newRoutesV1()
