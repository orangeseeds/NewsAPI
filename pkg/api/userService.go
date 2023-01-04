package api

type UserRepository interface {
	CreateUser(AuthUserRequest) error
	FindUserByEmailAndPassword(string, string) (*User, error)
	SetUserReadArticle(string, string) (*Article, error)
	SetUserBookmarksArticle(string, string) (*Article, error)
}

type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		storage: userRepo,
	}
}

func (u *userService) Create(user AuthUserRequest) error {
	err := u.storage.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) Login(email string, password string) (*User, error) {
	user, err := u.storage.FindUserByEmailAndPassword(email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) Read(email string, articleId string) (*Article, error) {
	article, err := u.storage.SetUserReadArticle(email, articleId)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (u *userService) Bookmark(email string, articleId string) (*Article, error) {
	article, err := u.storage.SetUserBookmarksArticle(email, articleId)
	if err != nil {
		return nil, err
	}
	return article, nil
}
