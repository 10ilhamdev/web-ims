package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"ims/app/facades"
)

type M20260701000002CreateGuestContentsTable struct{}

func (r *M20260701000002CreateGuestContentsTable) Signature() string {
	return "20260701000002_create_guest_contents_table"
}

func (r *M20260701000002CreateGuestContentsTable) Up() error {
	if !facades.Schema().HasTable("guest_contents") {
		if err := facades.Schema().Create("guest_contents", func(table schema.Blueprint) {
			table.ID()
			table.Integer("page_id")
			table.String("key")
			table.Text("value_id")
			table.Text("value_en")
			table.DateTime("created_at")
			table.DateTime("updated_at")
			table.DateTime("deleted_at").Nullable()
			table.Unique("key")
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260701000002CreateGuestContentsTable) Down() error {
	return facades.Schema().DropIfExists("guest_contents")
}
