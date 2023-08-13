package course

import "github.com/go-chi/chi/v5"

func AddCourseRoutes(r chi.Router, controller *CourseController) {

	r.Get("/api/courses/{courseID}", controller.FindCourseByID)
	r.Get("/api/courses/search/", controller.SearchCourse)
}
