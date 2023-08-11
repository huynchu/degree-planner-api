package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/huynchu/degree-planner-api/config"
	"github.com/huynchu/degree-planner-api/internal/course"
	"github.com/huynchu/degree-planner-api/internal/degree-csv"
	mymiddleware "github.com/huynchu/degree-planner-api/internal/middleware"
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
	db, err := storage.BootstrapMongo(env.MONGODB_URI, env.MONGODB_NAME, 10*time.Second)
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}
	fmt.Println("connected to mongodb...")

	// create s3 client
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}
	s3Client := s3.NewFromConfig(cfg)

	// create chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// add course routes
	courseStorage := course.NewCourseStorage(db)
	courseController := course.NewCourseController(courseStorage)
	course.AddCourseRoutes(r, courseController)

	// add degree csv routes
	degreeCsvStorage := degree.NewDegreeCsvStorage("degree-csv", storage.NewS3FileStorage(s3Client))
	degreeCsvController := degree.NewDegreeCsvController(degreeCsvStorage)
	r.Post("/degree-csv", degreeCsvController.UploadDegreeCsv)

	r.Group(func(r chi.Router) {
		r.Use(mymiddleware.EnsureValidToken)
		r.Get("/api/private", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("private"))
		})
	})

	// start server
	fmt.Println("starting server...")
	http.ListenAndServe(":8080", r)
}
