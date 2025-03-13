package authschema

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginCommand struct {
	Login     string
	Password  string
	IpAdderss string
}
