package course

import "github.com/go-chi/chi/v5"

func AddCourseRoutes(r chi.Router, controller *CourseController) {

	courseRouter := chi.NewRouter()

	courseRouter.Get("/courses/{id}", controller.FindCourseByID)
	courseRouter.Post("/courses", controller.CreateCourse)

	r.Mount("/api", courseRouter)
}
