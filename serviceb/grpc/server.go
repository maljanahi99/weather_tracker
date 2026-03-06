package grpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	pb "github.com/maljanahi99/weather_tracker/protos"
)

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

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
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("POSTGRES_USER", "dev_user")
	dbPassword := getEnvOrDefault("POSTGRES_PASSWORD", "dev_password")
	dbName := getEnvOrDefault("POSTGRES_DB", "dev_db")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
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
		VALUES ($1, $2, $3)
	`
	if _, err := db.ExecContext(ctx, sqlStatement, req.Latitude, req.Longitude, data.Current.Temperature); err != nil {
		return nil, fmt.Errorf("failed to insert weather report: %w", err)
	}

	return &pb.WeatherResponse{
		Current: &pb.Current{
			Temperature: data.Current.Temperature,
		},
	}, nil
}
