package config

import "os"

type Config struct {
	BucketKey          string
	BucketSecret       string
	BucketName         string
	ServicePort        string
	CassandraClusterIP string
	CassandraKeyspace  string
}

func LoadConfig() *Config {
	return &Config{
		BucketKey:          os.Getenv("BUCKET_KEY"),
		BucketSecret:       os.Getenv("BUCKET_SECRET"),
		BucketName:         os.Getenv("BUCKET_NAME"),
		ServicePort:        os.Getenv("SERVICE_PORT"),
		CassandraClusterIP: os.Getenv("CASSANDRA_CLUSTER_IP"),
		CassandraKeyspace:  os.Getenv("CASSANDRA_KEYSPACE"),
	}
}
