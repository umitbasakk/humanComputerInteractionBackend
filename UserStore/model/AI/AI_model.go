package model

type AIRequest struct {
	StartedDate   string `json:"startedDate" xml:"startedDate" form:"startedDate" query:"startedDate"`
	EndDate       string `json:"endDate" xml:"endDate" form:"endDate" query:"endDate"`
	HashTag       string `json:"hashTag" xml:"hashTag" form:"hashTag" query:"hashTag"`
	Category      int64  `json:"category" xml:"category" form:"category" query:"category"`
	QuantityLimit int64  `json:"quantityLimit" xml:"quantityLimit" form:"quantityLimit" query:"quantityLimit"`
}
