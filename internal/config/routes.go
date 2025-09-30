package routes

import "net/http"
import c "backend/internal/controller"

func CreateRoutes() {
	http.HandleFunc("/", c.IndexHelloWorld)
}
