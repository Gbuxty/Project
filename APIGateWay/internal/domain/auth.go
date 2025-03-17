package domain

type User struct {
    ID    string `json:"id"`
    Email string `json:"email"`
}

type LogoutRequest struct {
    UserID string `json:"user_id"`
}

type AuthRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

type ConfirmationRequest struct {
    Email            string `json:"email"`
    ConfirmationCode string `json:"confirmation_code"`
}