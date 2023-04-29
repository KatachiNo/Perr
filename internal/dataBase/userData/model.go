package userData

type UserData struct { //Users Table
	Id          int    `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone-number"`
	Country     string `json:"country"`
	City        string `json:"city"`
	Index       string `json:"index"`
	Street      string `json:"street"`
	NumberHouse string `json:"number-house"`
	Note        string `json:"note"`
	FirstName   string `json:"FirstName"`
	MiddleName  string `json:"middle-name"`
	LastName    string `json:"last-name"`
}
