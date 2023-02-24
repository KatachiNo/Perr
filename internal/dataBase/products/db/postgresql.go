package db

import (
	"context"
	db "github.com/KatachiNo/Perr/internal/dataBase/products"
	"github.com/KatachiNo/Perr/pkg/client/postgresql"
)

type Storage interface {
	ProductsGetAll(ctx context.Context) (db.Products, error)
	ProductsAddItem(ctx context.Context, pTable db.Products) error
	ProductsUpdateItem(ctx context.Context, pTable db.Products) error
	ProductsDeleteItem(ctx context.Context, pTable db.Products) error
}

type repo struct {
	client postgresql.Client
}
