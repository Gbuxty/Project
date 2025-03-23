package domain

type CreatePostRequest struct {
	Content  string `json:"content"`
	ImageURL string `json:"image_url"`

}

type PostResponse struct {
	Content   string `json:"content"`
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
}

type AllPostResponse struct {
	Posts      []*PostResponse `json:"posts"`
	TotalPosts int             `json:"total_posts"`
}
