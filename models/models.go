package models

import "time"

type User struct {
	Id          int64
	Username    string
	Fullname    string
	PhoneNumber string
	IsAdmin     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Product struct {
	Id          int64
	Name        string
	Price       float64
	Description string
	Quantity    int
	ImageUrl    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Cart struct {
	UserID int64
	Items  []CartItem
}

type CartItem struct {
	UserID    int64
	ProductID int64
	Quantity  int
}
