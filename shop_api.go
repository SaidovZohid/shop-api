package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Router /shop/customer/create [post]
// @Summary Create customer
// @Description Create Customer
// @Tags customer
// @Accept json
// @Produce json
// @Param book body CreateOrGetCustomer true "Customer"
// @Success 200 {object} CustomerID
// @Failure 500 {object} ResponseError
func (h *handler) CreateCustomer(ctx *gin.Context) {
	var c CreateOrGetCustomer
	err := ctx.ShouldBindJSON(&c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return 
	}
	customer_id, err := Dbmanager.CreateCustomer(&c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"customer_id": customer_id,
	})
}

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

// @Router /shop/customer/update/{id} [put]
// @Summary Update customer by it's id
// @Description Update customer
// @Tags customer
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param book body CreateOrGetCustomer true "Customer" 
// @Success 200 {object} CreateOrGetCustomer
// @Failure 500 {object} ResponseError
func (h *handler) UpdateCustomer(ctx *gin.Context) {
	var c Customer
	var err error
	c.Id, err = strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return
	}
	err = ctx.ShouldBindJSON(&c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return
	}
	customer, err := Dbmanager.UpdateCustomer(&c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
	}
	
	ctx.JSON(http.StatusOK, customer)
}

// @Router /shop/customer/delete/{id} [delete]
// @Summary Delete customer by it's id
// @Description Delete customer
// @Tags customer
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} ResponseOK
// @Failure 500 {object} ResponseError
func (h *handler) DeleteCustomer(ctx *gin.Context) {
	customer_id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return 
	}
	err = Dbmanager.DeleteCustomer(customer_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, ResponseOK{
		Message: "Succesfully deleted!",
	})
}

// @Router /shop/customer/all [get]
// @Summary Get all Customer
// @Description Get All Customer
// @Tags customer
// @Produce json
// @Param limit query int true "limit"
// @Param page query int true "page"
// @Success 200 {object} GetAllCustomer
// @Failure 500 {object} ResponseError
func (h *handler) GetAllCustomer(ctx *gin.Context) {
	params, err := validateParams(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return
	}
	customers, err := Dbmanager.GetCustomers(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return 
	}
	ctx.JSON(http.StatusOK, customers)
}

func validateParams(ctx *gin.Context) (*CustomerParams, error) {
	var (
		limit int64 = 10
		page int64 = 1
		err error
	)
	if ctx.Query("limit") != "" {
		limit, err = strconv.ParseInt(ctx.Query("limit"), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	if ctx.Query("page") != "" {
		page, err = strconv.ParseInt(ctx.Query("page"), 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return &CustomerParams{
		Limit: limit,
		Page: page,
	}, nil
}

// @Router /shop/category/create [post]
// @Summary Create Category
// @Description Get All Customer
// @Tags category
// @Accept json
// @Produce json
// @Param category body CreateCategory true "Category"
// @Success 200 {object} Category
// @Failure 500 {object} ResponseError
func (h *handler) CreateCategory(ctx *gin.Context) {
	var c CreateCategory
	ctx.ShouldBindJSON(&c)
	category, err := Dbmanager.CreateCategory(&c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return 
	}
	ctx.JSON(http.StatusOK, category)
}

// @Router /shop/category/{id} [get]
// @Summary Get Category by ID
// @Description Get Category by ID
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "Id"
// @Success 200 {object} Category
// @Failure 500 {object} ResponseError
func (h *handler) GetCategory(ctx *gin.Context) {
	category_id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return 
	}
	category, err := Dbmanager.GetCategory(category_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, category)
}

// @Router /shop/category/update/{id} [put]
// @Summary Update Category with id
// @Description Update category with it's ID
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "Id"
// @Param category body CreateCategory true "Category"
// @Success 200 {object} Category
// @Failure 500 {object} ResponseError
func (h *handler) UpdateCategory(ctx *gin.Context) {
	var c Category
	var err error
	c.Id, err = strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return 
	}
	ctx.ShouldBindJSON(&c)
	category, err := Dbmanager.UpdateCategory(&c)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, category)
}