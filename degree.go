package myapp

type Semester struct {
	ID      string
	Name    string
	Type    string
	Year    int
	Courses []string
}

type Degree struct {
	ID        string
	Name      string
	Semesters []Semester
}
