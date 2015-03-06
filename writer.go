package redisio

import (
	"github.com/adjust/redis"
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
	var todo int
	var batch []string

	for logLine := range writer.inputChannel {
		batch = append(batch, logLine)

		// are we still on a batch run?
		if todo > 0 {
			todo--
			continue
		}

		// batch run done, flushing
		writer.pushToRedis(batch)
		batch = []string{}

		// fetch next batch run size
		todo = len(writer.inputChannel) - 1 // otherwise we'd lose the last line
	}
}

func (writer *Writer) pushToRedis(logLines []string) {
	err := writer.redisClient.RPush(writer.listName, logLines...).Err()
	if err != nil {
		panic(err)
	}
}
