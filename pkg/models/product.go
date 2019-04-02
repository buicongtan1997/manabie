package models

type Product struct {
	ID       uint `gorm:"primary_key"`
	Quantity uint `gorm:"not null"`
}
