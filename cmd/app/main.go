package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/huynchu/degree-planner-api/config"
	"github.com/huynchu/degree-planner-api/internal/storage"
)

func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// load config
	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}
	// connect to db
	fmt.Println("connecting to mongodb...")
	_, err = storage.BootstrapMongo(env.MONGODB_URI, env.MONGODB_NAME, 10*time.Second)
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// create chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// start server
	http.ListenAndServe(":8080", r)
}
