package productPriceStory

import "context"

type Storage interface {
	ProductPriceAddItems(ctx context.Context, p []ProductPriceStoryTable) error
	ProductPriceTableDeleteItem(ctx context.Context, id int) error

	ProductPriceTableFindOne(ctx context.Context, id int) (ProductPriceStoryTable, error)
	ProductPriceTableGetAll(ctx context.Context) ([]ProductPriceStoryTable, error)
}
