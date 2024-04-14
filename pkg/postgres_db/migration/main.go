package migration

import (
	db_con "avito_test/pkg/postgres_db/connection"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"os"
)

func getMigrator(ctx context.Context, conn *pgx.Conn) *migrate.Migrator {
	migrator, err := migrate.NewMigrator(ctx, conn, "avito_test")
	if err != nil {
		fmt.Printf("Unable to create a migrator: %v\n", err)
	}

	err = migrator.LoadMigrations(os.DirFS("./scripts/migration"))
	if err != nil {
		fmt.Printf("Unable to load migrations: %v\n", err)
	}
	return migrator
}

func InitDB() {
	ctx := context.Background()
	conn := db_con.Connect(ctx)
	defer conn.Close(ctx)

	migrator := getMigrator(ctx, conn)

	err := migrator.Migrate(ctx)
	if err != nil {
		fmt.Printf("Unable to migrate: %v\n", err)
	}

	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		fmt.Printf("Unable to get current schema version: %v\n", err)
	}

	fmt.Printf("Migration done. Current schema version: %v\n", ver)
}

func DropDB() {
	ctx := context.Background()
	conn := db_con.Connect(ctx)
	defer conn.Close(ctx)

	migrator := getMigrator(ctx, conn)

	err := migrator.MigrateTo(ctx, 0)
	if err != nil {
		fmt.Printf("Unable to migrate: %v\n", err)
	}

	ver, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		fmt.Printf("Unable to get current schema version: %v\n", err)
	}

	fmt.Printf("Migration done. Current schema version: %v\n", ver)
}

func run_migration() {
	DropDB()
	InitDB()
}
