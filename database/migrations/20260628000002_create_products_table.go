package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"ims/app/facades"
)

type M20260628000002CreateProductsTable struct{}

func (r *M20260628000002CreateProductsTable) Signature() string {
	return "20260628000002_create_products_table"
}

func (r *M20260628000002CreateProductsTable) Up() error {
	if !facades.Schema().HasTable("products") {
		if err := facades.Schema().Create("products", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.Text("description")
			table.Double("price")
			table.Double("original_price")
			table.Double("discount")
			table.Text("features")
			table.String("image")
			table.DateTime("created_at")
			table.DateTime("updated_at")
			table.DateTime("deleted_at").Nullable()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260628000002CreateProductsTable) Down() error {
	return facades.Schema().DropIfExists("products")
}
