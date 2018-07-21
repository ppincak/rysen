package server

type Routes struct {
	symbols    string
	config     string
	statistics string
	live       string
}

func newRoutesV1() *Routes {
	return &Routes{
		symbols:    "/v1/symbols",
		config:     "/v1/config",
		statistics: "/v1/statistics",
		live:       "v1/live",
	}
}

var RoutesV1 *Routes = newRoutesV1()
