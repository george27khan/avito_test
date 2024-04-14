package connection

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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

func getPGconnStrMigr() (connStr string) {
	// loads DB settings from .env into the system
	if err := godotenv.Load("./.env"); err != nil {
		log.Print("No .env file found")
	}
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPwd := os.Getenv("POSTGRES_PASSWORD")
	host := "localhost"

	//"postgres://username:password@localhost:5432/database_name"
	connStr = fmt.Sprintf("postgres://%s:%s@%s:5432/%s", dbUser, dbPwd, host, dbName)
	fmt.Println("connStr ", connStr)
	return
}

var (
	pgOnce sync.Once
)

func Pool(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, getPGconnStrMigr())
}

func ConnectLocal(ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx, getPGconnStrMigr())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}

// ConnectPool возвращает соединение из пула
func ConnectPool(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := PGPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// OpenTx функция возвращает открытую транзакцию
func OpenTx(ctx context.Context, conn *pgxpool.Conn) (pgx.Tx, error) {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// ConnectPoolTrx возвращает соединение с открытой транзакцией из пула
func ConnectPoolTrx(ctx context.Context) (*pgxpool.Conn, pgx.Tx, error) {
	conn, err := ConnectPool(ctx)
	if err != nil {
		return nil, nil, err
	}
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, err
	}
	return conn, tx, nil
}

// создание пула соединений при запуске приложения
func init() {
	var err error
	PGPool, err = Pool(context.Background())
	if err != nil {
		fmt.Println(err)
	}
}
