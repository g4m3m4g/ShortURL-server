package models

import "gorm.io/gorm"

type SimpleUrl struct {
	gorm.Model
	OriginalUrl string `gorm:"unique"`
	ShortUrl    string `gorm:"unique"`
}
