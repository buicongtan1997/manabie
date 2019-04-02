package v1_test

import (
	"github.com/buicongtan1997/manabie/pkg/api/v1/dto/request"
	"github.com/buicongtan1997/manabie/pkg/api/v1/dto/response"
	"github.com/buicongtan1997/manabie/pkg/common/constants"
	"github.com/buicongtan1997/manabie/pkg/database"
	"github.com/buicongtan1997/manabie/pkg/mappers"
	"github.com/buicongtan1997/manabie/pkg/services/v1"
	"github.com/buicongtan1997/manabie/test"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
)

var productService = v1.NewProductService()

var wg sync.WaitGroup

func TestProductService_PurchaseProduct(t *testing.T) {
	productMapper := mappers.NewProductMapper(database.GetMainDatabase())
	//prepare database test
	//init 100 product with each quantity: 10
	test.Before()

	t.Run("Request 1 times", func(t *testing.T) {
		req := &request.PurchaseProduct{
			request.Product{
				ProductID: 1,
				Quantity:  5,
			},
			request.Product{
				ProductID: 2,
				Quantity:  8,
			},
			request.Product{
				ProductID: 3,
				Quantity:  10,
			},
		}

		resp, code := productService.PurchaseProduct("Request 1 times", req)
		assert.Equal(t, constants.ResponseCodeSuccess, code)
		assert.Equal(t, true, resp.Successful)

		for _, productRequest := range *req {
			product, _ := productMapper.FindById(productRequest.ProductID)
			assert.Equal(t, uint(test.DefaultQuantity - productRequest.Quantity), product.Quantity)
		}
	})

	t.Run("Request invalid id", func(t *testing.T) {
		req := &request.PurchaseProduct{
			request.Product{
				ProductID: 1111,
				Quantity:  5,
			},
		}

		resp, code := productService.PurchaseProduct("Request invalid id", req)
		assert.Equal(t, constants.ResponseCodeNotFound, code)
		assert.Equal(t, false, resp.Successful)
	})

	t.Run("Not Enough Product First Elem", func(t *testing.T) {
		req := &request.PurchaseProduct{
			request.Product{
				ProductID: 10,
				Quantity:  15,
			},
			request.Product{
				ProductID: 11,
				Quantity:  8,
			},
		}

		resp, code := productService.PurchaseProduct("Not Enough Product First Elem", req)
		assert.Equal(t, constants.ResponseCodeNotEnoughProduct, code)
		assert.Equal(t, false, resp.Successful)

		for _, productRequest := range *req {
			product, _ := productMapper.FindById(productRequest.ProductID)
			assert.Equal(t, uint(test.DefaultQuantity), product.Quantity)
		}
	})

	t.Run("Not Enough Product Second Elem", func(t *testing.T) {
		req := &request.PurchaseProduct{
			request.Product{
				ProductID: 15,
				Quantity:  1,
			},
			request.Product{
				ProductID: 16,
				Quantity:  15,
			},
		}

		resp, code := productService.PurchaseProduct("Not Enough Product Second Elem", req)
		assert.Equal(t, constants.ResponseCodeNotEnoughProduct, code)
		assert.Equal(t, false, resp.Successful)

		for _, productRequest := range *req {
			product, _ := productMapper.FindById(productRequest.ProductID)
			assert.Equal(t, uint(test.DefaultQuantity), product.Quantity)
		}
	})

	t.Run("Request 2 times success", func(t *testing.T) {
		req1 := &request.PurchaseProduct{
			request.Product{
				ProductID: 21,
				Quantity:  1,
			},
			request.Product{
				ProductID: 22,
				Quantity:  3,
			},
			request.Product{
				ProductID: 23,
				Quantity:  5,
			},
		}

		resp, code := productService.PurchaseProduct("Request 1", req1)
		assert.Equal(t, constants.ResponseCodeSuccess, code)
		assert.Equal(t, true, resp.Successful)

		for _, productRequest := range *req1 {
			product, _ := productMapper.FindById(productRequest.ProductID)
			assert.Equal(t, uint(test.DefaultQuantity - productRequest.Quantity), product.Quantity)
		}

		resp, code = productService.PurchaseProduct("Request 2", req1)
		assert.Equal(t, constants.ResponseCodeSuccess, code)
		assert.Equal(t, true, resp.Successful)

		for _, productRequest := range *req1 {
			product, _ := productMapper.FindById(productRequest.ProductID)
			assert.Equal(t, uint(test.DefaultQuantity - (productRequest.Quantity *2)), product.Quantity)
		}
	})

	t.Run("2 Request at the same times", func(t *testing.T) {
		var (
			req1, req2 *request.PurchaseProduct
			resp1, resp2 response.PurchaseProduct
			code1, code2 constants.ResponseCode
		)
		wg.Add(1)
		go func() {
			req1 = &request.PurchaseProduct{
				request.Product{
					ProductID: 30,
					Quantity:  5,
				},
				request.Product{
					ProductID: 31,
					Quantity:  6,
				},
				request.Product{
					ProductID: 32,
					Quantity:  7,
				},
			}
			resp1, code1 = productService.PurchaseProduct("Request 1", req1)
			wg.Done()
		}()

		func() {
			req2 = &request.PurchaseProduct{
				request.Product{
					ProductID: 30,
					Quantity:  7,
				},
				request.Product{
					ProductID: 31,
					Quantity:  5,
				},
				request.Product{
					ProductID: 32,
					Quantity:  4,
				},
			}
			resp2, code2 = productService.PurchaseProduct("Request 2", req2)
		}()

		wg.Wait()
		assert.NotEqual(t, code1, code2)
		assert.NotEqual(t, resp1.Successful, resp2.Successful)

		if resp1.Successful {
			log.Println("Success 1")
			for _, productRequest := range *req1 {
				product, _ := productMapper.FindById(productRequest.ProductID)
				assert.Equal(t, uint(test.DefaultQuantity - (productRequest.Quantity)), product.Quantity)
			}
		} else if resp2.Successful {
			log.Println("Success 2")
			for _, productRequest := range *req2 {
				product, _ := productMapper.FindById(productRequest.ProductID)
				assert.Equal(t, uint(test.DefaultQuantity - (productRequest.Quantity)), product.Quantity)
			}
		}
	})

	//clear all
	test.TearDown()
}
