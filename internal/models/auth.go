package models

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Username     string `json:"username"`
	UserRoleID   int    `json:"userRoleID"`
}

type RegistrationData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
