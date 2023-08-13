package degree

import (
	"errors"

	"github.com/huynchu/degree-planner-api/internal/course"
	"github.com/huynchu/degree-planner-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrSemesterIndexOutOfBounds = errors.New("semester index out of bounds")
)

type DegreeService struct {
	degreeStorage *DegreeStorage

	courseService *course.CourseService
}

func NewDegreeService(ds *DegreeStorage, cs *course.CourseService) *DegreeService {
	return &DegreeService{
		degreeStorage: ds,
		courseService: cs,
	}
}

func (ds *DegreeService) CreateDegree(name string) (string, error) {
	return ds.degreeStorage.CreateDegree(name)
}

func (ds *DegreeService) FindDegreeByID(id string) (*DegreeDB, error) {
	return ds.degreeStorage.FindDegreeByID(id)
}

func (ds *DegreeService) AddSemester(degreeID string, semesterName string, semesterIndex int) error {
	degree, err := ds.degreeStorage.FindDegreeByID(degreeID)
	if err != nil {
		return err
	}

	if semesterIndex < 0 || semesterIndex > len(degree.Semesters) {
		return ErrSemesterIndexOutOfBounds
	}

	semester := Semester{
		Name:    semesterName,
		Courses: []primitive.ObjectID{},
	}

	degree.Semesters = utils.Insert(degree.Semesters, semesterIndex, semester)

	err = ds.degreeStorage.UpdateSemesters(degreeID, degree.Semesters)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DegreeService) MoveSemester(degreeID string, semesterIndex int, newIndex int) error {
	degree, err := ds.degreeStorage.FindDegreeByID(degreeID)
	if err != nil {
		return err
	}

	if semesterIndex < 0 || semesterIndex >= len(degree.Semesters) ||
		newIndex < 0 || newIndex >= len(degree.Semesters) {
		return ErrSemesterIndexOutOfBounds
	}

	if semesterIndex == newIndex {
		return nil
	}

	degree.Semesters = utils.Move(degree.Semesters, semesterIndex, newIndex)

	err = ds.degreeStorage.UpdateSemesters(degreeID, degree.Semesters)
	return err
}

func (ds *DegreeService) AddCourseToSemester(degreeID string, semesterIndex int, courseID string) error {
	// Check if degree exists
	degree, err := ds.degreeStorage.FindDegreeByID(degreeID)
	if err != nil {
		return err
	}

	// Check if semester exists
	if semesterIndex < 0 || semesterIndex >= len(degree.Semesters) {
		return ErrSemesterIndexOutOfBounds
	}

	// Check if course exists
	course, err := ds.courseService.FindCourseByID(courseID)
	if err != nil {
		return err
	}

	// Add course to semester
	degree.Semesters[semesterIndex].Courses = append(degree.Semesters[semesterIndex].Courses, course.ID)

	// Update degree
	err = ds.degreeStorage.UpdateSemesters(degreeID, degree.Semesters)
	return err
}
