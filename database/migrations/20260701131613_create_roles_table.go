package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ims/app/facades"
)

type M20260701131613CreateRolesTable struct{}

// Signature The unique signature for the migration.
func (r *M20260701131613CreateRolesTable) Signature() string {
	return "20260701131613_create_roles_table"
}

// Up Run the migrations.
func (r *M20260701131613CreateRolesTable) Up() error {
	if !facades.Schema().HasTable("roles") {
		return facades.Schema().Create("roles", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("table_name")
			table.String("model_name")
			table.Text("fields")
			table.Text("relations")
			table.TimestampsTz()
			table.Unique("name")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260701131613CreateRolesTable) Down() error {
	return facades.Schema().DropIfExists("roles")
}
