package bootstrap

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ims/database/migrations"
)

func Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateJobsTable{},
		&migrations.M20260628000001CreateUsersTable{},
		&migrations.M20260628000002CreateProductsTable{},
		&migrations.M20260628000003CreateOrdersTable{},
		&migrations.M20260701000001CreateActivityLogsTable{},
		&migrations.M20260701000004CreateCmsPagesTable{},
		&migrations.M20260701000002CreateGuestContentsTable{},
		&migrations.M20260701131613CreateRolesTable{},
	}
}
