package degree

import (
	"errors"

	"github.com/huynchu/degree-planner-api/internal/course"
	"github.com/huynchu/degree-planner-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrSemesterIndexOutOfBounds      = errors.New("semester index out of bounds")
	ErrCourseAlreadyExistsInSemester = errors.New("course already exists in semester")
	ErrCourseDoesNotExistInSemester  = errors.New("course does not exist in semester")
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

func (ds *DegreeService) FindDegreeByID(id string) (*DegreeAggregated, error) {
	degree, err := ds.degreeStorage.FindDegreeByID(id)
	if err != nil {
		return nil, err
	}

	degreeAggregated := DegreeAggregated{
		ID:        degree.ID,
		Name:      degree.Name,
		Semesters: []SemesterAggregated{},
		Owner:     degree.Owner,
	}

	for _, semester := range degree.Semesters {
		semesterAggregated := SemesterAggregated{
			Name:    semester.Name,
			Courses: []course.CourseDB{},
		}
		for _, courseID := range semester.Courses {
			course, err := ds.courseService.FindCourseByID(courseID.Hex())
			if err != nil {
				return nil, err
			}
			semesterAggregated.Courses = append(semesterAggregated.Courses, *course)
		}
		degreeAggregated.Semesters = append(degreeAggregated.Semesters, semesterAggregated)
	}

	return &degreeAggregated, nil
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

func (ds *DegreeService) DeleteSemester(degreeID string, semesterIndex int) error {
	degree, err := ds.degreeStorage.FindDegreeByID(degreeID)
	if err != nil {
		return err
	}

	if semesterIndex < 0 || semesterIndex >= len(degree.Semesters) {
		return ErrSemesterIndexOutOfBounds
	}

	degree.Semesters = utils.Remove(degree.Semesters, semesterIndex)

	err = ds.degreeStorage.UpdateSemesters(degreeID, degree.Semesters)
	return err
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

	// Check if course already exists in semester
	for _, courseID := range degree.Semesters[semesterIndex].Courses {
		if courseID == course.ID {
			return ErrCourseAlreadyExistsInSemester
		}
	}

	// Add course to semester
	degree.Semesters[semesterIndex].Courses = append(degree.Semesters[semesterIndex].Courses, course.ID)

	// Update degree
	err = ds.degreeStorage.UpdateSemesters(degreeID, degree.Semesters)
	return err
}

func (ds *DegreeService) RemoveCourseFromSemester(degreeID string, semesterIndex int, courseID string) error {
	// Check if degree exists
	degree, err := ds.degreeStorage.FindDegreeByID(degreeID)
	if err != nil {
		return err
	}

	// Check if semester exists
	if semesterIndex < 0 || semesterIndex >= len(degree.Semesters) {
		return ErrSemesterIndexOutOfBounds
	}

	// Remove course from semester
	semester := &degree.Semesters[semesterIndex]
	courseIndex := -1
	for i, course := range semester.Courses {
		if course.Hex() == courseID {
			courseIndex = i
			break
		}
	}
	if courseIndex == -1 {
		return ErrCourseDoesNotExistInSemester
	}
	semester.Courses = utils.Remove(semester.Courses, courseIndex)

	// Update degree
	err = ds.degreeStorage.UpdateSemesters(degreeID, degree.Semesters)
	return err
}
