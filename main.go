package main

import (
	"github.com/micro-tok/video-service/pkg/config"
)

func main() {
	cfg := config.LoadConfig()

	//TODO: init cassandra conn

	//TODO: deploy grpcServer

}
