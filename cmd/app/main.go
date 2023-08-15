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
	"github.com/huynchu/degree-planner-api/internal/auth"
	"github.com/huynchu/degree-planner-api/internal/course"
	"github.com/huynchu/degree-planner-api/internal/degree"
	degreecsv "github.com/huynchu/degree-planner-api/internal/degree-csv"
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

	// Run course data worker
	// courseDataWorker := workers.NewCourseDataWorker(db)
	// go courseDataWorker.Run()

	// create chi router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// add course routes
	courseStorage := course.NewCourseStorage(db)
	courseService := course.NewCourseService(courseStorage)
	courseController := course.NewCourseController(courseService)
	course.AddCourseRoutes(r, courseController)

	// add degree routes
	degreeStorage := degree.NewDegreeStorage(db)
	degreeService := degree.NewDegreeService(degreeStorage, courseService)
	degreeController := degree.NewDegreeController(degreeService)
	degree.AddDegreeRoutes(r, degreeController)

	// add degree csv routes
	degreeCsvStorage := degreecsv.NewDegreeCsvStorage("degree-csv", storage.NewS3FileStorage(s3Client))
	degreeCsvController := degreecsv.NewDegreeCsvController(degreeCsvStorage)
	r.Post("/degree-csv", degreeCsvController.UploadDegreeCsv)

	// auth0 endpoints
	r.Group(func(r chi.Router) {
		r.Use(mymiddleware.EnsureValidToken())
		r.Get("/api/private", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("private"))
		})
	})

	// Google login endpoint
	auth.InitializeOAuthGoogle()
	r.Get("/api/auth/login/google", auth.HandleGoogleLogin)
	r.HandleFunc("/api/auth/google/callback", auth.CallBackFromGoogle)

	// start server
	fmt.Println("starting server on port 8080...")
	http.ListenAndServe(":8080", r)
}
