package http

type User struct {
	Id    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}

type ListUserResponseType struct {
	Users []*User `json:"users,omitempty"`
}