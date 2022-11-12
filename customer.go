package main

import (
	"time"
)

type Customer struct {
	Id          int64     `json:"id"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	PhoneNumber string    `json:"phone_number"`
	Gender      bool      `json:"gender"`
	BirthDate   time.Time `json:"birth_date"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type GetAllCustomer struct {
	Customers []*CreateOrGetCustomer `json:"customers"`
	Count     int         `json:"count"`
}

type CustomerParams struct {
	Limit   int64  `json:"limit"`
	Page    int64  `json:"page"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func (d *DBManager) CreateCustomer(customer *CreateOrGetCustomer) (*CreateOrGetCustomer, error) {
	var result CreateOrGetCustomer
	queryInsert := `
		INSERT INTO customer(
			firstname,
			lastname,
			phone_number,
			gender,
			birth_date,
			balance
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, firstname, lastname, phone_number, gender, birth_date, balance, created_at
	`
	err := d.db.QueryRow(
		queryInsert,
		customer.FirstName,
		customer.LastName,
		customer.PhoneNumber,
		customer.Gender,
		customer.BirthDate,
		customer.Balance,
	).Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Gender,
		&result.BirthDate,
		&result.Balance,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *DBManager) GetCustomer(customer_id int64) (*CreateOrGetCustomer, error) {
	var result CreateOrGetCustomer
	queryGet := `select id, firstname, lastname, phone_number, gender, birth_date, balance, created_at from customer where id  = $1 and deleted_at is null`
	err := d.db.QueryRow(
		queryGet,
		customer_id,
	).Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Gender,
		&result.BirthDate,
		&result.Balance,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *DBManager) UpdateCustomer(customer *Customer) (*CreateOrGetCustomer, error) {
	var result CreateOrGetCustomer
	queryUpdate := `
	UPDATE customer SET
		firstname = $1, 
		lastname = $2,  
		phone_number = $3, 
		gender = $4, 
		birth_date = $5, 
		balance = $6, 
		updated_at = $7 
	WHERE id = $8
	RETURNING id, firstname, lastname, phone_number, gender, birth_date, balance, created_at
	`
	err := d.db.QueryRow(
		queryUpdate,
		customer.FirstName,
		customer.LastName,
		customer.PhoneNumber,
		customer.Gender,
		customer.BirthDate,
		customer.Balance,
		time.Now(),
		customer.Id,
	).Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.PhoneNumber,
		&result.Gender,
		&result.BirthDate,
		&result.Balance,
		&result.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *DBManager) DeleteCustomer(customer_id int64) error {
	queryDelete := `
		delete from customer where id = $1
	`
	_, err := d.db.Exec(queryDelete, customer_id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBManager) GetCustomers(param *CustomerParams) (*GetAllCustomer, error){
	var customers GetAllCustomer
	customers.Customers = make([]*CreateOrGetCustomer, 0)
	filter := " WHERE deleted_at is null  "
	offset := (param.Page - 1) * param.Limit
	if param.Name != "" {
		filter += " AND firstname ilike '%s'" + "%" + param.Name + "%"
	}
	if param.Surname != "" {
		filter += " AND lastname ilike '%s'" + "%" + param.Surname + "%"
	}
	queryGetAll := `
		SELECT 
			id,
			firstname,
			lastname,
			phone_number,
			gender,
			birth_date,
			balance,
			created_at
		FROM customer 
	` + filter + `
		LIMIT $1 OFFSET $2
	`
	rows, err := d.db.Query(
		queryGetAll,
		param.Limit,
		offset,
	)
	defer rows.Close()
	for rows.Next() {
		var result CreateOrGetCustomer
		err := rows.Scan(
			&result.Id,
			&result.FirstName,
			&result.LastName,
			&result.PhoneNumber,
			&result.Gender,
			&result.BirthDate,
			&result.Balance,
			&result.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		customers.Customers = append(customers.Customers, &result)
	}
	if err != nil {
		return nil, err
	}
	queryCount := `SELECT count(*) FROM customer` + filter
	err = d.db.Get(&customers.Count, queryCount)
	if err != nil {
		return nil, err
	}
	return &customers, nil
}
