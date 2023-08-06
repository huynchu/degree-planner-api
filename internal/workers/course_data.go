package workers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/huynchu/degree-planner-api/config"
	"github.com/huynchu/degree-planner-api/internal/course"
)

func RunCourseDataCronWorker() {
}

type courseJson struct {
	Csre string `json:"csre"` // i.e 1200
	Name string `json:"name"` // i.e Data Structures
	Sbj  string `json:"subj"` // i.e CSCI
	// Desc string `json:"description"`
}

func fetchCourseData(url string, courseData map[string]*course.CourseDB) error {
	// Make the GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	defer response.Body.Close()

	// Read the response body
	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	// Convert response body to json
	var courseDataMap map[string]courseJson
	err = json.Unmarshal(data, &courseDataMap)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return err
	}

	// populate courseData
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
	return nil
}

type CoursePrerequisiteJson struct {
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

func fetchCoursePrereqData(url string, courseData map[string]*course.CourseDB) error {
	// Make the GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("error: %v", err)
		return err
	}
	defer response.Body.Close()

	// Read the response body
	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return err
	}

	coursePrereqDataMap := make(map[string]CoursePrerequisiteJson)
	err = json.Unmarshal(data, &coursePrereqDataMap)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
		return err
	}

	// populate courseData
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
	return nil
}

func FetchAllData() {
	// load env
	env, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	courseData := make(map[string]*course.CourseDB)
	fetchCourseData(env.COURSE_DATA_URL, courseData)
	fetchCoursePrereqData(env.COURSE_PREREQ_DATA_URL, courseData)

	for _, c := range courseData {
		fmt.Println(c)
	}
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
