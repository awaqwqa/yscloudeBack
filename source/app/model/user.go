package model

type User struct {
	ID       int    `gorm:"primaryKey"`
	UserName string `gorm:"ot null;unique;size:255"`
	Password string `gorm:"not null"`
	Mobile   string `gorm:"unique;not null;"`
}
