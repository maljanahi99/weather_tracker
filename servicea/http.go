package servicea

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	pb "github.com/maljanahi99/weather_tracker/protos"
)

type HTTPServer struct {
	ServiceB pb.WeatherServiceClient
}

func (s *HTTPServer) Weather(router *mux.Router) {
	router.
		HandleFunc("/weather", s.weatherHandler).
		Methods("GET").
		Name("weather")
}

func (s *HTTPServer) weatherHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	resp, err := s.ServiceB.GetWeather(ctx, &pb.GetWeatherRequest{
		Latitude:  52.52,
		Longitude: 13.41,
	})

	if err != nil {
		http.Error(w, "grpc call failed: "+err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func StartServer(serviceB pb.WeatherServiceClient) error {
	router := mux.NewRouter()

	s := &HTTPServer{ServiceB: serviceB}
	s.Weather(router)

	return http.ListenAndServe(":8080", router)
}
