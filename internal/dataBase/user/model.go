package user

type User struct { //Users Table
	Id                 string `json:"id"`
	Login              string `json:"login"`
	CategoryOfUser     string `json:"category-of-user"`
	DateOfRegistration string `json:"date-of-registration"`
	Password           string `json:"password"`
	PasswordHash       string `json:"-"`
	Salt               string `json:"-"`
	Algorithm          string `json:"-"`
}

type Token struct {
	TokenString string `json:"token"`
}
