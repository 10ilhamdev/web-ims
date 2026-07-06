package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"ims/app/facades"
)

type M20260701000004CreateCmsPagesTable struct{}

func (r *M20260701000004CreateCmsPagesTable) Signature() string {
	return "20260701000004_create_cms_pages_table"
}

func (r *M20260701000004CreateCmsPagesTable) Up() error {
	if !facades.Schema().HasTable("cms_pages") {
		if err := facades.Schema().Create("cms_pages", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("type")
			table.Integer("order").Default(1)
			table.DateTime("created_at")
			table.DateTime("updated_at")
			table.DateTime("deleted_at").Nullable()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260701000004CreateCmsPagesTable) Down() error {
	return facades.Schema().DropIfExists("cms_pages")
}
