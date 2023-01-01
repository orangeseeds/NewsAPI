package api

type NewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	UserId             string `json:"user_id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	PasswordResetToken string `json:"password_reset_token"`
}

type UserService interface {
	Create(user NewUserRequest) error
	All() (interface{}, error)
}

type UserRepository interface {
	CreateUser(NewUserRequest) error
	AllUser() (interface{}, error)
}

type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		storage: userRepo,
	}
}

func (u *userService) Create(user NewUserRequest) error {
	err := u.storage.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) All() (interface{}, error) {
	users, err := u.storage.AllUser()
	if err != nil {
		return nil, err
	}
	return users, nil
}
