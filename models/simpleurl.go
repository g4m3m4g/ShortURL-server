package models

import "gorm.io/gorm"

type SimpleUrl struct {
	gorm.Model
	OriginalUrl string `gorm:"type:varchar(600);unique"`
	ShortUrl    string `gorm:"unique"`
}
