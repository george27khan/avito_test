package tag_feature

import (
	db "avito_test/pkg/db_avito_banner"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type TagFeauture struct {
	TagFeatureId int64
	TagId        int64
	FeatureId    int64
	CreatedDt    time.Time
}

// Insert функция для добавление записи в таблицу
func InsertAll(ctx context.Context, tx pgx.Tx, tags []int64, feature int64) error {
	query := "INSERT INTO avito_banner.tag_feature(tag_id, feature_id) VALUES ($1, $2)"
	for _, tag := range tags {
		if _, err := tx.Exec(ctx, query, tag, feature); err != nil {
			return err
		}
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
	query := "delete from avito_banner.tag_feature t where t.tag_feature_id = $1"
	_, err = conn.Exec(ctx, query, tf.TagFeatureId)
	if err := conn.Ping(ctx); err != nil {
		return err
	}
	return nil
}

// CheckData функция проверки существования связки тег-фича-баннер в базе
func CheckData(ctx context.Context, tags []int64, feature int64) error {
	var (
		row pgx.Row
		cnt int
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "select 1 from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.tag_id = $1 and tf.feature_id = $2"

	for _, tag := range tags {
		row = conn.QueryRow(ctx, query, tag, feature)
		if err := row.Scan(cnt); err != nil {
			return fmt.Errorf("Баннер для тега %v и фичи %v уже определен.")
		}
	}
	return nil
}
