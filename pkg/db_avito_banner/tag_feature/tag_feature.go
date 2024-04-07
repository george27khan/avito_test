package tag_feature

import (
	db "avito_test/pkg/db_avito_banner"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type TagFeauture struct {
	Id        int64
	IdTag     int64
	IdFeature int64
	CreatedDt time.Time
}

// Insert функция для добавление записи в таблицу
func (tf *TagFeauture) Insert(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "INSERT INTO avito_banner.tag_feature(id_tag, id_feature) VALUES (@id, @id_tag, @id_feature)"

	args := pgx.NamedArgs{
		"id_tag":     tf.IdTag,
		"id_feature": tf.IdFeature}

	if res, err := conn.Exec(ctx, query, args); err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println(res)
	}
	return nil
}

// Delete функция для удаления записи из таблицы
func (tf *TagFeauture) Delete(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "delete from avito_banner.tag_feature t where t.id = $1"
	_, err = conn.Exec(ctx, query, tf.Id)
	if err := conn.Ping(ctx); err != nil {
		return err
	}
	return nil
}
