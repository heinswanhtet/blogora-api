package types

type ContextKey string

type ContextData struct {
	UserId string
	Email  string
}

type GoogleOAuthResp struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Sub     string `json:"sub"`
	Error   string `json:"error"`
}
