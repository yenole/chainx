package model

type Chain struct {
	ID   uint   `gorm:"PRIMARY_KEY"`
	CID  uint   `gorm:"NOT NULL"`
	URL  string `gorm:"NOT NULL"`
	Time int64  `gorm:"autoCreateTime"`
}
