package course

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseController struct {
	storage *CourseStorage
}

func NewCourseController(storage *CourseStorage) *CourseController {
	return &CourseController{
		storage: storage,
	}
}

type CreateCourseRequest struct {
	Name          string     `json:"name"`
	Code          string     `json:"code"`
	Prerequisites [][]string `json:"prerequisites"`
	Corequisites  []string   `json:"corequisites"`
}

func (s *CourseController) CreateCourse(w http.ResponseWriter, r *http.Request) {
	// decode json body
	var createCourseReq CreateCourseRequest
	err := json.NewDecoder(r.Body).Decode(&createCourseReq)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	// create course
	course := &CourseDB{
		Name:          createCourseReq.Name,
		Code:          createCourseReq.Code,
		Prerequisites: createCourseReq.Prerequisites,
		Corequisites:  createCourseReq.Corequisites,
	}

	// insert course into db
	insertedID, err := s.storage.CreateCourse(course)
	if err != nil {
		http.Error(w, "database insert error: create course", http.StatusInternalServerError)
		return
	}

	// encode json response
	res := struct {
		ID string `json:"id"`
	}{
		ID: insertedID,
	}

	// Respond with json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (s *CourseController) FindCourseByID(w http.ResponseWriter, r *http.Request) {
	// extract url params
	courseID := chi.URLParam(r, "courseID")

	// fetch course from db
	course, err := s.storage.FindCourseByID(courseID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "course not found", http.StatusNotFound)
			return
		}
		http.Error(w, "database fetch error: fetch course", http.StatusInternalServerError)
		return
	}

	// Respond with json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}
