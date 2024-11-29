package model

import "time"

type User struct {
	Id         int
	Name       string    `json:"name" xml:"name" form:"name" query:"name"`
	Username   string    `json:"username" xml:"username" form:"username" query:"username"`
	Password   string    `json:"password" xml:"password" form:"password" query:"password"`
	Email      string    `json:"email" xml:"email" form:"email" query:"email"`
	Phone      string    `json:"phone" xml:"phone" form:"phone" query:"phone"`
	Token      string    `json:"token" xml:"token" form:"token" query:"token"`
	Created_at time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	Updated_at time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
}

type UserResponse struct {
	Name       string    `json:"name" xml:"name" form:"name" query:"name"`
	Username   string    `json:"username" xml:"username" form:"username" query:"username"`
	Email      string    `json:"email" xml:"email" form:"email" query:"email"`
	Phone      string    `json:"phone" xml:"phone" form:"phone" query:"phone"`
	Token      string    `json:"token" xml:"token" form:"token" query:"token"`
	Created_at time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	Updated_at time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
}

type PasswordRequest struct {
	CurrentPassword string `json:"currentPassword" xml:"currentPassword" form:"currentPassword" query:"currentPassword"`
	NewPassword     string `json:"newPassword" xml:"newPassword" form:"newPassword" query:"newPassword"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name" xml:"name" form:"name" query:"name"`
	Username string `json:"username" xml:"username" form:"username" query:"username"`
	Email    string `json:"email" xml:"email" form:"email" query:"email"`
}

type ResponseUser struct {
	Token string `json:"token" xml:"token" form:"token" query:"token"`
}
