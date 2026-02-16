package grpc

import (
	"context"

	pb "github.com/maljanahi99/weather_tracker/protos"
)

type Server struct {
	pb.UnimplementedWeatherServiceServer
}

func (s *Server) GetWeather(ctx context.Context, req *pb.GetWeatherRequest) (*pb.WeatherResponse, error) {
	return &pb.WeatherResponse{
		Current: &pb.Current{
			Temperature: req.Latitude + req.Longitude,
		},
	}, nil
}
