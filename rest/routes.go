package rest

type Routes struct {
	symbols    string
	config     string
	statistics string
}

func newRoutesV1() *Routes {
	return &Routes{
		symbols:    "/symbols",
		config:     "/config",
		statistics: "/statistics",
	}
}

var RoutesV1 *Routes = newRoutesV1()
