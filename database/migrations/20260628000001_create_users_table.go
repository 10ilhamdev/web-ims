package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"ims/app/facades"
)

type M20260628000001CreateUsersTable struct{}

func (r *M20260628000001CreateUsersTable) Signature() string {
	return "20260628000001_create_users_table"
}

func (r *M20260628000001CreateUsersTable) Up() error {
	if !facades.Schema().HasTable("users") {
		if err := facades.Schema().Create("users", func(table schema.Blueprint) {
			table.ID()
			table.String("name")
			table.String("email")
			table.String("password")
			table.String("role").Default("client")
			table.DateTime("created_at")
			table.DateTime("updated_at")
			table.DateTime("deleted_at").Nullable()
			table.Unique("email")
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260628000001CreateUsersTable) Down() error {
	return facades.Schema().DropIfExists("users")
}
