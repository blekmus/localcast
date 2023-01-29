package models

import "html/template"

type Episode struct {
	ID              int `gorm:"column:id"`
	PodcastId       int
	Title           string
	Source          string `gorm:"column:link"`
	Date            int    `gorm:"column:published"`
	DateString      string
	FileSize        int
	MimeType        string
	FileName        string `gorm:"column:download_filename"`
	Duration        int    `gorm:"column:total_time"`
	Description     string `gorm:"column:description_html"`
	DescriptionHtml template.HTML
}

func (Episode) TableName() string {
	return "episode"
}
