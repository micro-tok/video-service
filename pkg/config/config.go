package config

import (
	"os"
	"strconv"
)

type Config struct {
	BucketKey          string
	BucketSecret       string
	BucketName         string
	ServicePort        string
	CassandraClusterIP string
	CassandraKeyspace  string
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
}

func LoadConfig() *Config {

	RedisAddr := os.Getenv("REDIS_ADDR")
	RedisPassword := os.Getenv("REDIS_PASSWORD")
	RedisDBString := os.Getenv("REDIS_DB")

	if RedisAddr == "" {
		RedisAddr = "localhost:6379"
	}

	var RedisDB int

	if RedisDBString == "" {
		RedisDB = 0
	} else {
		RedisDB, _ = strconv.Atoi(RedisDBString)
	}

	return &Config{
		BucketKey:          os.Getenv("BUCKET_KEY"),
		BucketSecret:       os.Getenv("BUCKET_SECRET"),
		BucketName:         os.Getenv("BUCKET_NAME"),
		ServicePort:        os.Getenv("SERVICE_PORT"),
		CassandraClusterIP: os.Getenv("CASSANDRA_CLUSTER_IP"),
		CassandraKeyspace:  os.Getenv("CASSANDRA_KEYSPACE"),
		RedisAddr:          RedisAddr,
		RedisPassword:      RedisPassword,
		RedisDB:            RedisDB,
	}
}
