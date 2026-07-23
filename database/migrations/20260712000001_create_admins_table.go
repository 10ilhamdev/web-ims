package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"ims/app/facades"
)

type M20260712000001CreateAdminsTable struct{}

func (r *M20260712000001CreateAdminsTable) Signature() string {
	return "20260712000001_create_admins_table"
}

func (r *M20260712000001CreateAdminsTable) Up() error {
	if !facades.Schema().HasTable("admins") {
		if err := facades.Schema().Create("admins", func(table schema.Blueprint) {
			table.ID()
			table.BigInteger("user_id").Unsigned()
			table.String("phone").Nullable()
			table.String("department").Nullable()
			table.DateTime("created_at")
			table.DateTime("updated_at")
			
			table.Foreign("user_id").References("id").On("users").CascadeOnDelete()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260712000001CreateAdminsTable) Down() error {
	return facades.Schema().DropIfExists("admins")
}
