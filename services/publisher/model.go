package publisher

type Model struct {
	KafkaTopic string `json:"kafkaTopic"`
	ReadTopic  string `json:"readTopic"`
}
