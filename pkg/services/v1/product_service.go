package v1

import (
	"fmt"
	"github.com/buicongtan1997/manabie/pkg/api/v1/dto/request"
	"github.com/buicongtan1997/manabie/pkg/api/v1/dto/response"
	"github.com/buicongtan1997/manabie/pkg/api/v1/services"
	"github.com/buicongtan1997/manabie/pkg/common/constants"
	"github.com/buicongtan1997/manabie/pkg/database"
	"github.com/buicongtan1997/manabie/pkg/logger"
	"github.com/buicongtan1997/manabie/pkg/mappers"
	"go.uber.org/zap"
)

type productService struct {
}

func NewProductService() services.IProductService {
	return &productService{}
}

func (*productService) PurchaseProduct(requestId string, req *request.PurchaseProduct) (resp response.PurchaseProduct, code constants.ResponseCode) {
	//open transaction
	tx := database.GetMainDatabase().Begin()

	//commit or rollback
	defer func() {
		if code == constants.ResponseCodeSuccess {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	productMapper := mappers.NewProductMapper(tx)
	for _, productRequest := range *req {
		//check id and quantity is missing
		if productRequest.ProductID == 0 || productRequest.Quantity == 0 {
			code = constants.ResponseCodeInvalidArgs
			return
		}

		//find product by id and lock this record for update
		product, err := productMapper.FindByIdForUpdate(productRequest.ProductID)
		if err != nil { //not found
			code = constants.ResponseCodeNotFound
			logger.Log.Error("ERR finding product: " + err.Error(), zap.String("request-id", requestId))
			return
		}

		//check quantity is enough
		if product.Quantity < productRequest.Quantity {
			code = constants.ResponseCodeNotEnoughProduct
			logger.Log.Error(fmt.Sprintf("ProductID %d not engough quanity", productRequest.ProductID), zap.String("request-id", requestId))
			return
		}

		//update quantity of product
		err = productMapper.DecreaseQuantity(productRequest.ProductID, productRequest.Quantity)
		if err != nil {
			code = constants.ResponseCodeUnknown
			logger.Log.Error("ERR update quantity: " + err.Error(), zap.String("request-id", requestId))
			return
		}
	}

	resp.Successful = true
	return
}
