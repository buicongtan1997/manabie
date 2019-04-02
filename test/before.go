package test

import (
	"github.com/buicongtan1997/manabie/pkg/database"
	"github.com/buicongtan1997/manabie/pkg/models"
	"log"
	"os"
	"strings"
)

const DefaultQuantity = 10

func Before(){
	env := os.Getenv("MANABIE_ENV")
	if !strings.Contains(env, "test") {
		log.Fatalf("ENV is: [%s], Please set MANABIE_ENV=test ", env)
	}
	log.Println("Preparing data testing...")
	//drop table
	database.GetMainDatabase().DropTable(&models.Product{})
	//migrate new
	database.GetMainDatabase().AutoMigrate(&models.Product{})

	//init data
	for i := 0; i < 50; i ++ {
		productId := uint(i + 1)
		database.GetMainDatabase().FirstOrCreate(&models.Product{
			ID: productId,
			Quantity: DefaultQuantity,
		}, "id = ?", productId)
	}
}
