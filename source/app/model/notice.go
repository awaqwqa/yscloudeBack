package model

type Notice struct {
	ID    uint `gorm:"primaryKey"`
	Value string
}
