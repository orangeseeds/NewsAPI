package api

type AuthUserRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	UserId             string `json:"user_id,omitempty"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	Password           string `json:"password,omitempty"`
	PasswordResetToken string `json:"password_reset_token,omitempty"`
}

type Article struct {
	ArticleId     string   `json:"article_id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	PublishedDate string   `json:"published_date"`
	URLtoMedia    []string `json:"url_to_media"`
	URltoArticle  string   `json:"url_to_article"`
}

type SaveArticleRequest struct {
	ArticleId string `json:"article_id"`
}

type UserService interface {
	Create(user AuthUserRequest) error
	Login(email string, password string) (*User, error)
	Read(email string, articleId string) (*Article, error)
	Bookmark(email string, articleId string) (*Article, error)
}

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
