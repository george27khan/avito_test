package banner

import (
	db "avito_test/pkg/db_avito_banner"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type Banner struct {
	Id         int64
	JsonValue  []byte
	IsActive   bool
	IdFeature  int64
	CreatedDt  time.Time
	UpdatesdDt time.Time
}

// Insert функция для добавление записи в таблицу
func (b *Banner) Insert(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "INSERT INTO avito_banner.banner(id, json_value, is_active, id_feature, created_dt, updated_dt) " +
		"VALUES (@id, @json_value, @is_active, @id_feature, @created_dt, @updated_dt)"

	args := pgx.NamedArgs{
		"id":         b.Id,
		"json_value": b.JsonValue,
		"is_active":  b.IsActive,
		"id_feature": b.IdFeature,
		"created_dt": b.CreatedDt,
		"updated_dt": b.UpdatesdDt}

	if res, err := conn.Exec(ctx, query, args); err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println(res)
	}
	return nil
}

// Delete функция для удаления записи из таблицы
func (b *Banner) Delete(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "delete from avito_banner.banner t where t.id = $1"
	_, err = conn.Exec(ctx, query, b.Id)
	if err := conn.Ping(ctx); err != nil {
		return err
	}
	return nil
}
