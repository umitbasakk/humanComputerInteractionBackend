package model

import "time"

type User struct {
	Id         int
	Name       string    `json:"name" xml:"name" form:"name" query:"name"`
	Username   string    `json:"username" xml:"username" form:"username" query:"username"`
	Password   string    `json:"password" xml:"password" form:"password" query:"password"`
	Email      string    `json:"email" xml:"email" form:"email" query:"email"`
	Token      string    `json:"token" xml:"token" form:"token" query:"token"`
	Created_at time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	Updated_at time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
}

type ResponseUser struct {
	Token string `json:"token" xml:"token" form:"token" query:"token"`
}

type MessageHandler struct {
	Message string
	Data    interface{}
}
