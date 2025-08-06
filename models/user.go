package models

type User struct {
	Name     string `json:"name" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"required,gte=18,lte=120"`
	Password string `json:"password" binding:"required,password"`
}
