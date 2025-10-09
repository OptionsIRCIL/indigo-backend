package main

import (
	c "backend/internal/config"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Printf("Hello, World!\n")
	c.CreateRoutes()

	fmt.Printf("Serving on :8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
