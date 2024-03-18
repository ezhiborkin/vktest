package models

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required" example:"ivanov@mail.ru""`
	Password string `json:"password" binding:"required" example:"123456ksksksksk"`
}

type UserCreate struct {
	Email    string `json:"email" binding:"required" example:"ivanov@mail.ru""`
	Role     string `json:"role" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"123456ksksksksk"`
}
