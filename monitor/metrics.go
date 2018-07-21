package monitor

type Metrics interface {
	Map() map[string]interface{}
}
