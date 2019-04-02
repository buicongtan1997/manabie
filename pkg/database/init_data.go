package database

import "github.com/buicongtan1997/manabie/pkg/models"

func initData(){
	for i := 0; i < 100; i ++ {
		productId := uint(i + 1)
		db.FirstOrCreate(&models.Product{
			ID: productId,
			Quantity: 10,
		}, "id = ?", productId)
	}
}
