package main

import (
	"fmt"
	"time"
)

type Product struct {
	Id            int64           `json:"id"`
	CategoryID    int             `json:"category_id"`
	Name          string          `json:"name"`
	Price         float64         `json:"price"`
	ImageUrl      string          `json:"image_url"`
	ProductImages []*ProductImage `json:"product_images"`
	CreatedAt     time.Time       `json:"created_at"`
}

type ProductImage struct {
	Id             int    `json:"id"`
	ImageUrl       string `json:"image_url"`
	SequenceNumber int    `json:"sequence_number"`
	ProductID      int    `json:"product_id"`
}

type ProductParams struct {
	Limit      int64  `json:"limit"`
	Page       int64  `json:"page"`
	CategoryId int64  `json:"category_id"`
	Name       string `json:"name"`
}

type AllProducts struct {
	Products []*Product `json:"products"`
	Count    int        `json:"count"`
}

func (d *DBManager) CreateProduct(product *Product) (int64, error) {
	queryInsertProduct := `
		INSERT INTO products(
			category_id,
			name,
			price,
			image_url
		) VALUES ($1, $2, $3, $4) RETURNING id
	`
	var product_id int64
	err := d.db.QueryRow(
		queryInsertProduct,
		product.CategoryID,
		product.Name,
		product.Price,
		product.ImageUrl,
	).Scan(&product_id)
	if err != nil {
		return 0, err
	}
	queryInsertImages := `
		INSERT INTO product_images(
			image_url,
			sequence_number,
			product_id
		) VALUES ($1, $2, $3) 
	`
	for _, v := range product.ProductImages {
		_, err := d.db.Exec(
			queryInsertImages,
			v.ImageUrl,
			v.SequenceNumber,
			product_id,
		)
		if err != nil {
			return 0, err
		}
	}
	return product_id, nil
}

func (d *DBManager) GetProduct(product_id int64) (*Product, error) {
	queryGetProduct := `
		SELECT 
			id,
			category_id,
			name,
			price,
			image_url,
			created_at
		FROM products WHERE id = $1
	`
	queryGetImages := `
		SELECT 
			id,
			image_url,
			sequence_number,
			product_id
		FROM product_images WHERE product_id = $1
	`
	var result Product
	err := d.db.QueryRow(
		queryGetProduct,
		product_id,
	).Scan(
		&result.Id,
		&result.CategoryID,
		&result.Name,
		&result.Price,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	rows, err := d.db.Query(queryGetImages, product_id)
	for rows.Next() {
		var image ProductImage
		err := rows.Scan(
			&image.Id,
			&image.ImageUrl,
			&image.SequenceNumber,
			&image.ProductID,
		)
		if err != nil {
			return nil, err
		}
		result.ProductImages = append(result.ProductImages, &image)
	}
	return &result, nil
}

func (d *DBManager) UpdateProduct(p *Product) (*Product, error) {
	queryUpdateProduct := `
		UPDATE products SET 
			category_id = $1,
			name = $2,
			price = $3,
			image_url = $4
		WHERE id = $5
		RETURNING id, category_id, name, price, image_url, created_at
	`
	var result Product
	err := d.db.QueryRow(
		queryUpdateProduct,
		p.CategoryID,
		p.Name,
		p.Price,
		p.ImageUrl,
		p.Id,
	).Scan(
		&result.Id,
		&result.CategoryID,
		&result.Name,
		&result.Price,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	queryDeleteImages := `
		DELETE FROM product_images WHERE product_id = $1
	`
	_, err = d.db.Exec(queryDeleteImages, p.Id)
	if err != nil {
		return nil, err
	}
	queryInsertImages := `
		INSERT INTO product_images(
			image_url,
			sequence_number,
			product_id
		) VALUES ($1, $2, $3)  
		RETURNING id, image_url, sequence_number, product_id
	`
	for _, v := range p.ProductImages {
		var image ProductImage
		err = d.db.QueryRow(
			queryInsertImages,
			v.ImageUrl,
			v.SequenceNumber,
			result.Id,
		).Scan(
			&image.Id,
			&image.ImageUrl,
			&image.SequenceNumber,
			&image.ProductID,
		)
		if err != nil {
			return nil, err
		}
		result.ProductImages = append(result.ProductImages, &image)
	}
	return &result, nil
}

func (d *DBManager) DeleteProduct(product_id int64) error {
	queryDeleteImages := "DELETE FROM product_images WHERE product_id = $1"
	queryDeleteProduct := "DELETE FROM products WHERE id = $1"
	_, err := d.db.Exec(queryDeleteImages, product_id)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(queryDeleteProduct, product_id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBManager) GetAllProducts(params *ProductParams) (*AllProducts, error) {
	var result AllProducts
	result.Products = make([]*Product, 0)
	filter := " WHERE true "
	offset := (params.Page - 1) * params.Limit
	if params.Name != "" {
		filter += " AND name ilike '%s'" + "%" + params.Name + "%"
	}
	if params.CategoryId != 0 {
		filter += fmt.Sprintf(" AND category_id = %d", params.CategoryId)
	}
	queryGetProducts := `
		SELECT 
			id, 
			category_id,
			name,
			price,
			image_url,
			created_at
		FROM products
	` + filter + `
		LIMIT $1 OFFSET $2
	`
	queryGetImages := `
		SELECT 
			id,
			image_url,
			sequence_number,
			product_id
		FROM product_images 
		WHERE product_id = $1
	`
	rows, err := d.db.Query(queryGetProducts, params.Limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.Id,
			&product.CategoryID,
			&product.Name,
			&product.Price,
			&product.ImageUrl,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		imageRows, err := d.db.Query(queryGetImages, product.Id)
		for imageRows.Next() {
			var image ProductImage
			err := imageRows.Scan(
				&image.Id,
				&image.ImageUrl,
				&image.SequenceNumber,
				&image.ProductID,
			)
			if err != nil {
				return nil, err
			}
			product.ProductImages = append(product.ProductImages, &image)
		}
		result.Products = append(result.Products, &product)
	}
	queryCount := "SELECT count(1) FROM products" + filter
	err = d.db.QueryRow(queryCount).Scan(&result.Count)
	if err != nil {
		return nil, err
	}
	return &result, nil
}