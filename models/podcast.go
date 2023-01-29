package models

type Podcast struct {
	ID             int
	Title          string
	Description    string
	Feed           string `gorm:"column:url"`
	Website        string `gorm:"column:link"`
	DownloadFolder string
}

func (Podcast) TableName() string {
	return "podcast"
}
