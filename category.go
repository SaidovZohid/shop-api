package main

import (
	"time"
)

type Category struct {
	Id        int64       `json:"id"`
	Name      string    `json:"name"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCategory struct {
	Name      string    `json:"name"`
	ImageUrl  string    `json:"image_url"`
}

type AllCategories struct {
	Categories []*Category
}

func (d *DBManager) CreateCategory(category *CreateCategory) (*Category, error) {
	var result Category
	query := `
		INSERT INTO categories (
			name,
			image_url
		) VALUES ($1, $2)
		RETURNING id, name, image_url, created_at
	`
	err := d.db.QueryRow(query, category.Name, category.ImageUrl).Scan(
		&result.Id,
		&result.Name,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *DBManager) GetCategory(category_id int64) (*Category, error) {
	var result Category
	query := `
		SELECT 
			id,
			name,
			image_url,
			created_at
		FROM categories WHERE id = $1
	`
	err := d.db.QueryRow(query, category_id).Scan(
		&result.Id,
		&result.Name,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *DBManager) UpdateCategory(category *Category) (*Category, error) {
	var result Category
	query := `
		UPDATE categories SET 
		name = $1,
		image_url = $2
		WHERE id = $3
		RETURNING id, name, image_url, created_at
	`
	err := d.db.QueryRow(query, category.Name, category.ImageUrl, category.Id).Scan(
		&result.Id,
		&result.Name,
		&result.ImageUrl,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *DBManager) DeleteCategory(category_id int64) error {
	query := `
		DELETE FROM categories WHERE id = $1
	`
	_, err := d.db.Exec(query, category_id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBManager) GetCategories() (*AllCategories, error) {
	query := `
		SELECT 
			id,
			name,
			image_url,
			created_at
		FROM categories
	`
	var result AllCategories
	result.Categories = make([]*Category, 0)
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var category Category
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&category.ImageUrl,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result.Categories = append(result.Categories, &category)
	}
	return &result, nil
}