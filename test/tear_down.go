package test

import (
	"github.com/buicongtan1997/manabie/pkg/database"
	"github.com/buicongtan1997/manabie/pkg/models"
)

func TearDown(){
	database.GetMainDatabase().DropTable(&models.Product{})
}
