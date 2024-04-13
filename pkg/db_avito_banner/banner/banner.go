package banner

import (
	db "avito_test/pkg/db_avito_banner"
	bch "avito_test/pkg/db_avito_banner/banner_content_hist"
	tf "avito_test/pkg/db_avito_banner/tag_feature"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Banner struct {
	BannerId  int64     `json:"banner_id"`
	TagIds    []int64   `json:"tag_ids"`
	FeatureId int64     `json:"feature_id"`
	Content   string    `json:"content"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Insert функция для создания баннера
func (b *Banner) Insert(ctx context.Context) (int64, error) {
	var bannerId int64
	conn, tx, err := db.ConnectPoolTrx(ctx)
	if err != nil {
		return 0, fmt.Errorf(" Insert ошибка подключения к БД (%v)", err.Error())
	}
	defer conn.Release()
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	query := "INSERT INTO avito_banner.banner(content, is_active, feature_id) VALUES (@Content, @is_active, @feature_id) RETURNING banner_id"
	args := pgx.NamedArgs{
		"banner_id":  b.BannerId,
		"Content":    b.Content,
		"is_active":  b.IsActive,
		"feature_id": b.FeatureId,
	}
	if err := tx.QueryRow(ctx, query, args).Scan(&bannerId); err != nil {
		return 0, fmt.Errorf(" Insert - ошибка создания записи баннера (%v)", err.Error())
	}
	content := bch.BannerContentHist{BannerId: bannerId, Content: b.Content}
	//добавление json в историю
	if err := content.Insert(ctx, tx); err != nil {
		return 0, fmt.Errorf(" Insert (%v)", err.Error())
	}
	//добавление тегов
	if err := tf.Insert(ctx, tx, b.TagIds, b.FeatureId); err != nil {
		return 0, fmt.Errorf(" Insert (%v)", err.Error())
	}
	return bannerId, nil
}

// Delete функция для удаления баннера
func Delete(ctx context.Context, tx pgx.Tx, bannerId int64) error {
	query := "delete from avito_banner.banner t where t.banner_id = $1"
	if _, err := tx.Exec(ctx, query, bannerId); err != nil {
		return fmt.Errorf(" Delete: ошибка при удалении банера (%v)", err.Error())
	}
	return nil
}

// CheckTagFeature функция проверки существования связки тег-фича-баннер в базе
func (b *Banner) CheckTagFeature(ctx context.Context) error {
	var (
		row pgx.Row
		cnt int8
	)
	conn, err := db.PGPool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("CheckTag ошибка подключения к БД (%v)", err.Error())
	}
	defer conn.Release()

	query := "select count(1) from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.tag_id = $1 and tf.feature_id = $2"

	for _, tag := range b.TagIds {
		row = conn.QueryRow(ctx, query, tag, b.FeatureId)
		if err := row.Scan(&cnt); err != nil {
			return fmt.Errorf("ошибка при поиске создаваемого банера (%v)", err.Error())
		} else if cnt > 0 {
			return fmt.Errorf("баннер для тега %v и фичи %v уже определен", tag, b.FeatureId)
		}
	}
	return nil
}

// Exist функция проверки существования баннера в базе
func Exist(ctx context.Context, conn *pgxpool.Conn, bannerId int64) error {
	var cnt int8
	query := "select count(1) from avito_banner.banner b where b.banner_id = $1"
	row := conn.QueryRow(ctx, query, bannerId)

	if err := row.Scan(&cnt); err != nil || cnt == 0 {
		return fmt.Errorf(" Exist - баннер не найден")
	}
	return nil
}

// Get функция возращает баннер по id
func Get(ctx context.Context, conn *pgxpool.Conn, bannerId int64) (banner Banner, err error) {
	query := "select b.banner_id, b.feature_id, b.content, b.is_active, b.created_at, b.updated_at " +
		"from avito_banner.banner b " +
		"where b.banner_id = $1"
	if err = conn.QueryRow(ctx, query, bannerId).
		Scan(&banner.BannerId, &banner.FeatureId, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
		return banner, fmt.Errorf(" Get ошибка при поиске баннера (%v)", err.Error())
	}
	return
}

// GetContentByTagFeature функция возращает json баннер по тегу и фиче
func GetContentByTagFeature(ctx context.Context, conn *pgxpool.Conn, tag int64, feature int64) (content string, isActive bool, err error) {
	query := "select b.content, b.is_active " +
		"from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.tag_id = $1 and tf.feature_id = $2"
	row := conn.QueryRow(ctx, query, tag, feature)
	if err = row.Scan(&content, &isActive); err != nil {
		err = fmt.Errorf("в GetContentByTagFeature ошибка, баннер для тега %v и фичи %v отсутствует (%v)", tag, feature, err.Error())
		return
	}
	return
}

// GetAllContentByTagFeature функция возращает все версии баннера по тегу и фиче
func GetAllContentByTagFeature(ctx context.Context, conn *pgxpool.Conn, tag int64, feature int64) (string, bool, error) {
	var (
		content  string
		isActive bool
	)
	query := "select json_agg(b.content, h.version), b.is_active " +
		"from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"join avito_banner.banner_content_hist h on h.banner_id = b.banner_id " +
		"where tf.tag_id = $1 and tf.feature_id = $2"
	row := conn.QueryRow(ctx, query, tag, feature)
	if err := row.Scan(&content, &isActive); err != nil {
		return "", false, fmt.Errorf("баннер для тега %v и фичи %v отсутствует", tag, feature)
	}
	return content, isActive, nil
}

// GetByTagFeature функция возращает баннер по тегу и фиче
func GetByTagFeature(ctx context.Context, tag int64, feature int64) (Banner, error) {
	var banner Banner
	conn, err := db.PGPool.Acquire(ctx)
	if err != nil {
		return banner, fmt.Errorf("GetByTagFeature ошибка подключения к БД: %v", err.Error())
	}
	defer conn.Release()

	query := "select b.banner_id, array_agg(tf.tag_id), b.feature_id, b.content::text, b.is_active, b.created_at, b.updated_at " +
		"from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.tag_id = $1 and tf.feature_id = $2 " +
		"group by b.banner_id,b.feature_id, b.content::text, b.is_active, b.created_at, b.updated_at"
	row := conn.QueryRow(ctx, query, tag, feature)
	if err := row.Scan(&banner.BannerId, &banner.TagIds, &banner.FeatureId, &banner.Content, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt); err != nil {
		return banner, fmt.Errorf("баннер для тега %v и фичи %v отсутствует", tag, feature)
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
	if err != nil {
		return bannerList, fmt.Errorf("GetByTag ошибка подключения к БД: %v", err.Error())
	}
	defer conn.Release()

	query := "select b.banner_id, array_agg(tf.tag_id) as tag_ids, b.feature_id, b.content::text, b.is_active, b.created_at, b.updated_at " +
		"from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.tag_id = $1 " +
		"group by b.banner_id,b.feature_id, b.content::text, b.is_active, b.created_at, b.updated_at"
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
			return []Banner{}, fmt.Errorf("баннеры для тега %v отсутствуют", tag)
		}
		bannerList = append(bannerList, banner)
	}
	return bannerList, nil
}

// GetByFeature функция возращает баннер по фиче
func GetByFeature(ctx context.Context, feature int64, limit int64, offset int64) ([]Banner, error) {
	var (
		bannerList []Banner
		banner     Banner
		rows       pgx.Rows
		queryErr   error
	)
	conn, err := db.PGPool.Acquire(ctx)
	if err != nil {
		return bannerList, fmt.Errorf("GetByFeature ошибка подключения к БД (%v)", err.Error())
	}
	defer conn.Release()

	query := "select b.banner_id, array_agg(tf.tag_id) as tag_ids, b.feature_id, b.content::text, b.is_active, b.created_at, b.updated_at " +
		"from avito_banner.tag_feature tf " +
		"join avito_banner.banner b on b.feature_id = tf.feature_id " +
		"where tf.feature_id = $1 " +
		"group by b.banner_id,b.feature_id, b.content::text, b.is_active, b.created_at, b.updated_at"
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
			return []Banner{}, fmt.Errorf("GetByFeature баннеры для фичи %v отсутствуют (%v)", feature, err.Error())
		}
		bannerList = append(bannerList, banner)
	}
	return bannerList, nil
}

// UpdateField функция обновления полей банера
func (b *Banner) UpdateField(ctx context.Context, tx pgx.Tx, fieldName string, val interface{}) error {
	query := fmt.Sprintf("update avito_banner.banner "+
		"set %v = $1, updated_at = current_timestamp where banner_id = $2", fieldName)
	if _, err := tx.Exec(ctx, query, val, b.BannerId); err != nil {
		return fmt.Errorf("UpdateField ошибка обновления данных (%v)", err.Error())
	}
	//обновление версионности
	if fieldName == "content" {
		content := bch.BannerContentHist{BannerId: b.BannerId, Content: val}
		if err := content.Insert(ctx, tx); err != nil {
			return fmt.Errorf(" UpdateField (%v)", err.Error())
		}
	}
	return nil
}
