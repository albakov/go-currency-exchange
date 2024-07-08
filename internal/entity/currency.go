package entity

type Currency struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	FullName string `json:"name"`
	Sign     string `json:"sign"`
}
