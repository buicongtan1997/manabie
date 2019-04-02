package services

import (
	"github.com/buicongtan1997/manabie/pkg/api/v1/dto/request"
	"github.com/buicongtan1997/manabie/pkg/api/v1/dto/response"
	"github.com/buicongtan1997/manabie/pkg/common/constants"
)

type IProductService interface {
	PurchaseProduct(requestId string, req *request.PurchaseProduct) (resp response.PurchaseProduct, code constants.ResponseCode)
}
