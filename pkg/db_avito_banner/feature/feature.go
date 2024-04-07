package feature

import (
	db "avito_test/pkg/db_avito_banner"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type Feature struct {
	Id        int64
	CreatedDt time.Time
}

// Insert функция для добавление записи в таблицу
func (f *Feature) Insert(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "INSERT INTO avito_banner.feature(id, created_dt) VALUES (@id, @created_dt)"

	args := pgx.NamedArgs{
		"id":         f.Id,
		"created_dt": f.CreatedDt,
	}
	if res, err := conn.Exec(ctx, query, args); err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println(res)
	}
	return nil
}

// Delete функция для удаления записи из таблицы
func (f *Feature) Delete(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "delete from avito_banner.tag t where t.id = $1"
	_, err = conn.Exec(ctx, query, f.Id)
	if err := conn.Ping(ctx); err != nil {
		return err
	}
	return nil
}
