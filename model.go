package main

type CreateOrGetCustomer struct {
	FirstName   string  `json:"firstname"`
	LastName    string  `json:"lastname"`
	PhoneNumber string  `json:"phone_number"`
	Gender      bool    `json:"gender"`
	BirthDate   string  `json:"birth_date"`
	Balance     float64 `json:"balance"`
}

type CreateCategory struct {
	Name      string    `json:"name"`
	ImageUrl  string    `json:"image_url"`
}

type ResponseOK struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Message string `json:"message"`
}

type CustomerID int64
