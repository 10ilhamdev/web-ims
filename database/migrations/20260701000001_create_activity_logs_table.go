package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"ims/app/facades"
)

type M20260701000001CreateActivityLogsTable struct{}

func (r *M20260701000001CreateActivityLogsTable) Signature() string {
	return "20260701000001_create_activity_logs_table"
}

func (r *M20260701000001CreateActivityLogsTable) Up() error {
	if !facades.Schema().HasTable("activity_logs") {
		if err := facades.Schema().Create("activity_logs", func(table schema.Blueprint) {
			table.ID()
			table.Integer("user_id")
			table.String("activity")
			table.String("device")
			table.String("ip")
			table.DateTime("created_at")
			table.DateTime("updated_at")
			table.DateTime("deleted_at").Nullable()
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *M20260701000001CreateActivityLogsTable) Down() error {
	return facades.Schema().DropIfExists("activity_logs")
}
