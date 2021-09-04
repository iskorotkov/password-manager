package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	openapi "github.com/iskorotkov/password-manager/go"
	"github.com/iskorotkov/password-manager/internal/database/postgres"
	"github.com/iskorotkov/password-manager/internal/services"
)

func main() {
	log.Printf("Server started")

	db, err := postgres.New(os.Getenv("CONN_STRING"))
	if err != nil {
		log.Fatal(err)
	}

	passwordService := services.NewPasswordService(db)
	passwordController := openapi.NewDefaultApiController(passwordService)

	apiRouter := openapi.NewRouter(passwordController)
	apiRouter.HandleFunc("/api/v1/passwords", dummyHandler).Methods(http.MethodOptions)
	apiRouter.HandleFunc("/api/v1/passwords/{id}", dummyHandler).Methods(http.MethodOptions)
	apiRouter.Use(
		addCORSMiddleware(),
	)

	rootRouter := http.NewServeMux()
	rootRouter.Handle("/openapi/v1/openapi.yaml", http.StripPrefix("/openapi/v1/", http.FileServer(http.Dir("./api"))))
	rootRouter.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./static/swagger/"))))
	rootRouter.Handle("/", apiRouter)

	log.Fatal(http.ListenAndServe(":8080", rootRouter))
}

func dummyHandler(_ http.ResponseWriter, _ *http.Request) {
}

func addCORSMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodOptions || req.Header.Get("Sec-Fetch-Mode") == "cors" {
				w.Header().Set("Access-Control-Allow-Methods", "*")
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			next.ServeHTTP(w, req)
		})
	}
}
