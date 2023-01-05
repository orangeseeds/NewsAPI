package api

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	// Username string `json:"username"`
	Password string `json:"password,omitempty" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type ReadArticleRequest struct {
	ArticleId string `json:"article_id" validate:"required"`
}
type BookmarkArticleRequest struct {
	ArticleId string `json:"article_id" validate:"required"`
}

type FollowSourceRequest struct {
	SourceId string `json:"source_id" validate:"required"`
}
type BlockSourceRequest struct {
	SourceId string `json:"source_id" validate:"required"`
}

type User struct {
	UserId             string `json:"user_id,omitempty"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	Password           string `json:"password,omitempty"`
	PasswordResetToken string `json:"password_reset_token,omitempty"`
}

type Article struct {
	ArticleId     string   `json:"article_id,omitempty"`
	Title         string   `json:"title,omitempty"`
	Content       string   `json:"content,omitempty"`
	PublishedDate string   `json:"published_date,omitempty"`
	URLtoMedia    []string `json:"url_to_media,omitempty"`
	URltoArticle  string   `json:"url_to_article,omitempty"`
}

type Source struct {
	SourceId string `json:"source_id,omitempty"`
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

type Creator struct {
	CreatorId string `json:"creator_id,omitempty"`
	Name      string `json:"name,omitempty"`
}

type Category struct {
	CategoryId string `json:"category_id,omitempty"`
	Name       string `json:"name,omitempty"`
	ImageURL   string `json:"image_url,omitempty"`
}
