package workers

import (
	"testing"

	"github.com/huynchu/degree-planner-api/internal/course"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCourseCompare(t *testing.T) {
	course1 := course.CourseDB{
		Name: "Introduction to Algorithms",
		Code: "CSCI-2300",
		Prerequisites: [][]string{
			{"CSCI-1200"},
			{"CSCI-2200", "MATH-2800"},
			{"MATH-1010", "MATH-1500", "MATH-1020", "MATH-2010"},
		},
		Corequisites:  []string{},
		CrossListings: []string{},
	}

	course1BSON := bson.M{
		"name":          "Introduction to Algorithms",
		"code":          "CSCI-2300",
		"prerequisites": [][]string{{"CSCI-1200"}, {"CSCI-2200", "MATH-2800"}, {"MATH-1010", "MATH-1500", "MATH-1020", "MATH-2010"}},
		"corequisites":  []string{},
		"crossListings": []string{},
	}

	course1BSONEncoded, err := bson.Marshal(course1BSON)
	if err != nil {
		t.Error(err)
	}

	var course1BSONDecoded course.CourseDB
	err = bson.Unmarshal(course1BSONEncoded, &course1BSONDecoded)
	if err != nil {
		t.Error(err)
	}

	if !course1.Equal(&course1BSONDecoded) {
		t.Error("course1 != course1BSONDecoded")
	}
}
