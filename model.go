package main

import "time"

type CreateOrGetCustomer struct {
	Id          int64     `json:"id"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	PhoneNumber string    `json:"phone_number"`
	Gender      bool      `json:"gender"`
	BirthDate   time.Time `json:"birth_date"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
}

type ResponseOK struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Message string `json:"message"`
}
