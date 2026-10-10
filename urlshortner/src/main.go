package main
import (
"net/http"
"fmt"
"urlshorner/controller"
"urlshorner/models"
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
	fmt.Println("listening on :8080")
	fmt.Println(http.ListenAndServe(":8000", mux))
}