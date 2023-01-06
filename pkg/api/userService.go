package api

type UserService interface {
	Create(user CreateUserRequest) (string, error)
	Login(email string, password string) (*User, error)
	Read(userId string, articleId string) (*Article, error)
	Bookmark(userId string, articleId string) (*Article, error)
	UnBookmark(userId string, articleId string) (*Article, error)
	FollowSource(userId string, sourceId string) (*Source, error)
	UnFollowSource(userId string, sourceId string) (*Source, error)
	FollowCategory(userId string, categoryId string) (*Category, error)
	UnFollowCategory(userId string, categoryId string) (*Category, error)
}

type UserRepository interface {
	CreateUser(CreateUserRequest) (string, error)
	FindUserByEmailAndPassword(string, string) (*User, error)
	SetUserReadArticle(string, string) (*Article, error)
	SetUserBookmarksArticle(string, string) (*Article, error)
	DelUserBookmarksArticle(string, string) (*Article, error)
	SetUserFollowsSource(string, string) (*Source, error)
	DelUserFollowsSource(string, string) (*Source, error)
	SetUserFollowsCategory(string, string) (*Category, error)
	DelUserFollowsCategory(string, string) (*Category, error)
}

type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		storage: userRepo,
	}
}

func (u *userService) Create(user CreateUserRequest) (string, error) {
	userId, err := u.storage.CreateUser(user)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (u *userService) Login(email string, password string) (*User, error) {
	user, err := u.storage.FindUserByEmailAndPassword(email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) Read(userId string, articleId string) (*Article, error) {
	article, err := u.storage.SetUserReadArticle(userId, articleId)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (u *userService) Bookmark(userId string, articleId string) (*Article, error) {
	article, err := u.storage.SetUserBookmarksArticle(userId, articleId)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (u *userService) UnBookmark(userId string, articleId string) (*Article, error) {
	article, err := u.storage.DelUserBookmarksArticle(userId, articleId)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (u *userService) FollowSource(userId string, sourceId string) (*Source, error) {
	source, err := u.storage.SetUserFollowsSource(userId, sourceId)
	if err != nil {
		return nil, err
	}
	return source, nil
}

func (u *userService) UnFollowSource(userId string, sourceId string) (*Source, error) {
	source, err := u.storage.DelUserFollowsSource(userId, sourceId)
	if err != nil {
		return nil, err
	}
	return source, nil
}

func (u *userService) FollowCategory(userId string, categoryId string) (*Category, error) {
	category, err := u.storage.SetUserFollowsCategory(userId, categoryId)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (u *userService) UnFollowCategory(userId string, categoryId string) (*Category, error) {
	category, err := u.storage.DelUserFollowsCategory(userId, categoryId)
	if err != nil {
		return nil, err
	}
	return category, nil
}
