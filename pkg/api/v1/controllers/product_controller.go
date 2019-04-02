package controllers

import (
	"github.com/buicongtan1997/manabie/pkg/api/v1/dto/request"
	"github.com/buicongtan1997/manabie/pkg/api/v1/services"
	"github.com/buicongtan1997/manabie/pkg/common/constants"
	"github.com/buicongtan1997/manabie/pkg/logger"
	"github.com/buicongtan1997/manabie/pkg/protocol/rest/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

type productController struct {
	svc services.IProductService
}

func RegisterProductController(router *gin.RouterGroup, svc services.IProductService) {
	ctrl := &productController{svc}

	router.POST("purchases", ctrl.purchaseEndpoint)
}

func (rcv *productController) purchaseEndpoint(ctx *gin.Context) {
	req := new(request.PurchaseProduct)
	err := ctx.BindJSON(req)
	if err != nil {
		logger.Log.Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"successful": false})
		return
	}

	resp, code := rcv.svc.PurchaseProduct(middleware.GetReqID(ctx), req)
	switch code {
	case constants.ResponseCodeNotEnoughProduct:
		ctx.JSON(http.StatusUnprocessableEntity, resp)
		return
	case constants.ResponseCodeInvalidArgs, constants.ResponseCodeNotFound:
		ctx.JSON(http.StatusBadRequest, resp)
		return
	case constants.ResponseCodeUnknown:
		ctx.Status(http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, resp)
}
