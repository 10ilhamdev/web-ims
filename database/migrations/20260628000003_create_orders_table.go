package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"ims/app/facades"
)

type M20260628000003CreateOrdersTable struct{}

func (r *M20260628000003CreateOrdersTable) Signature() string {
	return "20260628000003_create_orders_table"
}

func (r *M20260628000003CreateOrdersTable) Up() error {
	if !facades.Schema().HasTable("orders") {
		if err := facades.Schema().Create("orders", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("user_id")
			table.UnsignedBigInteger("product_id")
			table.Text("requirements")
			table.Double("price")
			table.String("status").Default("pending")
			table.DateTime("created_at")
			table.DateTime("updated_at")
			table.DateTime("deleted_at").Nullable()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260628000003CreateOrdersTable) Down() error {
	return facades.Schema().DropIfExists("orders")
}
