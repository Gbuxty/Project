package domain

type CreatePostRequest struct {
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`
}

type Post struct {
	Content   string `json:"content"`
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
}