package server

type Routes struct {
	// Crypto symbols
	symbols string
	// Schemas resource
	schema string
	// System statistics
	statistics string
	// Live feed
	live string
	// Feeds resource
	feeds string
	// Subscribtion to specific feed
	subscribeToFeed string
}

func newRoutesV1() *Routes {
	return &Routes{
		symbols:         "/v1//symbols/:exchange",
		schema:          "/v1/schemas",
		statistics:      "/v1/statistics",
		live:            "/v1/live",
		feeds:           "/v1/feeds",
		subscribeToFeed: "/v1/feeds/subscribe",
	}
}

var RoutesV1 *Routes = newRoutesV1()
