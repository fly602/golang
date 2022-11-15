package model

type ProductSize struct {
	ID       int64  `grom:"primary_key;not_null;auto_increment" json:"id"`
	SizeName string `json:"size_name"`
	SizeCode string `json:"size_code"`
}
