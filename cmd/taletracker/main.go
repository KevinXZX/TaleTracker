package main

import "taletracker.com/internal/http"

func main() {
	t := http.TaleServer{Config: &http.TaleConfig{
		DevelopmentMode: false,
		AllowedOrigins:  []string{"localhost:3000"},
	}}
	t.Start()
}
