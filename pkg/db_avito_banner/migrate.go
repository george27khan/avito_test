package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

var PGPool *pgxpool.Pool

func getPGconnStr() (connStr string) {
	// loads DB settings from .env into the system
	if err := godotenv.Load("./.env"); err != nil {
		log.Print("No .env file found")
	}
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPwd := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")

	//"postgres://username:password@localhost:5432/database_name"
	connStr = fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPwd, host, dbName)
	fmt.Println("connStr ", connStr)
	return
}

var (
	pgOnce sync.Once
)

func Pool(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, getPGconnStr())
}

func Connect(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, getPGconnStr())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}

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
	conn := Connect(ctx)
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
	conn := Connect(ctx)
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

func init() {
	var err error
	PGPool, err = Pool(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	InitDB()
	//DropDB()
}
