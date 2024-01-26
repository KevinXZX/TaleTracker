package main

import (
	"fmt"
	"taletracker.com/internal/http"
)

func main() {
	fmt.Println("Hello, world!")
	t := http.TaleServer{Config: &http.TaleConfig{
		DevelopmentMode: false,
		AllowedOrigins:  []string{"localhost:3000"},
	}}
	t.Start()
}
