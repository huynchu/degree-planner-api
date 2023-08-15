package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrDatabaseFetchUser  = errors.New("error fetching user from database")
	ErrDatabaseCreateUser = errors.New("error creating user in database")
)

type UserService struct {
	userStorage *UserStorage
}

func NewUserService(userStorage *UserStorage) *UserService {
	return &UserService{userStorage: userStorage}
}

func (s *UserService) FindUserByEmail(email string) (*UserDB, error) {
	return s.userStorage.FindUserByEmail(email)
}

func (s *UserService) CreateNewUser(email string) (string, error) {
	return s.userStorage.CreateNewUser(email)
}
