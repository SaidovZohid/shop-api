package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


// @Router /shop/customer/{id} [get]
// @Summary Get customer Info by id
// @Description Get Customer Info by id
// @Tags customer
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} CreateOrGetCustomer
// @Failure 500 {object} ResponseError
func (h *handler) GetCustomer(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return 
	}
	customer, err := h.db.GetCustomer(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}