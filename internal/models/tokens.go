package models

type TokenPayload struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	RoleID string `json:"roleID"`
}

type Tokens struct {
	AcceptToken  string `json:"acceptToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshToken struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}
