package domain

type User struct {
	ID     int
	UserID string
	Name   string
}

func NewUser(id int, userID, name string) *User {
	return &User{
		ID:     id,
		UserID: userID,
		Name:   name,
	}
}
