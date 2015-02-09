package redisio

import (
	"log"
	"strings"
	"testing"
	"time"

	. "github.com/adjust/gocheck"
	"github.com/adjust/redis"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct {
	redisClient *redis.Client
}

var _ = Suite(&TestSuite{})

func (suite *TestSuite) SetUpSuite(c *C) {
	suite.redisClient = redis.NewTCPClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   9,
	})
	c.Assert(suite.redisClient.Ping().Err(), IsNil)
}

func (suite *TestSuite) SetUpTest(c *C) {
	suite.redisClient.FlushDb()
	suite.redisClient.Ping()
}

func (suite *TestSuite) TestRedisLogWriter(c *C) {
	writer, err := NewWriter(suite.redisClient, "test123")
	c.Assert(err, IsNil)

	writer.Write([]byte("dingdong"))
	time.Sleep(time.Millisecond)

	value, err := suite.redisClient.LPop("test123").Result()
	c.Assert(err, IsNil)

	c.Check(value, Equals, "dingdong")
}

func (suite *TestSuite) TestRedisLogWriterForStdLog(c *C) {
	writer, err := NewWriter(suite.redisClient, "test123")
	c.Assert(err, IsNil)

	log.SetOutput(writer)
	log.Println("hello world")
	time.Sleep(time.Millisecond)

	value, err := suite.redisClient.LPop("test123").Result()
	c.Assert(err, IsNil)

	c.Check(strings.Contains(value, "hello world"), Equals, true)
}

func (suite *TestSuite) TestOrder(c *C) {
	writer, err := NewWriter(suite.redisClient, "test123")
	c.Assert(err, IsNil)

	writer.Write([]byte("1"))
	writer.Write([]byte("2"))
	writer.Write([]byte("3"))
	writer.Write([]byte("4"))
	writer.Write([]byte("5"))

	time.Sleep(time.Millisecond)

	value, err := suite.redisClient.LPop("test123").Result()
	c.Assert(err, IsNil)
	c.Check(value, Equals, "1")

	value, err = suite.redisClient.LPop("test123").Result()
	c.Assert(err, IsNil)
	c.Check(value, Equals, "2")

	value, err = suite.redisClient.LPop("test123").Result()
	c.Assert(err, IsNil)
	c.Check(value, Equals, "3")

	value, err = suite.redisClient.LPop("test123").Result()
	c.Assert(err, IsNil)
	c.Check(value, Equals, "4")

	value, err = suite.redisClient.LPop("test123").Result()
	c.Assert(err, IsNil)
	c.Check(value, Equals, "5")
}
