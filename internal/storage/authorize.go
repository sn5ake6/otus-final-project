package storage

type Authorize struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

func NewAuthorize(login string, password string, ip string) Authorize {
	return Authorize{
		Login:    login,
		Password: password,
		IP:       ip,
	}
}
