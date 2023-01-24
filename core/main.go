package main

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func main() {
	fmt.Print("Test")
}

type ProductsTable struct {
	Id              int
	Category        int
	QuantityOfGoods int
	Lastprice       decimal.Decimal
	PictureAddress  string
}

type CategoryTable struct {
	Id           int
	CategoryId   int
	CategoryName string
}

type ProductPriceStory struct {
	Id    int
	Date  string
	Price decimal.Decimal
}

type UserPassword struct{
	Username string
	Password string
	
}