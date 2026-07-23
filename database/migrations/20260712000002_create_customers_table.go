package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"ims/app/facades"
)

type M20260712000002CreateCustomersTable struct{}

func (r *M20260712000002CreateCustomersTable) Signature() string {
	return "20260712000002_create_customers_table"
}

func (r *M20260712000002CreateCustomersTable) Up() error {
	if !facades.Schema().HasTable("customers") {
		if err := facades.Schema().Create("customers", func(table schema.Blueprint) {
			table.ID()
			table.BigInteger("user_id").Unsigned()
			table.String("phone").Nullable()
			table.String("company_name").Nullable()
			table.String("address").Nullable()
			table.DateTime("created_at")
			table.DateTime("updated_at")
			
			table.Foreign("user_id").References("id").On("users").CascadeOnDelete()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260712000002CreateCustomersTable) Down() error {
	return facades.Schema().DropIfExists("customers")
}
