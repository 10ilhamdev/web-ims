package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"

	"ims/app/facades"
)

type M20260711000000CreateChatsTables struct{}

// Signature The unique signature for the migration.
func (r *M20260711000000CreateChatsTables) Signature() string {
	return "20260711000000_create_chats_tables"
}

// Up Run the migrations.
func (r *M20260711000000CreateChatsTables) Up() error {
	if !facades.Schema().HasTable("chats") {
		err := facades.Schema().Create("chats", func(table schema.Blueprint) {
			table.ID()
			table.String("token")
			table.String("name")
			table.String("email")
			table.Integer("user_id").Nullable()
			table.TimestampsTz()
		})
		if err != nil {
			return err
		}
	}

	if !facades.Schema().HasTable("chat_messages") {
		err := facades.Schema().Create("chat_messages", func(table schema.Blueprint) {
			table.ID()
			table.Integer("chat_id")
			table.String("sender_type") // "user" or "admin"
			table.Text("message")
			table.TimestampsTz()
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20260711000000CreateChatsTables) Down() error {
	err := facades.Schema().DropIfExists("chat_messages")
	if err != nil {
		return err
	}
	return facades.Schema().DropIfExists("chats")
}
