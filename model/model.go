package model

type Product struct {
	Title        string
	CurrentPrice string
	OldPrice     string
	Url          string
	DiscountType string
}

type Products []Product
