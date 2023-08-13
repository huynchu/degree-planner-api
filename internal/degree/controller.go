package degree

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type DegreeController struct {
	degreeService *DegreeService
}

func NewDegreeController(dsrv *DegreeService) *DegreeController {
	return &DegreeController{
		degreeService: dsrv,
	}
}

type CreateDegreeRequest struct {
	Name string `json:"name"`
}

func (dc *DegreeController) CreateDegree(w http.ResponseWriter, r *http.Request) {
	// decode json body
	var createDegreeReq CreateDegreeRequest
	err := json.NewDecoder(r.Body).Decode(&createDegreeReq)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	// create degree
	id, err := dc.degreeService.CreateDegree(createDegreeReq.Name)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "database insert error: create degree", http.StatusInternalServerError)
		return
	}

	// encode json response
	res := struct {
		ID string `json:"id"`
	}{
		ID: id,
	}

	// Respond with json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (dc *DegreeController) FindDegreeByID(w http.ResponseWriter, r *http.Request) {
	// extract url params
	degreeID := chi.URLParam(r, "degreeID")

	// fetch course from db
	degree, err := dc.degreeService.FindDegreeByID(degreeID)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(degree)
}

type AddSemesterRequest struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

func (dc *DegreeController) AddSemester(w http.ResponseWriter, r *http.Request) {
	// extract url params
	degreeID := chi.URLParam(r, "degreeID")

	// decode json body
	var addSemesterReq AddSemesterRequest
	err := json.NewDecoder(r.Body).Decode(&addSemesterReq)
	if err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	// add semester
	err = dc.degreeService.AddSemester(degreeID, addSemesterReq.Name, addSemesterReq.Index)
	if err != nil {
		if err == ErrSemesterIndexOutOfBounds {
			http.Error(w, "semester index out of bounds", http.StatusBadRequest)
			return
		}
		fmt.Println(err)
		http.Error(w, "database insert error: add semester", http.StatusInternalServerError)
		return
	}

	// Respond with json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("added semester successfully")
}
