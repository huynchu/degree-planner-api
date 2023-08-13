package course

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseController struct {
	courseService *CourseService
}

func NewCourseController(csrv *CourseService) *CourseController {
	return &CourseController{
		courseService: csrv,
	}
}

func (cc *CourseController) FindCourseByID(w http.ResponseWriter, r *http.Request) {
	// extract url params
	courseID := chi.URLParam(r, "courseID")

	// fetch course from db
	course, err := cc.courseService.FindCourseByID(courseID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "course not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "database fetch error: fetch course", http.StatusInternalServerError)
		return
	}

	// Respond with json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(course)
}

func (cc *CourseController) SearchCourse(w http.ResponseWriter, r *http.Request) {
	// extract query params
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "missing query param", http.StatusBadRequest)
		return
	}
	maxLimit := 10
	limit := 5
	if limitQuery := r.URL.Query().Get("limit"); limitQuery != "" {
		tmp, err := strconv.Atoi(limitQuery)
		if err != nil || tmp <= 0 || tmp > maxLimit {
			http.Error(w, fmt.Sprintf("invalid limit param: limit must be greater than 0 and then than %v", maxLimit), http.StatusBadRequest)
			return
		}
		limit = tmp
	}

	// seach course
	courses, err := cc.courseService.SearchCourse(query, limit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "course not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "database fetch error: fetch course", http.StatusInternalServerError)
		return
	}

	// Respond with json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(courses)
}
