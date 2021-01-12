package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Broker []string `json:"broker"`
	Topics []Topic  `json:"topics"`
	Groups string   `json:"groups"`
}

type Topic struct {
	Topic  string   `json:"topic"`
	Prefix []string `json:"prefix"`
}

func contain(data []string, target string) bool {
	for _, i := range data {
		if i == target {
			return true
		}
	}
	return false
}

func getTopic(prefix string, listTopic []Topic) string {
	for _, i := range listTopic {
		if contain(i.Prefix, prefix) {
			return i.Topic
		}
	}
	return "Topic tidak ditemukan"
}

// Config for Apache Kafka
func configTransaction(iso string) (broker string, group string, topic string) {
	log.Printf("Get config for current request: %s\n", iso)

	file, _ := os.Open("./config.json")
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	json.Unmarshal(b, &config)

	prefix := iso[4:8]
	selectedTopic := getTopic(prefix, config.Topics)
	log.Println("Broker: ", config.Broker[0])
	log.Println("Group: ", config.Groups)
	log.Println("Topic: ", selectedTopic)

	return config.Broker[0], config.Groups, selectedTopic
}

var (
	broker = "localhost:9092"
	group  = "test-go"
	topic1 = "channel-kafka"
	topic2 = "kafka-biller"
)
