package categoryTable

import "context"

type Storage interface {
	CategoryTableAddItems(ctx context.Context, p []CategoryTable) error
	CategoryTableUpdateItem(ctx context.Context, p CategoryTable) error
	CategoryTableDeleteItem(ctx context.Context, id int) error

	CategoryTableFindOne(ctx context.Context, id int) (CategoryTable, error)
	CategoryTableGetAll(ctx context.Context) ([]CategoryTable, error)
}
