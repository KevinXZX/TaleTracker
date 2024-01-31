package main

import (
	"fmt"
	"taletracker.com/internal/http"
	"taletracker.com/internal/taledb"
)

func main() {
	fmt.Println("Hello, world!")
	// TODO: Enable SSL, this is for testing only
	database, err := taledb.OpenDatabase("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("error opening database: %v\n", err))
	}
	fmt.Println("Database opened successfully")
	t := http.TaleServer{
		Config: &http.TaleConfig{
			DevelopmentMode: false,
			AllowedOrigins:  []string{"localhost:3000"},
		},
		Db: database,
	}
	err = t.Start()
	if err != nil {
		fmt.Printf("error starting server: %v\n", err)
	}
	fmt.Println("Goodbye, world!")
}
