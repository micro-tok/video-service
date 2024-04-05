package main
grpc
import (
	"net"

	"github.com/micro-tok/video-service/pkg/config"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	//TODO: init cassandra conn

	lis, err := net.Listen("tcp", "localhost:"+cfg.ServicePort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()

	//TODO: deploy grpcServer

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

}
