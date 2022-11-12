package main

import (
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/stretchr/testify/require"
)

func createcustomer(t *testing.T) *CreateOrGetCustomer {
	result, err := dbManager.CreateCustomer(&CreateOrGetCustomer{
		FirstName:   faker.Name(),
		LastName:    faker.LastName(),
		PhoneNumber: faker.Phonenumber(),
		Gender:      true,
		BirthDate:   time.Date(2004, 12, 2, 13, 00, 00, 00, time.Local),
		Balance:     10000,
	})
	require.NoError(t, err)
	return result
}

func deleteCustomer(t *testing.T, customer_id int64) {
	err := dbManager.DeleteCustomer(customer_id)
	require.NoError(t, err)
}

func TestCreateCustomer(t *testing.T) {
	customer, err := dbManager.CreateCustomer(&CreateOrGetCustomer{
		FirstName:   faker.Name(),
		LastName:    faker.LastName(),
		PhoneNumber: faker.Phonenumber(),
		Gender:      true,
		BirthDate:   time.Date(2004, 12, 2, 13, 00, 00, 00, time.Local),
		Balance:     10000,
	})
	require.NoError(t, err)
	deleteCustomer(t, customer.Id)
	require.NotEmpty(t, customer)
}

func TestGetCustomer(t *testing.T) {
	customer := createcustomer(t)
	require.NotEmpty(t, customer)
	customer, err := dbManager.GetCustomer(customer.Id)
	require.NoError(t, err)
	deleteCustomer(t, customer.Id)
	require.NotEmpty(t, customer)
}

func TestUpdateCustomer(t *testing.T) {
	customer := createcustomer(t)
	require.NotEmpty(t, customer)
	c, err := dbManager.UpdateCustomer(&Customer{
		Id:          customer.Id,
		FirstName:   faker.Name(),
		LastName:    faker.LastName(),
		PhoneNumber: faker.Phonenumber(),
		Gender:      false,
		BirthDate:   time.Date(2004, 12, 2, 13, 00, 00, 00, time.Local),
		Balance:     10000,
	})
	require.NoError(t, err)
	deleteCustomer(t, c.Id)
	require.NotEmpty(t, c)
}

func TestDeleteCustomer(t *testing.T) {
	customer := createcustomer(t)
	err := dbManager.DeleteCustomer(customer.Id)
	require.NoError(t, err)
	require.NotEmpty(t, customer)
}

func TestGetCustomers(t *testing.T) {
	c := createcustomer(t)
	customers, err := dbManager.GetCustomers(&CustomerParams{
		Limit: 10,
		Page:  1,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(customers.Customers), 1)
	deleteCustomer(t, c.Id)
	require.NotEmpty(t, c)
}

func createcategory(t *testing.T) *Category {
	result, err := dbManager.CreateCategory(&CreateCategory{
		Name: "Iphone",
		ImageUrl: "~/zohid/image.jpg",
	})
	require.NoError(t, err)
	require.NotEmpty(t, result)
	return result
}

func deletecategory(t *testing.T, category_id int64) {
	err := dbManager.DeleteCategory(category_id)
	require.NoError(t, err)
}

func TestCreateCategory(t *testing.T) {
	result, err := dbManager.CreateCategory(&CreateCategory{
		Name: "Iphone",
		ImageUrl: "~/zohid/image.jpg",
	})
	require.NoError(t, err)
	deletecategory(t, result.Id)
	require.NotEmpty(t, result)
}

func TestGetCategory(t *testing.T) {
	category := createcategory(t)
	require.NotEmpty(t, category)
	c, err := dbManager.GetCategory(category.Id)
	require.NoError(t, err)
	require.NotEmpty(t, c)
	deletecategory(t, category.Id)
	require.NotEmpty(t, category)
}

func TestUpdateCategory(t *testing.T) {
	category := createcategory(t)
	c, err := dbManager.UpdateCategory(&Category{
		Id: category.Id,
		Name: "Apple",
		ImageUrl: "~/zohid/apple.jpg",
	})
	require.NoError(t, err)
	require.NotEmpty(t, c)
	deletecategory(t, category.Id)
	require.NotEmpty(t, category)
}

func TestDeleteCategory(t *testing.T) {
	category := createcategory(t)
	err := dbManager.DeleteCategory(category.Id)
	require.NotEmpty(t, category)
	require.NoError(t, err)
}

func TestGetCategories(t *testing.T) {
	category := createcategory(t)
	require.NotEmpty(t, category)
	categories, err := dbManager.GetCategories()
	require.GreaterOrEqual(t, len(categories.Categories), 1)
	require.NoError(t, err)
	require.NotEmpty(t, categories)
	deletecategory(t, category.Id)
	require.NotEmpty(t, category)
}

func createProduct(t *testing.T) int64 {
	product_id, err :=dbManager.CreateProduct(&Product{
		CategoryID: 1,
		Name: "Iphone 13 Pro Max",
		Price: 13000,
		ImageUrl: "~/zohid/iphone.jpg",
		ProductImages: []*ProductImage{
			{
				ImageUrl: "~/zohid/iphone1.jpg",
				SequenceNumber: 1,
			},
			{
				ImageUrl: "~/zohid/iphone2.jpg",
				SequenceNumber: 2,
			},
			{
				ImageUrl: "~/zohid/iphone3.jpg",
				SequenceNumber: 3,
			},
		},
	})
	require.NoError(t, err)
	return product_id
}

func deleteProduct(t *testing.T, product_id int64) {
	err := dbManager.DeleteProduct(product_id)
	require.NoError(t, err)
}

func TestCreateProduct(t *testing.T) {
	product_id := createProduct(t)
	deleteProduct(t, product_id)
}

func TestGetProduct(t *testing.T) {
	product_id := createProduct(t)
	product, err := dbManager.GetProduct(product_id)
	require.NoError(t, err)
	require.NotEmpty(t, product)
	deleteProduct(t, product_id)
}

func TestUpdateProduct(t *testing.T) {
	product_id := createProduct(t)
	product, err := dbManager.UpdateProduct(&Product{
		Id: product_id,
		CategoryID: 1,
		Name: "Iphone 14",
		Price: 1600,
		ImageUrl: "~/Desktop/iphone.jpg",
		ProductImages: []*ProductImage{
			{
				ImageUrl: "~/zohid/iphone1.jpg",
				SequenceNumber: 1,
			},
			{
				ImageUrl: "~/zohid/iphone3.jpg",
				SequenceNumber: 2,
			},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, product)
	deleteProduct(t, product_id)
}

func TestDeleteProduct(t *testing.T) {
	product_id := createProduct(t)
	err := dbManager.DeleteProduct(product_id)
	require.NoError(t, err)
}

func TestGetAllProducts(t *testing.T) {
	product_id := createProduct(t)
	products, err := dbManager.GetAllProducts(&ProductParams{
		Limit: 10,
		Page: 1,
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(products.Products), 1)
	deleteProduct(t, product_id)
}