package model

import "time"

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type File struct {
	ID           uint      `gorm:"primaryKey"`
	FileName     string    `gorm:"column:file_name"`
	Extension    string    `gorm:"column:extension"`
	Size         int       `gorm:"column:size"`
	CreationDate time.Time `gorm:"column:creation_date"` // Thêm trường ngày tạo
}

type Filter struct {
	Field    string `json:"field"`              // Trường cần lọc
	Operator string `json:"operator,omitempty"` // Toán tử so sánh
	Value    any    `json:"value"`              // Giá trị lọc
	Logic    string `json:"logic"`              // AND hoặc OR
}

type FilterRequest struct {
	Filters []Filter `json:"filters"`
}
