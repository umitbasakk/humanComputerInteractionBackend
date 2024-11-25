package model

import "time"

type Verify struct {
	Id           int
	Username     string    `json:"username" xml:"username" form:"username" query:"username"`
	VerifyCode   string    `json:"verifyCode" xml:"verifyCode" form:"verifyCode" query:"verifyCode"`
	VerifyStatus int       `json:"verifyStatus" xml:"verifyStatus" form:"verifyStatus" query:"verifyStatus"`
	Created_at   time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	Updated_at   time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
}

type VerifyRequest struct {
	VerifyCode string    `json:"verifyCode" xml:"verifyCode" form:"verifyCode" query:"verifyCode"`
	Created_at time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	Updated_at time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
}
