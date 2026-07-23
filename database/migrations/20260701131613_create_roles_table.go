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
			table.String("label")
			table.String("table_name").Nullable()
			table.String("relation_name").Nullable()
			table.Boolean("is_system").Default(false)
			table.Boolean("is_registerable").Default(false)
			table.String("badge_color").Nullable()
			table.Text("description").Nullable()
			table.String("dashboard_route").Nullable()
			table.String("dashboard_view").Nullable()
			table.Timestamps()
			table.Unique("name")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260701131613CreateRolesTable) Down() error {
	return facades.Schema().DropIfExists("roles")
}
