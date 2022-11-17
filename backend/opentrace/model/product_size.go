package model

type ProductSize struct {
	ID       int64  `gorm:"primary_key;not_null;auto_increment" json:"id"`
	SizeName string `json:"size_name"`
	SizeCode string `json:"size_code"`
}
