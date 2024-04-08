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
	BannerId  int64       `json:"banner_id"`
	TagIds    []int64     `json:"tag_ids"`
	FeatureId int64       `json:"feature_id"`
	Content   interface{} `json:"content"`
	IsActive  bool        `json:"is_active"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// Insert функция для добавление записи в таблицу
func (b *Banner) Insert(ctx context.Context) (int64, error) {
	var bannerId int64
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
	query := "INSERT INTO avito_banner.banner(content, is_active, feature_id) " +
		"VALUES (@Content, @is_active, @feature_id) RETURNING banner_id"

	args := pgx.NamedArgs{
		"banner_id":  b.BannerId,
		"Content":    b.Content,
		"is_active":  b.IsActive,
		"feature_id": b.FeatureId,
	}

	if err := tx.QueryRow(ctx, query, args).Scan(&bannerId); err != nil {
		return 0, err
	}
	if err := tf.InsertAll(ctx, tx, b.TagIds, b.FeatureId); err != nil {
		return 0, err
	}
	return bannerId, nil
}

// Delete функция для удаления записи из таблицы
func (b *Banner) Delete(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "delete from avito_banner.banner t where t.banner_id = $1"
	_, err = conn.Exec(ctx, query, b.BannerId)
	if err := conn.Ping(ctx); err != nil {
		return err
	}
	return nil
}

// Check функция проверки существования связки тег-фича-баннер в базе
func (b *Banner) Check(ctx context.Context) error {
	var (
		row pgx.Row
		cnt int8
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "select count(1) from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.tag_id = $1 and tf.feature_id = $2"

	for _, tag := range b.TagIds {
		row = conn.QueryRow(ctx, query, tag, b.FeatureId)
		if err := row.Scan(&cnt); err != nil {
			return fmt.Errorf("Ошибка при поиске создаваемого банера - %v.", err.Error())
		} else if cnt > 0 {
			return fmt.Errorf("Баннер для тега %v и фичи %v уже определен.", tag, b.FeatureId)
		}
	}
	return nil
}

// Get функция возращает баннер по тегу и фиче
func Get(ctx context.Context, tag int64, feature int64) (Banner, error) {
	var (
		banner Banner
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return banner, err
	}
	query := "select b.banner_id, array_agg(tf.tag_id), b.feature_id, b.content, b.is_active, b.created_at, b.updated_at " +
		"from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.tag_id = $1 and tf.feature_id = $2 " +
		"group by b.banner_id,b.feature_id, b.content, b.is_active, b.created_at, b.updated_at"
	row := conn.QueryRow(ctx, query, tag, feature)
	if err := row.Scan(&banner.BannerId, &banner.TagIds, &banner.FeatureId, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
		fmt.Println(err)
		return banner, fmt.Errorf("Баннер для тега %v и фичи %v отсутствует.", tag, feature)
	}
	return banner, nil
}

// GetByTag функция возращает баннер по тегу
func GetByTag(ctx context.Context, tag int64, limit int64, offset int64) ([]Banner, error) {
	var (
		bannerList []Banner
		banner     Banner
		rows       pgx.Rows
		queryErr   error
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return bannerList, err
	}
	query := "select b.banner_id, array_agg(tf.tag_id) as tag_ids, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at " +
		"from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.tag_id = $1 " +
		"group by b.banner_id,b.feature_id, b.content, b.is_active, b.created_at, b.updated_at"
	if limit > 0 && offset > 0 {
		query += " limit $2 offset $3"
		rows, queryErr = conn.Query(ctx, query, tag, limit, offset)
	} else if limit > 0 {
		query += " limit $2"
		rows, queryErr = conn.Query(ctx, query, tag, limit)
	} else if offset > 0 {
		query += " offset $2"
		rows, queryErr = conn.Query(ctx, query, tag, offset)
	} else {
		rows, queryErr = conn.Query(ctx, query, tag)
	}
	if queryErr != nil {
		return bannerList, queryErr
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&banner.BannerId, &banner.TagIds, &banner.FeatureId, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			fmt.Println(err)
			return []Banner{}, fmt.Errorf("Баннеры для тега %v отсутствуют.", tag)
		}
		bannerList = append(bannerList, banner)
	}
	return bannerList, nil
}

// GetByFeature функция возращает баннер по тегу
func GetByFeature(ctx context.Context, feature int64, limit int64, offset int64) ([]Banner, error) {
	var (
		bannerList []Banner
		banner     Banner
		rows       pgx.Rows
		queryErr   error
	)
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return bannerList, err
	}
	query := "select b.banner_id, array_agg(tf.tag_id) as tag_ids, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at " +
		"from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.feature_id = $1 " +
		"group by b.banner_id,b.feature_id, b.content, b.is_active, b.created_at, b.updated_at"
	if limit > 0 && offset > 0 {
		query += " limit $2 offset $3"
		rows, queryErr = conn.Query(ctx, query, feature, limit, offset)
	} else if limit > 0 {
		query += " limit $2"
		rows, queryErr = conn.Query(ctx, query, feature, limit)
	} else if offset > 0 {
		query += " offset $2"
		rows, queryErr = conn.Query(ctx, query, feature, offset)
	} else {
		rows, queryErr = conn.Query(ctx, query, feature)
	}
	if queryErr != nil {
		return bannerList, queryErr
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&banner.BannerId, &banner.TagIds, &banner.FeatureId, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
			fmt.Println(err)
			return []Banner{}, fmt.Errorf("Баннеры для фичи %v отсутствуют.", feature)
		}
		bannerList = append(bannerList, banner)
	}
	return bannerList, nil
}
