package main
import (
"net/http"
"fmt"
"urlshorner/controller"
"urlshorner/models"
"github.com/prometheus/client_golang/prometheus/promhttp"

)



func main() {


	store := models.NewStore()
	handler := controller.NewController(store)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("POST /generate", handler.HandleGenerateShortURL)
	mux.HandleFunc("GET /{code}", handler.HandleGetFullURL)
	mux.Handle("GET /metrics", promhttp.Handler())



	fmt.Println("listening on :8080")
	fmt.Println(http.ListenAndServe(":8000", handler.CounterMiddleware(mux)))
}