package course

type CourseService struct {
	courseStorage *CourseStorage
}

func NewCourseService(cs *CourseStorage) *CourseService {
	return &CourseService{
		courseStorage: cs,
	}
}

func (cs *CourseService) FindCourseByID(id string) (*CourseDB, error) {
	return cs.courseStorage.FindCourseByID(id)
}

func (cs *CourseService) SearchCourse(query string, limit int) ([]CourseDB, error) {
	return cs.courseStorage.FindCourseByNameOrCode(query, limit)
}
