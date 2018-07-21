package monitor

type Reporter interface {
	Statistics() []*Statistic
}
