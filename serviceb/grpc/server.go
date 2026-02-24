package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	pb "github.com/maljanahi99/weather_tracker/protos"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "dev_user"
	password = "dev_password"
	dbname   = "dev_db"
)

type Server struct {
	pb.UnimplementedWeatherServiceServer
}

type openMeteoResponse struct {
	Current struct {
		Temperature float32 `json:"temperature_2m"`
	} `json:"current"`
}

func (s *Server) GetWeather(ctx context.Context, req *pb.GetWeatherRequest) (*pb.WeatherResponse, error) {
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&current=temperature_2m",
		req.Latitude, req.Longitude,
	)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	var data openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse weather JSON: %w", err)
	}

	// db connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Open database connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// query: postgres insert query takes (latitude, longitude, temperature)
	sqlStatement := `
       INSERT INTO weather_reports (latitude, longitude, temperature)
       VALUES (52.52, 13.41, -3)`
	// use the connection to exec the query

	_, err = db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}

	return &pb.WeatherResponse{
		Current: &pb.Current{
			Temperature: data.Current.Temperature,
		},
	}, nil
}
