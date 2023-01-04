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
