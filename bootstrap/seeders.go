package bootstrap

import (
	"github.com/goravel/framework/contracts/database/seeder"

	"ims/database/seeders"
)

func Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
	}
}
