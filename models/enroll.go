package models

type EnrollForm struct {
	Members    []byte `json:"members"`
	Contest    string `json:"contest"`
	CreateTime string `json:"create_time"`
	School     string `json:"school"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

type EnrollInformation struct {
	Members    []byte    `json:"members"`
	Contest    string    `json:"contest"`
	CreateTime OftenTime `json:"create_time"`
	School     string    `json:"school"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Deleted    OftenTime `json:"deleted"`
}
