package entity

import "time"

type User struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	SurName   string    `db:"surname"`
	UserName  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	RefreshT  string    `db:"refresh_t"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}

type UserResponse struct {
	ID        string
	Name      string
	SurName   string
	UserName  string
	Email     string
	RefreshT  string
	AccessT   string
	Role      string
}
type UserListResponse struct {
	User      []User
	UserCount int
}
type DeleteRequest struct {
	ID        string
	DeletedAt time.Time
}
type UpdateRefresh struct {
	UserID       string
	RefreshToken string
}
type Response struct {
	Status bool
}
type Message struct {
	Message string
}
type UpdatePassword struct {
	UserID      string
	NewPassword string
}
type UserCreateRequst struct {
	Name     string
	Surname  string
	Username string
	Email    string
	Password string
}
type UserCreateResponse struct {
	ID string
}
type PaginatedResponse struct {
	Data       UserListResponse
	Page       int
	Limit      int
	TotalCount int
	TotalPages int
}
