package user

type UserController struct {
	userService *UserService
}

func NewCourseController(usrv *UserService) *UserController {
	return &UserController{
		userService: usrv,
	}
}
