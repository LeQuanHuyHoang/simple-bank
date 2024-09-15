package repo

import (
	"Go_Learn/pkg/model"
	"Go_Learn/pkg/utils"
	"fmt"
	"gorm.io/gorm"
)

type Repo struct {
	Postgres *gorm.DB
}

func NewRepo(pg *gorm.DB) IRepo {
	repo := &Repo{
		Postgres: pg,
	}
	return repo
}

type IRepo interface {
	SearchFiles(filterRequest []model.Filter) ([]model.File, error)
	ApplyFilters(query *gorm.DB, filters []model.Filter) *gorm.DB
}

// Hàm tìm kiếm tập tin dựa trên bộ lọc
func (r *Repo) SearchFiles(filterRequest []model.Filter) ([]model.File, error) {
	var files []model.File
	query := r.Postgres.Model(&model.File{})

	// Áp dụng các bộ lọc
	query = r.ApplyFilters(query, filterRequest)

	// Thực hiện truy vấn
	fmt.Println(query)
	err := query.Find(&files).Error
	return files, err
}

func (r *Repo) ApplyFilters(query *gorm.DB, filters []model.Filter) *gorm.DB {
	fmt.Println(filters)
	for idx, filter := range filters {

		switch filter.Field {
		case "file_name":
			condition := fmt.Sprintf("%s LIKE  ?", filter.Field)
			if idx == 0 || filter.Logic == "AND" {
				query = query.Where(condition, "%"+filter.Value.(string)+"%")
			} else if filter.Logic == "OR" {
				query = query.Or(condition, "%"+filter.Value.(string)+"%")
			}
		case "extension":
			condition := fmt.Sprintf("%s =  ?", filter.Field)
			if idx == 0 || filter.Logic == "AND" {
				query = query.Where(condition, filter.Value)
			} else if filter.Logic == "OR" {
				query = query.Or(condition, filter.Value)
			}
		case "size":
			condition := fmt.Sprintf("%s %s ?", filter.Field, filter.Operator)
			if idx == 0 || filter.Logic == "AND" {
				query = query.Where(condition, filter.Value)
			} else if filter.Logic == "OR" {
				query = query.Or(condition, filter.Value)
			}
		case "creation_date":
			date, err := utils.ParseDate(filter.Value.(string))
			if err != nil {
				fmt.Println("Invalid date format:", err)
				continue
			}

			condition := fmt.Sprintf("%s %s ?", filter.Field, filter.Operator)
			if idx == 0 || filter.Logic == "AND" {
				query = query.Where(condition, date)
			} else if filter.Logic == "OR" {
				query = query.Or(condition, date)
			}
		}

	}

	return query
}

//type FileSearchCriteria struct {
//	FileName      string    `json:"file_name,omitempty"`
//	Extension     string    `json:"extension,omitempty"`
//	MinSize       int64     `json:"min_size,omitempty"`
//	MaxSize       int64     `json:"max_size,omitempty"`
//	CreatedAfter  string    `json:"created_after,omitempty"`
//	CreatedBefore string    `json:"created_before,omitempty"`
//	ModifiedAfter string    `json:"modified_after,omitempty"`
//	ModifiedBefore string   `json:"modified_before,omitempty"`
//	AccessedAfter string    `json:"accessed_after,omitempty"`
//	AccessedBefore string   `json:"accessed_before,omitempty"`
//	FileContent   string    `json:"file_content,omitempty"`
//	SearchMode    string    `json:"search_mode,omitempty"` // AND/OR
//}
//
//// Hàm xử lý tìm kiếm trong database
//func searchFilesInDB(criteria FileSearchCriteria) ([]File, error) {
//	var files []File
//	query := db.Model(&File{})
//
//	// Điều kiện tìm kiếm tên file
//	if criteria.FileName != "" {
//		query = query.Where("file_name LIKE ?", "%"+criteria.FileName+"%")
//	}
//
//	// Điều kiện tìm kiếm phần mở rộng
//	if criteria.Extension != "" {
//		query = query.Where("extension = ?", criteria.Extension)
//	}
//
//	// Điều kiện tìm kiếm kích thước
//	if criteria.MinSize > 0 {
//		query = query.Where("size >= ?", criteria.MinSize)
//	}
//	if criteria.MaxSize > 0 {
//		query = query.Where("size <= ?", criteria.MaxSize)
//	}
//
//	// Điều kiện tìm kiếm ngày tạo
//	if criteria.CreatedAfter != "" {
//		createdAfter, _ := time.Parse("2006-01-02", criteria.CreatedAfter)
//		query = query.Where("created_at >= ?", createdAfter)
//	}
//	if criteria.CreatedBefore != "" {
//		createdBefore, _ := time.Parse("2006-01-02", criteria.CreatedBefore)
//		query = query.Where("created_at <= ?", createdBefore)
//	}
//
//	// Điều kiện tìm kiếm ngày sửa đổi
//	if criteria.ModifiedAfter != "" {
//		modifiedAfter, _ := time.Parse("2006-01-02", criteria.ModifiedAfter)
//		query = query.Where("modified_at >= ?", modifiedAfter)
//	}
//	if criteria.ModifiedBefore != "" {
//		modifiedBefore, _ := time.Parse("2006-01-02", criteria.ModifiedBefore)
//		query = query.Where("modified_at <= ?", modifiedBefore)
//	}
//
//	// Điều kiện tìm kiếm nội dung tập tin
//	if criteria.FileContent != "" {
//		query = query.Where("content LIKE ?", "%"+criteria.FileContent+"%")
//	}
//
//	// Thực hiện truy vấn và trả về kết quả
//	err := query.Find(&files).Error
//	return files, err
//}
