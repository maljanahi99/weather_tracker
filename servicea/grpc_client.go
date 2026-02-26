package servicea

import (
	pb "github.com/maljanahi99/weather_tracker/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewServiceBClient(addr string) (pb.WeatherServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	return pb.NewWeatherServiceClient(conn), conn, nil
}
