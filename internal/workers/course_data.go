package workers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/huynchu/degree-planner-api/config"
	"github.com/huynchu/degree-planner-api/internal/course"
	"go.mongodb.org/mongo-driver/mongo"
)

type CourseDataWorker struct {
	db *mongo.Database
}

func NewCourseDataWorker(db *mongo.Database) *CourseDataWorker {
	return &CourseDataWorker{
		db: db,
	}
}

func (w *CourseDataWorker) Run() {
	courseData := make(map[string]*course.CourseDB)
	populateCourseData(courseData)

	for _, c := range courseData {
		fmt.Println(c)
	}
}

type courseJson struct {
	Csre string `json:"csre"` // i.e 1200
	Name string `json:"name"` // i.e Data Structures
	Sbj  string `json:"subj"` // i.e CSCI
	// Desc string `json:"description"`
}

type coursePrerequisiteJson struct {
	// Atributes     []string     `json:"attributes"`
	Corequisites  []string     `json:"corequisites"`
	CrossListings []string     `json:"cross_listings"`
	Prerequisites Prerequisite `json:"prerequisites"`
}

type Prerequisite struct {
	Course string         `json:"course"`
	Type   string         `json:"type"`
	Nested []Prerequisite `json:"nested,omitempty"`
}

func populateCourseData(courseData map[string]*course.CourseDB) error {
	// Get course catalog data
	catalogData, err := fetchCourseCatalogData()
	if err != nil {
		fmt.Println("Error fetching course catalog data:", err)
		return err
	}

	// Decode course catalog data
	courseDataMap := make(map[string]courseJson)
	err = json.Unmarshal(catalogData, &courseDataMap)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return err
	}

	// Get course prereq data
	prereqData, err := fetchCoursePrereqData()
	if err != nil {
		fmt.Println("Error fetching course prereq data:", err)
		return err
	}

	// Decode course prereq data
	coursePrereqDataMap := make(map[string]coursePrerequisiteJson)
	err = json.Unmarshal(prereqData, &coursePrereqDataMap)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return err
	}

	// populate courseData with course catalog data
	for key, c := range courseDataMap {
		newDBCourse := &course.CourseDB{
			Code:          key,
			Name:          c.Name,
			Prerequisites: [][]string{},
			Corequisites:  []string{},
			CrossListings: []string{},
		}
		courseData[key] = newDBCourse
	}

	// populate courseData with course prereq data
	for key, cprq := range coursePrereqDataMap {
		c, ok := courseData[key]
		if ok {
			c.Corequisites = cprq.Corequisites
			c.Prerequisites = cprq.Prerequisites.TransformPrereq()
			c.CrossListings = cprq.CrossListings
		} else {
			if cprq.CrossListings != nil || len(cprq.CrossListings) > 0 {
				for _, crossListing := range cprq.CrossListings {
					crossListedCourse, ok := courseData[crossListing]
					if ok {
						course := course.CourseDB{
							Code:          key,
							Name:          crossListedCourse.Name,
							Prerequisites: cprq.Prerequisites.TransformPrereq(),
							Corequisites:  cprq.Corequisites,
							CrossListings: cprq.CrossListings,
						}
						courseData[crossListing] = &course
					}
					// else: cross listing also dont exist
				}
			}
			// else: course not found and no cross listing
		}
	}

	// no errors
	return nil
}

func fetchCourseCatalogData() ([]byte, error) {
	env, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filepath := filepath.Join(rootDir, "internal", "data", "catalog.json")

	switch env.GO_ENV {
	case "dev":
		return fetchJsonDataFromLocalFile(filepath)
	case "prod":
		return fetchJsonDataFromGitHub(env.COURSE_DATA_URL)
	default:
		return fetchJsonDataFromLocalFile(filepath)
	}
}

func fetchCoursePrereqData() ([]byte, error) {
	env, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	rootDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filepath := filepath.Join(rootDir, "internal", "data", "prereq_data.json")

	switch env.GO_ENV {
	case "dev":
		return fetchJsonDataFromLocalFile(filepath)
	case "prod":
		return fetchJsonDataFromGitHub(env.COURSE_PREREQ_DATA_URL)
	default:
		return fetchJsonDataFromLocalFile(filepath)
	}
}

func fetchJsonDataFromLocalFile(filePath string) ([]byte, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	dataBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil
}

func fetchJsonDataFromGitHub(url string) ([]byte, error) {
	// Make the GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p Prerequisite) TransformPrereqRecursive(res *[][]string) []string {
	if p.Type == "and" {
		for _, prereq := range p.Nested {
			*res = append(*res, prereq.TransformPrereqRecursive(res))
		}
		return []string{}
	} else if p.Type == "or" {
		or := []string{}
		for _, prereq := range p.Nested {
			or = append(or, prereq.TransformPrereqRecursive(res)...)
		}
		return or
	} else {
		return []string{p.Course}
	}
}

func (p Prerequisite) TransformPrereq() [][]string {
	if p.Type == "course" {
		return [][]string{{p.Course}}
	} else {
		res := [][]string{}
		p.TransformPrereqRecursive(&res)
		return res
	}
}
