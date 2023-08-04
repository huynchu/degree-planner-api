package myapp

type Course struct {
	ID            string
	Name          string
	Code          string
	Prerequisites []string
	Corequisites  []string
}

type CourseService interface {
	// Find a course by ID
	FindCourseByID(id string) (*Course, error)

	// Find courses by filter
	FindCourses(filter CourseFilter) ([]*Course, error)

	// Create a course
	CreateCourse(course *Course) error

	// Update a course
	UpdateCourse(id string, update CourseUpdate) (*Course, error)

	// Delete a course
	DeleteCourse(id string) error
}

type CourseFilter struct {
	Name string
}

type CourseUpdate struct {
	Name          string
	Code          string
	Prerequisites []string
	Corequisites  []string
}
