package main

import (
	"time"

	"github.com/huynchu/degree-planner-api/config"
	"github.com/huynchu/degree-planner-api/internal/storage"
	"github.com/huynchu/degree-planner-api/internal/workers"
)

func main() {
	// load config
	env, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	// connect to db
	db, err := storage.BootstrapMongo(env.MONGODB_URI, env.MONGODB_NAME, 10*time.Second)
	if err != nil {
		panic(err)
	}

	courseDataWorker := workers.NewCourseDataWorker(db)
	courseDataWorker.Run()
}
