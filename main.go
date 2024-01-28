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
	err := t.Start()
	if err != nil {
		fmt.Printf("error starting server: %v\n", err)
	}
	fmt.Println("Goodbye, world!")
}
