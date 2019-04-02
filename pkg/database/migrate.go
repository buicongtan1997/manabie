package database

import "github.com/buicongtan1997/manabie/pkg/models"

func migrate() {
	db.AutoMigrate(
		&models.Product{},
	)
}
