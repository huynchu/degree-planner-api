package degree

import "github.com/go-chi/chi/v5"

func AddDegreeRoutes(r chi.Router, controller *DegreeController) {
	r.Post("/api/degrees", controller.CreateDegree)
	r.Get("/api/degrees/{degreeID}", controller.FindDegreeByID)
	r.Post("/api/degrees/{degreeID}/semesters", controller.AddSemester)
	r.Put("/api/degrees/{degreeID}/semesters/{index}/move", controller.MoveSemester)
	r.Post("/api/degrees/{degreeID}/semesters/{index}/courses", controller.AddCourse)
}
