package domain

type RegisterRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Tokens TokenPair `json:"tokens"`
	User   User      `json:"user"`
}

type LogoutRequest struct {
	UserID string `json:"user_id"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ConfirmationRequest struct {
	Email            string `json:"email"`
	ConfirmationCode string `json:"confirmation_code"`
}
