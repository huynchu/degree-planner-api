package workers

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/huynchu/degree-planner-api/config"
	"github.com/huynchu/degree-planner-api/internal/course"
	"go.mongodb.org/mongo-driver/bson"
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

	courseCollection := w.db.Collection("courses")

	// Bulk update(upsert) courses
	models := []mongo.WriteModel{}
	for _, c := range courseData {
		filter := bson.M{"code": c.Code}
		update := bson.M{"$set": bson.M{
			"name":          c.Name,
			"prerequisites": c.Prerequisites,
			"corequisites":  c.Corequisites,
			"crossListings": c.CrossListings,
		}}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
	}

	result, err := courseCollection.BulkWrite(context.Background(), models)
	if err != nil {
		fmt.Println("Error bulk writing course data:", err)
		return
	}

	fmt.Println("Number of courses:", len(courseData))
	fmt.Println("Bulk write result:")
	fmt.Println("Matched", result.MatchedCount, "documents")
	fmt.Println("Upserted", result.UpsertedCount, "documents")
	fmt.Println("Modified", result.ModifiedCount, "documents")
	fmt.Println("UpsertedIds", result.UpsertedIDs)
}

type courseJson struct {
	Csre string `json:"csre"` // i.e 1200
	Name string `json:"name"` // i.e Data Structures
	Sbj  string `json:"subj"` // i.e CSCI
	// Desc string `json:"description"`
}

type coursePrerequisiteJson struct {
	// Atributes     []string     `json:"attributes"`
	Corequisites  []string      `json:"corequisites"`
	CrossListings []string      `json:"cross_listings"`
	Prerequisites *Prerequisite `json:"prerequisites"`
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
			if cprq.Corequisites != nil {
				c.Corequisites = cprq.Corequisites
			}
			if cprq.Prerequisites != nil {
				c.Prerequisites = cprq.Prerequisites.TransformPrereq()
			}
			if cprq.CrossListings != nil {
				c.CrossListings = cprq.CrossListings
			}
		} else {
			hasCrossListings := cprq.CrossListings != nil && len(cprq.CrossListings) > 0
			if hasCrossListings {
				courseCrossListing := getCrossListings(key, coursePrereqDataMap)
				for _, crossListing := range courseCrossListing {
					_, ok := courseData[crossListing]
					if ok {
						course := course.CourseDB{
							Code:          key,
							Name:          key + " (Cross-listed Course)",
							Prerequisites: [][]string{},
							Corequisites:  []string{},
							CrossListings: []string{},
						}
						if cprq.Corequisites != nil {
							course.Corequisites = cprq.Corequisites
						}
						if cprq.Prerequisites != nil {
							course.Prerequisites = cprq.Prerequisites.TransformPrereq()
						}
						if cprq.CrossListings != nil {
							course.CrossListings = cprq.CrossListings
						}
						courseData[key] = &course
						break
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

func getCrossListings(ccode string, courseData map[string]coursePrerequisiteJson) []string {
	visited := make(map[string]bool)

	queue := list.New()
	queue.PushBack(ccode)
	for queue.Len() > 0 {
		e := queue.Front()
		queue.Remove(e)
		courseCode := e.Value.(string)
		visited[courseCode] = true
		cdata, ok := courseData[courseCode]
		if ok {
			for _, crossListing := range cdata.CrossListings {
				inQueue := false
				for e := queue.Front(); e != nil; e = e.Next() {
					if e.Value.(string) == crossListing {
						inQueue = true
						break
					}
				}
				if !visited[crossListing] && !inQueue {
					queue.PushBack(crossListing)
				}
			}
		}
	}

	crossListings := make([]string, 0, len(visited))
	for key := range visited {
		if key != ccode {
			crossListings = append(crossListings, key)
		}
	}
	return crossListings
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
		tmp := strings.Split(p.Course, " ")
		course := tmp[0] + "-" + tmp[1]
		return []string{course}
	}
}

func (p Prerequisite) TransformPrereq() [][]string {
	if p.Type == "course" {
		tmp := strings.Split(p.Course, " ")
		course := tmp[0] + "-" + tmp[1]
		return [][]string{{course}}
	} else {
		res := [][]string{}
		p.TransformPrereqRecursive(&res)
		return res
	}
}
