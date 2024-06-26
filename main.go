package main

import (
	"net"

	_ "github.com/joho/godotenv/autoload"
	"github.com/micro-tok/video-service/pkg/cassandra"
	"github.com/micro-tok/video-service/pkg/config"
	"github.com/micro-tok/video-service/pkg/pb"
	"github.com/micro-tok/video-service/pkg/redis"
	"github.com/micro-tok/video-service/pkg/s3"
	"github.com/micro-tok/video-service/pkg/services"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	cassService := cassandra.NewCassandraService(cfg)

	s3Service := s3.NewAWSService(cfg)

	redisService := redis.NewRedisService(cfg)

	lis, err := net.Listen("tcp", "localhost:"+cfg.ServicePort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024 * 1024 * 1024 * 2),
	)

	// Register the service
	pb.RegisterVideoServiceServer(grpcServer, services.NewVideoService(cassService, s3Service, redisService))

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
