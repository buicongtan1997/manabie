package test_test

import (
	"bytes"
	"encoding/json"
	"github.com/buicongtan1997/manabie/pkg/api/v1/controllers"
	"github.com/buicongtan1997/manabie/pkg/api/v1/dto/request"
	"github.com/buicongtan1997/manabie/pkg/api/v1/dto/response"
	"github.com/buicongtan1997/manabie/pkg/database"
	"github.com/buicongtan1997/manabie/pkg/mappers"
	"github.com/buicongtan1997/manabie/pkg/services/v1"
	"github.com/buicongtan1997/manabie/test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var router = gin.New()

var wg sync.WaitGroup

func init() {
	api := router.Group("")
	controllers.RegisterProductController(api, v1.NewProductService())
}

func purseProduct(purchaseProduct request.PurchaseProduct) (product response.PurchaseProduct, code int) {
	requestByte, _ := json.Marshal(purchaseProduct)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/purchases", bytes.NewReader(requestByte))
	router.ServeHTTP(w, req)

	code = w.Code
	json.Unmarshal(w.Body.Bytes(), &product)
	return
}

func TestPurchasesProductApi(t *testing.T) {
	//prepare database test
	//init 100 product with each quantity: 10
	test.Before()

	productMapper := mappers.NewProductMapper(database.GetMainDatabase())

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

		resp, code := purseProduct(*req)
		assert.Equal(t, 200, code)
		assert.Equal(t, true, resp.Successful)

		for _, productRequest := range *req {
			product, _ := productMapper.FindById(productRequest.ProductID)
			assert.Equal(t, uint(test.DefaultQuantity - productRequest.Quantity), product.Quantity)
		}
	})

	t.Run("Request with id invalid", func(t *testing.T) {
		req := &request.PurchaseProduct{
			request.Product{
				ProductID: 1000,
				Quantity:  5,
			},
		}

		resp, code := purseProduct(*req)
		assert.Equal(t, 400, code)
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

		resp, code := purseProduct(*req)

		assert.Equal(t, 422, code)
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

		resp, code := purseProduct(*req)

		assert.Equal(t, 422, code)
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

		resp, code := purseProduct(*req1)
		assert.Equal(t, 200, code)
		assert.Equal(t, true, resp.Successful)

		for _, productRequest := range *req1 {
			product, _ := productMapper.FindById(productRequest.ProductID)
			assert.Equal(t, uint(test.DefaultQuantity - productRequest.Quantity), product.Quantity)
		}

		resp, code = purseProduct(*req1)
		assert.Equal(t, 200, code)
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
			code1, code2 int
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
			resp1, code1 = purseProduct(*req1)
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

			resp2, code2 = purseProduct(*req2)
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
