package degree

import "github.com/go-chi/chi/v5"

func AddDegreeRoutes(r chi.Router, controller *DegreeController) {
	// Degree routes
	r.Post("/api/degrees", controller.CreateDegree)
	r.Get("/api/degrees/{degreeID}", controller.FindDegreeByID)

	// Degree Semesters routes
	r.Post("/api/degrees/{degreeID}/semesters", controller.AddSemester)
	r.Put("/api/degrees/{degreeID}/semesters/{index}/move", controller.MoveSemester)
	r.Delete("/api/degrees/{degreeID}/semesters/{index}", controller.DeleteSemester)

	// Degree Semester Courses routes
	r.Post("/api/degrees/{degreeID}/semesters/{index}/courses", controller.AddCourseToSemester)
	r.Delete("/api/degrees/{degreeID}/semesters/{index}/courses/{courseID}", controller.RemoveCourseFromSemester)
}
