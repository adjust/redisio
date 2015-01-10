package redisio

import (
	redis "github.com/adjust/redis-latest-head"
)

type Writer struct {
	redisClient  *redis.Client
	listName     string
	inputChannel chan string
}

func NewWriter(redisClient *redis.Client, listName string) (writer *Writer, err error) {

	writer = &Writer{
		redisClient: redisClient,
		listName:    listName,
	}

	err = writer.redisClient.Ping().Err()
	if err != nil {
		return nil, err
	}
	writer.inputChannel = make(chan string, 10000)
	go writer.startConsumer()
	return writer, nil
}

func (writer *Writer) Write(p []byte) (n int, err error) {
	writer.inputChannel <- string(p)
	return len(p), nil
}

func (writer *Writer) startConsumer() {
	for logLine := range writer.inputChannel {
		err := writer.redisClient.RPush(writer.listName, logLine).Err()
		if err != nil {
			panic(err)
		}
	}
}
