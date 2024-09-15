package utils

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

func IsErrNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func ParseDate(dateStr string) (time.Time, error) {
	// Parse ngày theo định dạng "02-01-2006" (định dạng cho DD-MM-YYYY)
	date, err := time.Parse("2-1-2006", dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
