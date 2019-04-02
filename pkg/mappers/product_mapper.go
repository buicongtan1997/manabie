package mappers

import (
	"github.com/buicongtan1997/manabie/pkg/models"
	"github.com/jinzhu/gorm"
)

type productMapper struct {
	baseMapper
}

func NewProductMapper(db *gorm.DB) *productMapper {
	return &productMapper{
		baseMapper: baseMapper{db: db.Model(&models.Product{})},
	}
}

func (rcv *productMapper) FindById(id uint) (product models.Product, err error)  {
	err = rcv.db.First(&product, id).Error
	return
}

func (rcv *productMapper) FindByIdForUpdate(id uint) (product models.Product, err error)  {
	err = rcv.db.Set("gorm:query_option", "FOR UPDATE").First(&product, id).Error
	return
}

func (rcv *productMapper) DecreaseQuantity(id, quantity uint) error {
	return rcv.db.Where("id = ?", id).UpdateColumn("quantity", gorm.Expr("quantity - ?", quantity)).Error
}
