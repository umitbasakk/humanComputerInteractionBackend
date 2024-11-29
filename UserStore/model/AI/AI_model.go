package model

import "time"

type AIRequest struct {
	StartedDate   string `json:"startedDate" xml:"startedDate" form:"startedDate" query:"startedDate"`
	EndDate       string `json:"endDate" xml:"endDate" form:"endDate" query:"endDate"`
	HashTag       string `json:"hashTag" xml:"hashTag" form:"hashTag" query:"hashTag"`
	Category      string `json:"category" xml:"category" form:"category" query:"category"`
	QuantityLimit string `json:"quantityLimit" xml:"quantityLimit" form:"quantityLimit" query:"quantityLimit"`
}

type AIData struct {
	UserId        string    `json:"userId" xml:"userId" form:"userId" query:"userId"`
	StartedDate   string    `json:"startedDate" xml:"startedDate" form:"startedDate" query:"startedDate"`
	EndDate       string    `json:"endDate" xml:"endDate" form:"endDate" query:"endDate"`
	HashTag       string    `json:"hashTag" xml:"hashTag" form:"hashTag" query:"hashTag"`
	Category      int       `json:"category" xml:"category" form:"category" query:"category"`
	QuantityLimit int       `json:"quantityLimit" xml:"quantityLimit" form:"quantityLimit" query:"quantityLimit"`
	RequestStatus int       `json:"requestStatus" xml:"requestStatus" form:"requestStatus" query:"requestStatus"`
	Created_at    time.Time `json:"created_at" xml:"created_at" form:"created_at" query:"created_at"`
	Updated_at    time.Time `json:"updated_at" xml:"updated_at" form:"updated_at" query:"updated_at"`
}
