package banner

import (
	db "avito_test/pkg/db_avito_banner"
	tf "avito_test/pkg/db_avito_banner/tag_feature"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

type Banner struct {
	Id        int64
	Content   []byte
	IsActive  bool
	FeatureId int64
	Tags      []int64
	CreatedDt time.Time
	UpdatedDt time.Time
}

// Insert функция для добавление записи в таблицу
func (b *Banner) Insert(ctx context.Context) (int64, error) {
	var idBanner int64
	conn, err := db.PGPool.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	query := "INSERT INTO avito_banner.banner(id, content, is_active, feature_id) " +
		"VALUES (@id, @Content, @is_active, @feature_id) RETURNING id"

	args := pgx.NamedArgs{
		"id":         b.Id,
		"Content":    b.Content,
		"is_active":  b.IsActive,
		"feature_id": b.FeatureId,
	}

	if err := tx.QueryRow(ctx, query, args).Scan(&idBanner); err != nil {
		return 0, err
	}
	if err := tf.InsertAll(ctx, tx, b.Tags, b.FeatureId); err != nil {
		return 0, err
	}
	return idBanner, nil
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

// Check функция проверки существования связки тег-фича-баннер в базе
func (b *Banner) Check(ctx context.Context) error {
	var (
		row pgx.Row
		cnt int
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "select 1 from avito_banner.tag_feature tf" +
		"join avito_banner.banner b on b.feature_id = tf.feature_id" +
		"where tf.tag_id = $1 and tf.feature_id = $2)"

	for _, tag := range b.Tags {
		row = conn.QueryRow(ctx, query, tag, b.FeatureId)
		if err := row.Scan(cnt); err != nil {
			return fmt.Errorf("Баннер для тега %v и фичи %v уже определен.", tag, b.FeatureId)
		}
	}
	return nil
}

// Get функция возращает баннер
func Get(ctx context.Context, tag int64, feature int64) (Banner, error) {
	var (
		banner Banner
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return banner, err
	}
	query := "select b.id, tf.tag_id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at" +
		"from avito_banner.tag_feature tf" +
		"join avito_banner.banner b on b.feature_id = tf.feature_id" +
		"where tf.tag_id = $1 and tf.feature_id = $2"
	row := conn.QueryRow(ctx, query, tag, feature)
	if err := row.Scan(&banner.Id, &banner.Tags, &banner.FeatureId, &banner.Content, &banner.IsActive, &banner.CreatedDt, &banner.UpdatedDt); err != nil {
		return banner, fmt.Errorf("Баннер для тега %v и фичи %v отсутствует.", tag, feature)
	}
	return banner, nil
}
