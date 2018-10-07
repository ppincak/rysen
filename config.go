package main

type Config struct {
	MaxProcs      int      `json:"maxProcs"`
	ServerPort    int      `json:"serverPort"`
	DecryptionKey string   `json:"decryptionKey"`
	KafkaBrokers  []string `json:"kafkaBrokers"`
}
